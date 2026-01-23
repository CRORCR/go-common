package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/gpencil/go-common/sms"
)

// CustomProvider 自定义短信服务商示例
// 演示如何实现 SMSProvider 接口
type CustomProvider struct {
	apiKey    string
	apiSecret string
	redis     *redis.Client
}

func NewCustomProvider(apiKey, apiSecret string, redis *redis.Client) *CustomProvider {
	return &CustomProvider{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		redis:     redis,
	}
}

func (p *CustomProvider) Send(ctx context.Context, req *sms.SendRequest) (*sms.SendResponse, error) {
	// 实现你的短信发送逻辑
	// 1. 调用第三方API
	// 2. 处理响应
	// 3. 根据错误码判断是否可重试

	fmt.Printf("发送短信到 %s, 模板: %s\n", req.Phone, req.Template)

	// 示例：模拟成功响应
	return &sms.SendResponse{
		MsgID:     "custom_msg_123",
		Success:   true,
		ErrorCode: "",
		ErrorMsg:  "",
	}, nil
}

func (p *CustomProvider) Verify(ctx context.Context, req *sms.VerifyRequest) (*sms.VerifyResponse, error) {
	// 实现验证码验证逻辑
	key := fmt.Sprintf("sms:code:%s:%s", req.BizID, req.Phone)
	storedCode, err := p.redis.Get(ctx, key).Result()
	if err != nil {
		return &sms.VerifyResponse{
			Success: false,
			ErrMsg:  "验证码不存在或已过期",
		}, nil
	}

	if storedCode != req.Code {
		return &sms.VerifyResponse{
			Success: false,
			ErrMsg:  "验证码错误",
		}, nil
	}

	p.redis.Del(ctx, key)
	return &sms.VerifyResponse{Success: true}, nil
}

func (p *CustomProvider) QueryStatus(ctx context.Context, msgID string) (*sms.StatusResponse, error) {
	// 实现状态查询逻辑
	return &sms.StatusResponse{
		MsgID:  msgID,
		Status: sms.StatusDelivered,
	}, nil
}

func (p *CustomProvider) QueryStatusByPhone(ctx context.Context, phone string) ([]*sms.StatusResponse, error) {
	// 实现通过手机号查询短信状态的逻辑
	// 这里返回该手机号最近的短信记录
	return []*sms.StatusResponse{
		{
			MsgID:  fmt.Sprintf("msg_%s_1", phone),
			Phone:  phone,
			Status: sms.StatusDelivered,
		},
	}, nil
}

func main() {
	ctx := context.Background()

	// 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// 创建自定义短信服务商
	customProvider := NewCustomProvider("your-api-key", "your-api-secret", rdb)

	// 使用装饰器模式添加重试功能
	providerWithRetry := sms.NewRetryProvider(customProvider, &sms.RetryConfig{
		MaxRetries: 3,
	})

	// 创建短信客户端
	client := sms.NewClient(&sms.ClientConfig{
		Redis:    rdb,
		Provider: providerWithRetry,
	})

	// 发送短信
	resp, err := client.Send(ctx, &sms.SendRequest{
		Phone:    "13800138000",
		Template: "TEMPLATE_001",
		BizID:    "login",
	})

	if err != nil {
		log.Fatalf("发送失败: %v", err)
	}

	fmt.Printf("发送成功: %s\n", resp.MsgID)
}
