package sms

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter 限流器
type RateLimiter struct {
	redis  *redis.Client
	config *LimiterConfig
}

// NewRateLimiter 创建限流器
func NewRateLimiter(redis *redis.Client, config *LimiterConfig) *RateLimiter {
	if config == nil {
		config = DefaultLimiterConfig()
	}
	return &RateLimiter{
		redis:  redis,
		config: config,
	}
}

// CheckAndIncrement 检查并增加计数
func (l *RateLimiter) CheckAndIncrement(ctx context.Context, req *SendRequest) error {
	// 1. 检查手机号限流
	if err := l.checkPhoneLimit(ctx, req.Phone); err != nil {
		return err
	}

	// 2. 检查设备限流
	if req.DeviceID != "" {
		if err := l.checkDeviceLimit(ctx, req.DeviceID); err != nil {
			return err
		}
	}

	// 3. 检查IP限流
	if req.IP != "" {
		if err := l.checkIPLimit(ctx, req.IP); err != nil {
			return err
		}
	}

	// 所有检查通过，增加计数
	if err := l.incrementCounters(ctx, req); err != nil {
		return err
	}

	return nil
}

// checkPhoneLimit 检查手机号限流
func (l *RateLimiter) checkPhoneLimit(ctx context.Context, phone string) error {
	now := time.Now()

	// 检查3分钟限制
	if l.config.PhonePerMinute > 0 {
		key := fmt.Sprintf("sms:limiter:phone:minute:%s:%s", phone, now.Format("200601021504"))
		count, err := l.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		if count >= l.config.PhonePerMinute {
			return ErrPhoneRateLimit
		}
	}

	// 检查1小时限制
	if l.config.PhonePerHour > 0 {
		key := fmt.Sprintf("sms:limiter:phone:hour:%s:%s", phone, now.Format("2006010215"))
		count, err := l.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		if count >= l.config.PhonePerHour {
			return ErrPhoneRateLimit
		}
	}

	// 检查24小时限制
	if l.config.PhonePerDay > 0 {
		key := fmt.Sprintf("sms:limiter:phone:day:%s:%s", phone, now.Format("20060102"))
		count, err := l.redis.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		if count >= l.config.PhonePerDay {
			return ErrPhoneRateLimit
		}
	}

	return nil
}

// checkDeviceLimit 检查设备限流
func (l *RateLimiter) checkDeviceLimit(ctx context.Context, deviceID string) error {
	if l.config.DevicePerDay == 0 {
		return nil
	}

	now := time.Now()
	key := fmt.Sprintf("sms:limiter:device:day:%s:%s", deviceID, now.Format("20060102"))
	count, err := l.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return err
	}

	if count >= l.config.DevicePerDay {
		return ErrDeviceRateLimit
	}

	return nil
}

// checkIPLimit 检查IP限流
func (l *RateLimiter) checkIPLimit(ctx context.Context, ip string) error {
	if l.config.IPPerDay == 0 {
		return nil
	}

	now := time.Now()
	key := fmt.Sprintf("sms:limiter:ip:day:%s:%s", ip, now.Format("20060102"))
	count, err := l.redis.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		return err
	}

	if count >= l.config.IPPerDay {
		return ErrIPRateLimit
	}

	return nil
}

// incrementCounters 增加所有计数器
func (l *RateLimiter) incrementCounters(ctx context.Context, req *SendRequest) error {
	now := time.Now()

	// 使用 pipeline 批量执行
	pipe := l.redis.Pipeline()

	// 手机号计数器
	if l.config.PhonePerMinute > 0 {
		key := fmt.Sprintf("sms:limiter:phone:minute:%s:%s", req.Phone, now.Format("200601021504"))
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Minute)
	}

	if l.config.PhonePerHour > 0 {
		key := fmt.Sprintf("sms:limiter:phone:hour:%s:%s", req.Phone, now.Format("2006010215"))
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, time.Hour)
	}

	if l.config.PhonePerDay > 0 {
		key := fmt.Sprintf("sms:limiter:phone:day:%s:%s", req.Phone, now.Format("20060102"))
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, 24*time.Hour)
	}

	// 设备计数器
	if req.DeviceID != "" && l.config.DevicePerDay > 0 {
		key := fmt.Sprintf("sms:limiter:device:day:%s:%s", req.DeviceID, now.Format("20060102"))
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, 24*time.Hour)
	}

	// IP计数器
	if req.IP != "" && l.config.IPPerDay > 0 {
		key := fmt.Sprintf("sms:limiter:ip:day:%s:%s", req.IP, now.Format("20060102"))
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, 24*time.Hour)
	}

	_, err := pipe.Exec(ctx)
	return err
}

// GetPhoneCount 获取手机号当前计数（用于调试）
func (l *RateLimiter) GetPhoneCount(ctx context.Context, phone string, _type string) (int, error) {
	now := time.Now()
	var key string

	switch _type {
	case "minute":
		key = fmt.Sprintf("sms:limiter:phone:minute:%s:%s", phone, now.Format("200601021504"))
	case "hour":
		key = fmt.Sprintf("sms:limiter:phone:hour:%s:%s", phone, now.Format("2006010215"))
	case "day":
		key = fmt.Sprintf("sms:limiter:phone:day:%s:%s", phone, now.Format("20060102"))
	default:
		return 0, fmt.Errorf("invalid period: %s", _type)
	}

	count, err := l.redis.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return count, err
}
