package sms

import "errors"

var (
	// 限流相关错误
	ErrPhoneRateLimit  = errors.New("手机号发送频率超限")
	ErrDeviceRateLimit = errors.New("设备发送频率超限")
	ErrIPRateLimit     = errors.New("IP发送频率超限")

	// 配额相关错误
	ErrQuotaExceeded = errors.New("业务配额已用尽")

	// 参数错误
	ErrInvalidParams = errors.New("无效的参数")

	// 服务错误
	ErrProviderFailed = errors.New("短信服务商调用失败")
	ErrTimeout        = errors.New("请求超时")
	ErrNetworkError   = errors.New("网络错误")

	// 业务错误
	ErrCodeExpired      = errors.New("验证码已过期")
	ErrCodeNotMatch     = errors.New("验证码不匹配")
	ErrBalanceNotEnough = errors.New("余额不足")
)

// NewSMSError 创建短信错误
func NewSMSError(code, message string, retryable bool, err error) *SMSError {
	return &SMSError{
		Code:      code,
		Message:   message,
		Retryable: retryable,
		RawError:  err,
	}
}

// IsRetryableError 判断错误是否可重试
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}

	var smsErr *SMSError
	if errors.As(err, &smsErr) {
		return smsErr.Retryable
	}

	// 白名单错误，不重试
	switch {
	case errors.Is(err, ErrTimeout):
		return false
	case errors.Is(err, ErrPhoneRateLimit):
		return false
	case errors.Is(err, ErrDeviceRateLimit):
		return false
	case errors.Is(err, ErrIPRateLimit):
		return false
	case errors.Is(err, ErrBalanceNotEnough):
		return false
	case errors.Is(err, ErrInvalidParams):
		return false
	default:
		// 其他错误默认可重试
		return true
	}
}

// GetErrorType 根据错误码获取错误类型
func GetErrorType(code string) ErrorType {
	switch code {
	case "TIMEOUT", "timeout", "RequestTimeout":
		return ErrorTypeTimeout
	case "RateLimit", "Throttling", "FlowControl":
		return ErrorTypeRateLimit
	case "ServiceUnavailable":
		return ErrorTypeCircuitBreak
	case "InsufficientBalance":
		return ErrorTypeBalance
	case "InvalidParameter":
		return ErrorTypeFormat
	case "InvalidPhoneNumber":
		return ErrorTypeInvalidPhone
	default:
		return ErrorTypeOther
	}
}

// ShouldRetry 根据错误类型判断是否应该重试
func ShouldRetry(errType ErrorType) bool {
	switch errType {
	case ErrorTypeTimeout:
		return false
	case ErrorTypeRateLimit:
		return false
	case ErrorTypeCircuitBreak:
		return false
	case ErrorTypeBalance:
		return false
	case ErrorTypeFormat:
		return false
	case ErrorTypeInvalidPhone:
		return false
	default:
		return true
	}
}
