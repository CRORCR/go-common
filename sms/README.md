# SMS 短信服务

企业级短信发送服务，支持多个短信服务商、限流、配额管理、重试机制等功能。

## 功能特性

- ✅ **多服务商支持**：统一接口，支持多个短信服务商（阿里云、腾讯云等）
- ✅ **防刷机制**：支持手机号、设备、IP 三维度限流
- ✅ **业务配额**：按业务类型（登录/注册/支付等）灵活配置配额
- ✅ **智能重试**：基于错误类型的智能重试策略
- ✅ **装饰器模式**：重试功能与基础功能解耦，灵活组合
- ✅ **验证码管理**：内置验证码存储和验证功能
- ✅ **状态查询**：支持查询短信发送状态

## 安装

```bash
go get ysgit.lunalabs.cn/products/go-common/sms
```

需要的依赖：
```bash
go get github.com/redis/go-redis/v9
```

## 重要说明

### 模板参数 Params 的使用

**重要！** `Params` 是模板变量，**不是直接发送给用户的内容**！

短信发送流程：
1. **在短信平台配置模板**（如阿里云控制台）
   ```
   模板内容：您的验证码是${code}，请在5分钟内完成验证。
   ```

2. **代码中传递 Params**
   ```go
   Params: map[string]string{
       "code": "123456",  // 替换模板中的 ${code}
   }
   ```

3. **用户收到的短信**
   ```
   您的验证码是123456，请在5分钟内完成验证。
   ```

**用户不会看到 `{"code":"123456"}` 这样的内容！** 这是常见的误解。

### 国家代码说明

支持国际短信发送，默认为中国（+86）：

```go
// 方式1：使用默认国家代码（+86）
req := &sms.SendRequest{
    Phone:    "13800138000",
    Template: "SMS_123456",
    Params:   map[string]string{"code": "123456"},
}

// 方式2：指定国家代码
req := &sms.SendRequest{
    Phone:       "13800138000",
    CountryCode: "+86",  // 中国
    Template:    "SMS_123456",
    Params:      map[string]string{"code": "123456"},
}

// 方式3：使用辅助函数
req := sms.NewSendRequest("13800138000", "SMS_123456",
    map[string]string{"code": "123456"}).
    SetCountryCode("+1")  // 美国
```

### MsgID 的作用

**为什么需要 MsgID？**

1. **一个手机号可能发送多条短信**
   - 用户今天登录 3 次 = 3 条短信 = 3 个不同的 MsgID
   - 只有通过 MsgID 才能精确查询某一条短信的状态

2. **用于追踪和对账**
   - 服务商（如阿里云）返回的唯一标识
   - 可以用于查询发送状态、对账、退费等

3. **支持两种查询方式**
   ```go
   // 方式1：通过 MsgID 精确查询（推荐）
   status, _ := client.QueryStatus(ctx, msgID)

   // 方式2：通过手机号查询最近的短信（返回多条）
   statuses, _ := client.QueryStatusByPhone(ctx, "13800138000")
   ```

## 快速开始

### 1. 基础使用

```go
package main

import (
    "context"
    "log"

    "github.com/redis/go-redis/v9"
    "ysgit.lunalabs.cn/products/go-common/sms"
)

func main() {
    // 创建 Redis 客户端
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

    // 创建短信服务商（使用模拟的服务商）
    provider := sms.NewMockProvider(rdb)

    // 创建短信客户端
    client := sms.NewClient(&sms.ClientConfig{
        Redis:       rdb,
        Provider:    provider,
        EnableRetry: true,  // 启用重试
    })

    // 设置业务配额
    client.SetQuota("login", 10)  // 登录业务每天10条

    // 发送短信
    ctx := context.Background()
    resp, err := client.Send(ctx, &sms.SendRequest{
        Phone:    "13800138000",
        Template: "SMS_123456",
        Params:   map[string]string{"code": "123456"},
        BizID:    "login",
        DeviceID: "device_123",
        IP:       "192.168.1.1",
    })

    if err != nil {
        log.Fatal(err)
    }

    log.Printf("发送成功: %s", resp.MsgID)
}
```

### 2. 验证码验证

```go
// 验证短信验证码
verifyResp, err := client.Verify(ctx, &sms.VerifyRequest{
    Phone: "13800138000",
    Code:  "123456",
    BizID: "login",
})

if err != nil {
    log.Fatal(err)
}

if verifyResp.Success {
    log.Println("验证成功")
} else {
    log.Printf("验证失败: %s", verifyResp.ErrMsg)
}
```

