package sms

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Client 短信客户端
// 集成了限流、配额、重试等功能
type Client struct {
	provider     SMSProvider   // 短信服务商
	limiter      *RateLimiter  // 限流器
	quotaManager *QuotaManager // 配额管理
}

// ClientConfig 客户端配置
type ClientConfig struct {
	Redis         *redis.Client  // Redis 客户端（必须）
	Provider      SMSProvider    // 短信服务商（必须）
	LimiterConfig *LimiterConfig // 限流配置（可选，使用默认值）
	RetryConfig   *RetryConfig   // 重试配置（可选，使用默认值）
	EnableRetry   bool           // 是否启用重试（默认 false）
}

// NewClient 创建短信客户端
func NewClient(config *ClientConfig) *Client {
	if config.Redis == nil {
		panic("redis client is required")
	}
	if config.Provider == nil {
		panic("sms provider is required")
	}

	// 创建限流器
	limiter := NewRateLimiter(config.Redis, config.LimiterConfig)

	// 创建配额管理器
	quotaManager := NewQuotaManager(config.Redis)

	// 如果启用重试，包装 provider
	provider := config.Provider
	if config.EnableRetry {
		provider = NewRetryProvider(provider, config.RetryConfig)
	}

	return &Client{
		provider:     provider,
		limiter:      limiter,
		quotaManager: quotaManager,
	}
}

// Send 发送短信 会自动进行限流和配额检查
func (c *Client) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	// 1. 限流检查
	if err := c.limiter.CheckAndIncrement(ctx, req); err != nil {
		return nil, err
	}

	// 2. 配额检查
	if req.BizID != "" {
		if err := c.quotaManager.CheckAndIncrement(ctx, req.BizID); err != nil {
			return nil, err
		}
	}

	// 3. 发送短信
	return c.provider.Send(ctx, req)
}

// Verify 验证短信验证码
func (c *Client) Verify(ctx context.Context, req *VerifyRequest) (*VerifyResponse, error) {
	return c.provider.Verify(ctx, req)
}

// QueryStatus 查询短信发送状态
func (c *Client) QueryStatus(ctx context.Context, msgID string) (*StatusResponse, error) {
	return c.provider.QueryStatus(ctx, msgID)
}

// QueryStatusByPhone 查询手机号最近的短信状态（返回多条）
func (c *Client) QueryStatusByPhone(ctx context.Context, phone string) ([]*StatusResponse, error) {
	return c.provider.QueryStatusByPhone(ctx, phone)
}

// GetQuota 获取配额使用情况
func (c *Client) GetQuota(ctx context.Context, bizID string) (used int, max int, err error) {
	return c.quotaManager.GetQuota(ctx, bizID)
}

// SetQuota 设置业务配额
func (c *Client) SetQuota(bizID string, maxPerDay int) {
	c.quotaManager.SetQuota(bizID, maxPerDay)
}

// ResetQuota 重置配额
func (c *Client) ResetQuota(ctx context.Context, bizID string) error {
	return c.quotaManager.ResetQuota(ctx, bizID)
}

// GetPhoneCount 获取手机号发送次数
func (c *Client) GetPhoneCount(ctx context.Context, phone string, _type string) (int, error) {
	return c.limiter.GetPhoneCount(ctx, phone, _type)
}
