### 新功能

#### ✅ 完整实现阿里云短信服务商 (AliyunProvider)

基于阿里云官方 SDK（dysmsapi-20170525/v5）实现了完整的短信功能：

1. **发送短信** (`Send`)
   - 支持模板短信发送
   - 支持自定义签名（优先使用请求中的签名，否则使用配置的默认签名）
   - 支持国际短信（通过 CountryCode）
   - 自动存储验证码到 Redis
   - 完整的错误处理和错误码映射

2. **查询短信状态** (`QueryStatusByPhone`)
   - 支持通过手机号查询短信状态
   - 自动查询今天和昨天的短信记录
   - 解析阿里云返回的状态码（1-等待回执、2-发送失败、3-发送成功）
   - 注意：阿里云不支持仅通过 MsgID 查询，必须提供手机号

3. **错误处理**
   - 完整的阿里云错误码映射（15+ 种错误类型）
   - 区分可重试和不可重试错误
   - 支持从 SDK 错误中提取诊断建议

4. **安全性**
   - 禁用 SDK 内置重试，使用自定义重试机制
   - 避免重复重试导致的短信费用浪费

#### 示例代码

```go
// 创建阿里云服务商
provider, err := sms.NewAliyunProvider(redis, &sms.AliyunConfig{
    AccessKeyID:     "your-key",
    AccessKeySecret: "your-secret",
    SignName:        "你的签名",
})

// 发送短信
resp, err := client.Send(ctx, &sms.SendRequest{
    Phone:    "13800138000",
    Template: "SMS_123456",
    Params:   map[string]string{"code": "123456"},
    BizID:    "login",
})

// 查询状态
statuses, err := client.QueryStatusByPhone(ctx, "13800138000")
```

### 依赖更新

新增阿里云 SDK 依赖：
```
github.com/alibabacloud-go/dysmsapi-20170525/v5
github.com/alibabacloud-go/darabonba-openapi/v2
github.com/alibabacloud-go/tea-utils/v2
github.com/alibabacloud-go/tea
```

### 文档更新

- 新增 `sms/examples/aliyun_usage.go` - 完整的阿里云使用示例
- 更新 `README.md` - 添加阿里云使用说明
- 更新 `CHANGELOG.md` - 记录变更

### 注意事项

1. **阿里云查询限制**：
   - ❌ 不支持 `QueryStatus(msgID)` - 会返回错误
   - ✅ 必须使用 `QueryStatusByPhone(phone)`

2. **签名配置**：
   - 可以在 `AliyunConfig` 中配置默认签名
   - 也可以在每个请求中指定签名（优先级更高）

3. **国际短信**：
   - 需要在阿里云控制台开通国际短信权限
   - 设置 `CountryCode` 字段（如 "+1" 美国）

---

## 2024-12-08 优化

### 重要改进

#### 1. 添加国家代码支持
- `SendRequest` 新增 `CountryCode` 字段，支持国际短信
- 默认值为 `+86`（中国）
- 提供辅助方法：
  - `GetFullPhone()` - 获取完整手机号（国家代码+手机号）
  - `SetCountryCode()` - 链式设置国家代码
  - `NewSendRequest()` - 创建带默认值的请求

**示例：**
```go
// 国内短信（默认）
req := sms.NewSendRequest("13800138000", "SMS_123456",
    map[string]string{"code": "123456"})

// 国际短信
req := sms.NewSendRequest("1234567890", "SMS_123456",
    map[string]string{"code": "123456"}).
    SetCountryCode("+1")  // 美国
```

#### 2. 澄清模板参数 Params 的使用

**重要说明：** `Params` 是模板变量，不是直接发送给用户的内容！

- 在短信平台配置模板：`"您的验证码是${code}，请在5分钟内完成验证"`
- 代码传参：`Params: map[string]string{"code": "123456"}`
- 用户收到：`"您的验证码是123456，请在5分钟内完成验证"`

**用户不会看到** `{"code":"123456"}` 这样的 JSON 格式！

#### 3. 改进短信状态查询

##### SendResponse 优化
- 新增 `SentTime` 字段（time.Time 类型）
- `MsgID` 添加详细注释说明其重要性

##### StatusResponse 优化
- 新增 `Phone` 字段，方便查看是哪个号码的短信
- `MsgID` 添加注释说明：一个手机号可能有多条短信，MsgID 用于精确标识

##### 新增 QueryStatusByPhone 方法
支持通过手机号查询该号码的所有短信：

```go
// 方式1：通过 MsgID 精确查询（推荐）
status, err := client.QueryStatus(ctx, msgID)

// 方式2：通过手机号查询最近的短信（返回多条）
statuses, err := client.QueryStatusByPhone(ctx, "13800138000")
for _, status := range statuses {
    fmt.Printf("MsgID: %s, 状态: %d\n", status.MsgID, status.Status)
}
```

### 为什么 MsgID 必须保留？

1. **一个手机号可能发送多条短信**
   - 用户今天登录 3 次 = 3 条短信 = 3 个不同的 MsgID
   - 只有通过 MsgID 才能精确查询某一条短信的状态

2. **服务商需要 MsgID**
   - 阿里云、腾讯云等服务商返回的唯一标识
   - 用于查询发送状态、对账、退费等

3. **业务追踪**
   - 可以记录 MsgID 到数据库
   - 关联业务操作（如：订单ID -> MsgID）
   - 方便问题排查和统计

### 接口变更

#### SMSProvider 接口
新增方法：
```go
// QueryStatusByPhone 通过手机号查询最近的短信状态（可选实现）
// 注意：一个手机号可能有多条短信，此方法返回最近一条
QueryStatusByPhone(ctx context.Context, phone string) ([]*StatusResponse, error)
```

#### Client 接口
新增方法：
```go
// QueryStatusByPhone 查询手机号最近的短信状态（返回多条）
func (c *Client) QueryStatusByPhone(ctx context.Context, phone string) ([]*StatusResponse, error)
```

### 文件更新

- `types.go` - 核心数据结构优化
- `client.go` - 添加 QueryStatusByPhone 方法
- `provider_mock.go` - 实现 QueryStatusByPhone
- `provider_aliyun.go` - 实现 QueryStatusByPhone（示例代码）
- `retry.go` - 支持 QueryStatusByPhone 重试
- `README.md` - 添加详细说明文档

### 向后兼容性

所有改动都是向后兼容的：
- `CountryCode` 字段为可选，默认 `+86`
- 新增的方法不影响现有代码
- 原有 API 保持不变

### 升级建议

1. **如果发送国际短信**，添加 `CountryCode` 字段：
   ```go
   req.CountryCode = "+1"  // 美国
   ```

2. **如果需要通过手机号查询**，使用新方法：
   ```go
   statuses, err := client.QueryStatusByPhone(ctx, phone)
   ```

3. **更新文档理解**：理解 `Params` 是模板变量，不是发送内容
