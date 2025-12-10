package sms

import (
	"context"
	"time"
)

// SMSProvider 短信服务商接口
type SMSProvider interface {
	// Send 发送短信
	Send(ctx context.Context, req *SendRequest) (*SendResponse, error)

	// Verify 验证短信验证码
	Verify(ctx context.Context, req *VerifyRequest) (*VerifyResponse, error)

	// QueryStatus 通过消息ID查询短信发送状态
	QueryStatus(ctx context.Context, msgID string) (*StatusResponse, error)

	// QueryStatusByPhone 通过手机号查询最近的短信状态，一个手机号可能有多条短信，此方法返回最近一条
	QueryStatusByPhone(ctx context.Context, phone string) ([]*StatusResponse, error)
}

// SendRequest 发送短信请求
type SendRequest struct {
	Phone       string            // 手机号（不含国家代码，如：13800138000）
	CountryCode string            // 国家代码（如：+86，默认+86），用于国际短信
	Template    string            // 模板ID（在短信平台配置的模板编号）
	Params      map[string]string // 模板参数（用于替换模板中的${变量}）
	BizID       string            // 业务ID (login/register/pay等)
	DeviceID    string            // 设备ID（用于防刷）
	IP          string            // IP地址（用于防刷）
	SignName    string            // 签名名称（如：阿里云）
	OutID       string            // 外部ID，用于业务追踪
}

// SendResponse 发送短信响应
type SendResponse struct {
	MsgID     string // 消息ID（服务商返回的唯一标识，用于追踪和查询状态）
	Success   bool   // 是否成功
	ErrorCode string // 错误码
	ErrorMsg  string // 错误信息
}

// VerifyRequest 验证短信请求
type VerifyRequest struct {
	Phone string // 手机号
	Code  string // 验证码
	BizID string // 业务ID
}

// VerifyResponse 验证短信响应
type VerifyResponse struct {
	Success bool   // 是否验证成功
	ErrMsg  string // 错误信息
}

// StatusResponse 查询短信状态响应
type StatusResponse struct {
	MsgID       string        // 消息ID（必须！一个手机号可能有多条短信，通过MsgID精确标识）
	Phone       string        // 手机号（方便查看是哪个号码的短信）
	Status      MessageStatus // 消息状态
	SentTime    int64         // 发送时间（Unix时间戳）
	ReceiveTime int64         // 接收时间（Unix时间戳）
	ErrorMsg    string        // 错误信息
}

// MessageStatus 短信状态
type MessageStatus int

const (
	StatusUnknown   MessageStatus = 0 // 未知
	StatusPending   MessageStatus = 1 // 待发送
	StatusSent      MessageStatus = 2 // 已发送
	StatusDelivered MessageStatus = 3 // 已送达
	StatusFailed    MessageStatus = 4 // 发送失败
)

// SMSError 短信错误类型
type SMSError struct {
	Code      string // 错误码
	Message   string // 错误信息
	Retryable bool   // 是否可重试
	RawError  error  // 原始错误
}

func (e *SMSError) Error() string {
	if e.RawError != nil {
		return e.Message + ": " + e.RawError.Error()
	}
	return e.Message
}

// ErrorType 错误类型枚举
type ErrorType string

const (
	ErrorTypeTimeout      ErrorType = "timeout"       // 超时
	ErrorTypeRateLimit    ErrorType = "rate_limit"    // 限流
	ErrorTypeCircuitBreak ErrorType = "circuit_break" // 熔断
	ErrorTypeBalance      ErrorType = "balance"       // 余额不足
	ErrorTypeFormat       ErrorType = "format"        // 格式错误
	ErrorTypeInvalidPhone ErrorType = "invalid_phone" // 手机号无效
	ErrorTypeOther        ErrorType = "other"         // 其他错误
)

// LimiterConfig 限流配置
type LimiterConfig struct {
	// 基于手机号的限制
	PhonePerMinute int // 每分钟手机号限制
	PhonePerHour   int // 每小时手机号限制
	PhonePerDay    int // 每天手机号限制

	// 基于设备的限制
	DevicePerDay int // 每天设备限制

	// 基于IP的限制
	IPPerDay int // 每天IP限制
}

// DefaultLimiterConfig 默认限流配置
func DefaultLimiterConfig() *LimiterConfig {
	return &LimiterConfig{
		PhonePerMinute: 1,
		PhonePerHour:   3,
		PhonePerDay:    10,
		DevicePerDay:   10,
		IPPerDay:       10,
	}
}

// QuotaConfig 业务配额配置
type QuotaConfig struct {
	BizID     string // 业务ID
	MaxPerDay int    // 每个手机号每天最大次数
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries int           // 最大重试次数
	RetryDelay time.Duration // 重试延迟
}

// DefaultRetryConfig 默认重试配置
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxRetries: 3,
		RetryDelay: time.Second * 2,
	}
}

// ========== 辅助函数 ==========

// GetFullPhone 获取完整手机号（国家代码+手机号）
func (r *SendRequest) GetFullPhone() string {
	if r.CountryCode == "" {
		r.CountryCode = "+86" // 默认中国
	}
	return r.CountryCode + r.Phone
}

// SetCountryCode 设置国家代码（链式调用）
func (r *SendRequest) SetCountryCode(code string) *SendRequest {
	r.CountryCode = code
	return r
}

// NewSendRequest 创建发送请求（带默认值）
func NewSendRequest(phone, template string, params map[string]string) *SendRequest {
	return &SendRequest{
		Phone:       phone,
		CountryCode: "+86", // 默认中国
		Template:    template,
		Params:      params,
	}
}