## 配置说明

### 限流配置

```go
client := sms.NewClient(&sms.ClientConfig{
    Redis:    rdb,
    Provider: provider,
    LimiterConfig: &sms.LimiterConfig{
        PhonePerMinute: 1,  // 每分钟每手机号限制 1 条
        PhonePerHour:   3,  // 每小时每手机号限制 3 条
        PhonePerDay:    10, // 每天每手机号限制 10 条
        DevicePerDay:   10, // 每天每设备限制 10 条
        IPPerDay:       10, // 每天每IP限制 10 条
    },
})
```

**频率控制说明：**
- **3分钟内同一手机号最多1条**：通过 `PhonePerMinute` 控制
- **1小时内最多3条**：通过 `PhonePerHour` 控制
- **24小时最多10条**：通过 `PhonePerDay` 控制
- 所有限制通过 Redis 的过期时间自动管理

### 重试配置

```go
client := sms.NewClient(&sms.ClientConfig{
    Redis:       rdb,
    Provider:    provider,
    EnableRetry: true,
    RetryConfig: &sms.RetryConfig{
        MaxRetries: 3,                    // 最大重试 3 次
        RetryDelay: time.Second * 2,      // 重试延迟 2 秒
    },
})
```

**重试策略：**
- ✅ **会重试**：网络错误、未知错误等
- ❌ **不重试**：超时、限流、熔断、余额不足、格式错误

### 业务配额配置

```go
// 设置不同业务的配额
client.SetQuota("login", 10)     // 登录：每天10条
client.SetQuota("register", 5)   // 注册：每天5条
client.SetQuota("pay", 20)       // 支付：每天20条
client.SetQuota("reset_pwd", 3)  // 重置密码：每天3条
```

## 防刷机制

系统实现了三维度的防刷控制：

### 1. 手机号维度
- 3分钟内最多 1 条
- 1小时内最多 3 条
- 24小时内最多 10 条

### 2. 设备维度
- 24小时内每设备最多 10 条

### 3. IP维度
- 24小时内每IP最多 10 条

### 建议的额外防刷措施

除了系统内置的三维度控制，建议在应用层增加以下防刷措施：

1. **图形验证码前置**：发送短信前要求用户完成图形验证码验证
2. **行为分析**：监控用户行为模式，识别异常请求
3. **用户信誉评分**：基于用户历史行为建立信誉体系
4. **实名认证**：要求用户进行实名认证后才能发送短信
5. **黑名单机制**：维护恶意手机号/IP/设备的黑名单

## 自定义短信服务商

### 实现 SMSProvider 接口

```go
type MyProvider struct {
    // 你的配置
}

func (p *MyProvider) Send(ctx context.Context, req *sms.SendRequest) (*sms.SendResponse, error) {
    // 实现发送逻辑
    // 1. 调用第三方API
    // 2. 处理错误码，返回正确的错误类型
    return &sms.SendResponse{
        MsgID:   "msg_123",
        Success: true,
    }, nil
}

func (p *MyProvider) Verify(ctx context.Context, req *sms.VerifyRequest) (*sms.VerifyResponse, error) {
    // 实现验证逻辑
}

func (p *MyProvider) QueryStatus(ctx context.Context, msgID string) (*sms.StatusResponse, error) {
    // 实现状态查询逻辑
}
```

### 使用装饰器模式添加重试

```go
// 创建你的服务商
myProvider := NewMyProvider(config)

// 使用装饰器模式添加重试功能
providerWithRetry := sms.NewRetryProvider(myProvider, &sms.RetryConfig{
    MaxRetries: 3,
})

// 创建客户端
client := sms.NewClient(&sms.ClientConfig{
    Redis:    rdb,
    Provider: providerWithRetry,  // 使用带重试的 provider
})
```

## 错误处理

### 错误类型

```go
const (
    ErrorTypeTimeout      // 超时（不重试）
    ErrorTypeRateLimit    // 限流（不重试）
    ErrorTypeCircuitBreak // 熔断（不重试）
    ErrorTypeBalance      // 余额不足（不重试）
    ErrorTypeFormat       // 格式错误（不重试）
    ErrorTypeInvalidPhone // 手机号无效（不重试）
    ErrorTypeOther        // 其他错误（重试）
)
```

### 自定义错误处理

```go
resp, err := client.Send(ctx, req)
if err != nil {
    // 检查是否是限流错误
    if errors.Is(err, sms.ErrPhoneRateLimit) {
        return errors.New("发送太频繁，请稍后再试")
    }

    // 检查是否是配额超限
    if errors.Is(err, sms.ErrQuotaExceeded) {
        return errors.New("今日发送次数已达上限")
    }

    return err
}
```

