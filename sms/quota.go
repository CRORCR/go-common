package sms

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// QuotaManager 配额管理器
type QuotaManager struct {
	redis  *redis.Client
	quotas map[string]*QuotaConfig // bizID -> QuotaConfig
}

// NewQuotaManager 创建配额管理器
func NewQuotaManager(redis *redis.Client) *QuotaManager {
	return &QuotaManager{
		redis:  redis,
		quotas: make(map[string]*QuotaConfig),
	}
}

const (
	defaultMax = 3
)

// SetQuota 设置业务配额
func (q *QuotaManager) SetQuota(bizID string, maxPerDay int) {
	q.quotas[bizID] = &QuotaConfig{
		BizID:     bizID,
		MaxPerDay: maxPerDay,
	}
}

// CheckAndIncrement 检查并增加配额计数
func (q *QuotaManager) CheckAndIncrement(ctx context.Context, bizID string) error {
	quota, exists := q.quotas[bizID]
	if !exists {
		// 如果没有配置该业务的配额，默认使用3次
		quota = &QuotaConfig{
			BizID:     bizID,
			MaxPerDay: defaultMax,
		}
	}

	// 检查配额
	if err := q.checkQuota(ctx, bizID, quota.MaxPerDay); err != nil {
		return err
	}

	// 增加计数
	if err := q.incrementQuota(ctx, bizID); err != nil {
		return err
	}

	return nil
}

// checkQuota 检查配额
func (q *QuotaManager) checkQuota(ctx context.Context, bizID string, maxPerDay int) error {
	now := time.Now()
	key := fmt.Sprintf("sms:quota:%s:%s", bizID, now.Format("20060102"))

	count, err := q.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return err
	}

	if count >= maxPerDay {
		return ErrQuotaExceeded
	}

	return nil
}

// incrementQuota 增加配额计数
func (q *QuotaManager) incrementQuota(ctx context.Context, bizID string) error {
	now := time.Now()
	key := fmt.Sprintf("sms:quota:%s:%s", bizID, now.Format("20060102"))

	pipe := q.redis.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, 24*time.Hour)

	_, err := pipe.Exec(ctx)
	return err
}

// GetQuota 获取当前配额使用情况
func (q *QuotaManager) GetQuota(ctx context.Context, bizID string) (used int, max int, err error) {
	quota, exists := q.quotas[bizID]
	if !exists {
		max = defaultMax // 默认值
	} else {
		max = quota.MaxPerDay
	}

	now := time.Now()
	key := fmt.Sprintf("sms:quota:%s:%s", bizID, now.Format("20060102"))

	count, err := q.redis.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, max, nil
	}
	if err != nil {
		return 0, max, err
	}

	return count, max, nil
}

// ResetQuota 重置配额（用于测试或管理后台）
func (q *QuotaManager) ResetQuota(ctx context.Context, bizID string) error {
	now := time.Now()
	key := fmt.Sprintf("sms:quota:%s:%s", bizID, now.Format("20060102"))
	return q.redis.Del(ctx, key).Err()
}
