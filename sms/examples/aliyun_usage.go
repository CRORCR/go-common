package main

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"

	"github.com/gpencil/go-common/sms"
)

func main() {
	// 1. 创建 Redis 客户端
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}

	// 2. 创建阿里云短信服务商
	provider, err := sms.NewAliyunProvider(rdb, &sms.AliyunConfig{
		AccessKeyID:     "your-access-key-id",
		AccessKeySecret: "your-access-key-secret",
		SignName:        "你的签名",                  // 在阿里云短信控制台配置的签名
		Endpoint:        "dysmsapi.aliyuncs.com", // 可选，默认值
		CodeExpiry:      0,                       // 可选，0 表示使用默认 5 分钟
	})
	if err != nil {
		log.Fatalf("创建阿里云短信客户端失败: %v", err)
	}

	// 3. 创建短信客户端（带限流和重试）
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

	// 5. 发送短信验证码
	sendResp, err := client.Send(ctx, &sms.SendRequest{
		Phone:       "13800138000",
		CountryCode: "+86",        // 可选，默认 +86
		Template:    "SMS_123456", // 阿里云短信模板 ID
		Params: map[string]string{
			"code": "123456", // 模板参数，替换模板中的 ${code}
		},
		BizID:    "login",
		DeviceID: "device_123",
		IP:       "192.168.1.1",
		SignName: "",          // 可选，为空则使用 AliyunConfig 中的默认签名
		OutID:    "order_123", // 可选，用于业务追踪
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

	// 7. 查询短信状态（通过手机号）
	statuses, err := client.QueryStatusByPhone(ctx, "13800138000")
	if err != nil {
		log.Printf("查询状态失败: %v", err)
		return
	}

	fmt.Printf("查询到 %d 条短信记录:\n", len(statuses))
	for i, status := range statuses {
		fmt.Printf("  [%d] MsgID: %s, 状态: %d, 错误: %s\n",
			i+1, status.MsgID, status.Status, status.ErrorMsg)
	}

	// 8. 查看配额使用情况
	used, max, err := client.GetQuota(ctx, "login")
	if err != nil {
		log.Printf("查询配额失败: %v", err)
		return
	}

	fmt.Printf("登录业务配额使用情况: %d/%d\n", used, max)
}