## Redis Key 设计

系统使用以下 Redis Key 格式：

```
# 限流相关
sms:limiter:phone:minute:{phone}:{YYYYMMDDHHmm}  # 3分钟过期
sms:limiter:phone:hour:{phone}:{YYYYMMDDHH}      # 1小时过期
sms:limiter:phone:day:{phone}:{YYYYMMDD}         # 24小时过期
sms:limiter:device:day:{deviceID}:{YYYYMMDD}     # 24小时过期
sms:limiter:ip:day:{ip}:{YYYYMMDD}               # 24小时过期

# 配额相关
sms:quota:{bizID}:{phone}:{YYYYMMDD}             # 24小时过期

# 验证码相关
sms:code:{bizID}:{phone}                         # 5分钟过期（可配置）
```

## 支持的短信服务商

- ✅ **MockProvider**（模拟服务商，用于测试）
- ✅ **AliyunProvider**（阿里云短信，已完整实现）
- 📝 TencentProvider（腾讯云，待实现）
- 📝 其他服务商...

### 阿里云短信使用

#### 1. 安装依赖

```bash
go get github.com/alibabacloud-go/dysmsapi-20170525/v5
go get github.com/alibabacloud-go/darabonba-openapi/v2
go get github.com/alibabacloud-go/tea-utils/v2
go get github.com/alibabacloud-go/tea
```

#### 2. 创建阿里云服务商

```go
provider, err := sms.NewAliyunProvider(redis, &sms.AliyunConfig{
    AccessKeyID:     "your-access-key-id",     // 阿里云 AccessKey ID
    AccessKeySecret: "your-access-key-secret", // 阿里云 AccessKey Secret
    SignName:        "你的签名",                 // 默认签名（在阿里云控制台配置）
    Endpoint:        "dysmsapi.aliyuncs.com",  // 可选，默认值
    CodeExpiry:      5 * time.Minute,          // 可选，验证码过期时间
})
if err != nil {
    log.Fatal(err)
}
```

#### 3. 发送短信

```go
resp, err := client.Send(ctx, &sms.SendRequest{
    Phone:       "13800138000",
    CountryCode: "+86",        // 国家代码，默认 +86
    Template:    "SMS_123456", // 阿里云短信模板 ID
    Params: map[string]string{
        "code": "123456", // 模板变量
    },
    BizID:    "login",
    SignName: "自定义签名", // 可选，为空则使用 AliyunConfig 中的默认签名
    OutID:    "order_123", // 可选，用于业务追踪
})
```

#### 4. 查询短信状态

**注意**：阿里云的查询接口需要提供手机号，因此：
- ✅ 使用 `QueryStatusByPhone(phone)` - 推荐
- ❌ `QueryStatus(msgID)` - 不支持（会返回错误）

```go
// 查询某个手机号的短信状态（会查询今天和昨天的记录）
statuses, err := client.QueryStatusByPhone(ctx, "13800138000")
for _, status := range statuses {
    fmt.Printf("MsgID: %s, 状态: %d\n", status.MsgID, status.Status)
}
```

#### 5. 完整示例

参见：`sms/examples/aliyun_usage.go`

## 最佳实践

1. **生产环境使用真实服务商**
   ```go
   provider := sms.NewAliyunProvider(rdb, &sms.AliyunConfig{
       AccessKeyID:     "your-key",
       AccessKeySecret: "your-secret",
   })
   ```

2. **合理配置重试次数**：建议不超过 3 次

3. **配置合适的限流策略**：根据业务特点调整频率限制

4. **监控配额使用情况**
   ```go
   used, max, _ := client.GetQuota(ctx, phone, bizID)
   log.Printf("配额使用: %d/%d", used, max)
   ```

5. **使用上下文控制超时**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   resp, err := client.Send(ctx, req)
   ```

## 目录结构

```
sms/
├── types.go              # 核心数据结构和接口定义
├── errors.go             # 错误定义
├── client.go             # 短信客户端
├── limiter.go            # 限流器
├── quota.go              # 配额管理器
├── retry.go              # 重试装饰器
├── provider_mock.go      # 模拟服务商
├── provider_aliyun.go    # 阿里云服务商
├── examples/             # 使用示例
│   ├── basic_usage.go
│   └── custom_provider.go
└── README.md             # 文档
```

## License

Copyright © 2024 Luna Labs
