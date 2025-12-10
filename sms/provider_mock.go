package sms

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// MockProvider 模拟短信服务商（用于测试）
type MockProvider struct {
	redis       *redis.Client
	codeExpiry  time.Duration // 验证码过期时间
	successRate float64       // 成功率（用于模拟失败）
	mu          sync.RWMutex
}

// NewMockProvider 创建模拟短信服务商
func NewMockProvider(redis *redis.Client) *MockProvider {
	return &MockProvider{
		redis:       redis,
		codeExpiry:  5 * time.Minute,
		successRate: 0.95, // 95%成功率
	}
}

// Send 发送短信
func (p *MockProvider) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	// 模拟随机失败
	if rand.Float64() > p.successRate {
		return nil, NewSMSError("NETWORK_ERROR", "网络错误", true, nil)
	}

	// 生成6位验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 存储验证码到 Redis
	key := fmt.Sprintf("sms:code:%s:%s", req.BizID, req.Phone)
	err := p.redis.Set(ctx, key, code, p.codeExpiry).Err()
	if err != nil {
		return nil, err
	}

	msgID := fmt.Sprintf("mock_%d", time.Now().UnixNano())

	return &SendResponse{
		MsgID:     msgID,
		Success:   true,
		ErrorCode: "",
		ErrorMsg:  "",
	}, nil
}

// Verify 验证短信验证码
func (p *MockProvider) Verify(ctx context.Context, req *VerifyRequest) (*VerifyResponse, error) {
	key := fmt.Sprintf("sms:code:%s:%s", req.BizID, req.Phone)

	// 从 Redis 获取验证码
	storedCode, err := p.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return &VerifyResponse{
			Success: false,
			ErrMsg:  "验证码不存在或已过期",
		}, nil
	}
	if err != nil {
		return nil, err
	}

	// 验证码匹配
	if storedCode != req.Code {
		return &VerifyResponse{
			Success: false,
			ErrMsg:  "验证码错误",
		}, nil
	}

	// 验证成功，删除验证码
	p.redis.Del(ctx, key)

	return &VerifyResponse{
		Success: true,
		ErrMsg:  "",
	}, nil
}

// QueryStatus 查询短信发送状态
func (p *MockProvider) QueryStatus(ctx context.Context, msgID string) (*StatusResponse, error) {
	// 模拟查询
	return &StatusResponse{
		MsgID:       msgID,
		Status:      StatusDelivered,
		SentTime:    time.Now().Unix(),
		ReceiveTime: time.Now().Unix(),
		ErrorMsg:    "",
	}, nil
}

// QueryStatusByPhone 通过手机号查询短信状态（返回最近的短信）
func (p *MockProvider) QueryStatusByPhone(ctx context.Context, phone string) ([]*StatusResponse, error) {
	// 模拟返回该手机号最近的短信记录
	return []*StatusResponse{
		{
			MsgID:       fmt.Sprintf("mock_%s_1", phone),
			Phone:       phone,
			Status:      StatusDelivered,
			SentTime:    time.Now().Unix() - 300, // 5分钟前
			ReceiveTime: time.Now().Unix() - 299, // 5分钟前
			ErrorMsg:    "",
		},
	}, nil
}
