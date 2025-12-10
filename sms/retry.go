package sms

import (
	"context"
	"fmt"
	"time"
)

// RetryProvider 重试
type RetryProvider struct {
	provider SMSProvider  // 被装饰的provider
	config   *RetryConfig // 重试配置
}

// NewRetryProvider 创建重试装饰器 默认3次重试，2s退避
func NewRetryProvider(provider SMSProvider, config *RetryConfig) *RetryProvider {
	if config == nil {
		config = DefaultRetryConfig()
	}
	return &RetryProvider{
		provider: provider,
		config:   config,
	}
}

// Send 发送短信（带重试）
func (r *RetryProvider) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			// 重试前延迟
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(r.config.RetryDelay * time.Duration(attempt)): // 指数退避
			}
		}

		resp, err := r.provider.Send(ctx, req)

		// 成功则直接返回
		if err == nil && resp != nil && resp.Success {
			return resp, nil
		}

		// 记录错误
		lastErr = err

		// 不应该重试，直接返回错误
		if !r.shouldRetry(err, resp) {
			return resp, err
		}

		// 如果是最后一次尝试，不再继续
		if attempt == r.config.MaxRetries {
			break
		}
	}

	// 所有重试都失败
	return nil, fmt.Errorf("短信发送失败，已重试%d次: %w", r.config.MaxRetries, lastErr)
}

// Verify 验证短信验证码（不需要重试）
func (r *RetryProvider) Verify(ctx context.Context, req *VerifyRequest) (*VerifyResponse, error) {
	return r.provider.Verify(ctx, req)
}

// QueryStatus 查询短信发送状态（带重试）
func (r *RetryProvider) QueryStatus(ctx context.Context, msgID string) (*StatusResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(r.config.RetryDelay * time.Duration(attempt)):
			}
		}

		resp, err := r.provider.QueryStatus(ctx, msgID)

		if err == nil {
			return resp, nil
		}

		lastErr = err

		// 查询失败可以重试（网络问题等）
		if !IsRetryableError(err) {
			return nil, err
		}

		if attempt == r.config.MaxRetries {
			break
		}
	}

	return nil, fmt.Errorf("查询短信状态失败，已重试%d次: %w", r.config.MaxRetries, lastErr)
}

// QueryStatusByPhone 通过手机号查询短信状态（带重试）
func (r *RetryProvider) QueryStatusByPhone(ctx context.Context, phone string) ([]*StatusResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(r.config.RetryDelay * time.Duration(attempt)):
			}
		}

		resp, err := r.provider.QueryStatusByPhone(ctx, phone)

		if err == nil {
			return resp, nil
		}

		lastErr = err

		// 查询失败可以重试（网络问题等）
		if !IsRetryableError(err) {
			return nil, err
		}

		if attempt == r.config.MaxRetries {
			break
		}
	}

	return nil, fmt.Errorf("查询短信状态失败，已重试%d次: %w", r.config.MaxRetries, lastErr)
}

// shouldRetry 判断是否应该重试
func (r *RetryProvider) shouldRetry(err error, resp *SendResponse) bool {
	if err != nil {
		return IsRetryableError(err)
	}

	// 检查错误码
	if resp != nil && !resp.Success {
		errType := GetErrorType(resp.ErrorCode)
		return ShouldRetry(errType)
	}

	// 其他情况不重试
	return false
}
