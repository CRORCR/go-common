package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"

	"github.com/CRORCR/go-common/sms"
)

func main() {
	// 1. 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// 测试 Redis 连接
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	// 2. 创建短信服务商（这里使用模拟的服务商）
	provider := sms.NewMockProvider(rdb)

	// 3. 创建短信客户端
	client := sms.NewClient(&sms.ClientConfig{
		Redis:    rdb,
		Provider: provider,
		LimiterConfig: &sms.LimiterConfig{
			PhonePerMinute: 1,  // 每分钟1条
			PhonePerHour:   3,  // 每小时3条
			PhonePerDay:    10, // 每天10条
			DevicePerDay:   10, // 每设备每天10条
			IPPerDay:       10, // 每IP每天10条
		},
		EnableRetry: true, // 启用重试
		RetryConfig: &sms.RetryConfig{
			MaxRetries: 3,
		},
	})

	// 4. 设置业务配额
	client.SetQuota("login", 10)   // 登录业务每天10条
	client.SetQuota("register", 5) // 注册业务每天5条
	client.SetQuota("pay", 20)     // 支付业务每天20条

	// 5. 发送短信
	sendResp, err := client.Send(ctx, &sms.SendRequest{
		Phone:    "13800138000",
		Template: "SMS_123456",
		Params: map[string]string{
			"code": "123456",
		},
		BizID:    "login",
		DeviceID: "device_123",
		IP:       "192.168.1.1",
	})

	if err != nil {
		log.Printf("发送短信失败: %v", err)
		return
	}

	fmt.Printf("短信发送成功！MsgID: %s\n", sendResp.MsgID)

	// 6. 验证短信验证码
	verifyResp, err := client.Verify(ctx, &sms.VerifyRequest{
		Phone: "13800138000",
		Code:  "123456",
		BizID: "login",
	})

	if err != nil {
		log.Printf("验证失败: %v", err)
		return
	}

	if verifyResp.Success {
		fmt.Println("验证码验证成功！")
	} else {
		fmt.Printf("验证码验证失败: %s\n", verifyResp.ErrMsg)
	}

	// 7. 查询短信状态
	statusResp, err := client.QueryStatus(ctx, sendResp.MsgID)
	if err != nil {
		log.Printf("查询状态失败: %v", err)
		return
	}

	fmt.Printf("短信状态: %d\n", statusResp.Status)

	// 8. 查看配额使用情况
	used, max, err := client.GetQuota(ctx, "login")
	if err != nil {
		log.Printf("查询配额失败: %v", err)
		return
	}

	fmt.Printf("配额使用情况: %d/%d\n", used, max)
}
