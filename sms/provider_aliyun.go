package sms

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	g_json "github.com/gpencil/go-common/json"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi "github.com/alibabacloud-go/dysmsapi-20170525/v5/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/redis/go-redis/v9"
)

// AliyunProvider 阿里云短信服务商
type AliyunProvider struct {
	client     *dysmsapi.Client
	redis      *redis.Client
	codeExpiry time.Duration
	signName   string // 签名名称
}

// AliyunConfig 阿里云配置
type AliyunConfig struct {
	AccessKeyID     string // AccessKey ID
	AccessKeySecret string // AccessKey Secret
	SignName        string // 签名名称
	Endpoint        string
	CodeExpiry      time.Duration // 验证码过期时间，默认 5 分钟
}

// NewAliyunProvider 创建阿里云短信服务商
func NewAliyunProvider(redis *redis.Client, config *AliyunConfig) (*AliyunProvider, error) {
	if config.CodeExpiry == 0 {
		config.CodeExpiry = 5 * time.Minute
	}
	if config.Endpoint == "" {
		config.Endpoint = "dysmsapi.aliyuncs.com"
	}

	if config.AccessKeyID == "" || config.AccessKeySecret == "" {
		return nil, errors.New("AccessKey不能为空")
	}

	// 创建客户端配置
	clientConfig := &openapi.Config{
		AccessKeyId:     tea.String(config.AccessKeyID),
		AccessKeySecret: tea.String(config.AccessKeySecret),
		Endpoint:        tea.String(config.Endpoint),
	}

	// 创建客户端
	client, err := dysmsapi.NewClient(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("创建阿里云短信客户端失败: %w", err)
	}

	return &AliyunProvider{
		client:     client,
		redis:      redis,
		codeExpiry: config.CodeExpiry,
		signName:   config.SignName,
	}, nil
}

// Send 发送短信
func (p *AliyunProvider) Send(ctx context.Context, req *SendRequest) (*SendResponse, error) {
	// 签名
	signName := req.SignName
	if signName == "" {
		signName = p.signName
	}

	// 获取完整手机号（包含国家代码）
	phone := req.GetFullPhone()

	// 转换模板参数为 JSON
	paramsJSON, err := g_json.Marshal(req.Params)
	if err != nil {
		return nil, NewSMSError("PARAM_ERROR", "模板参数序列化失败", false, err)
	}

	// 创建发送请求
	sendRequest := &dysmsapi.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(signName),
		TemplateCode:  tea.String(req.Template),
		TemplateParam: tea.String(string(paramsJSON)),
	}

	// 如果有外部 ID，设置
	if req.OutID != "" {
		sendRequest.OutId = tea.String(req.OutID)
	}

	// 配置运行时选项
	runtime := &util.RuntimeOptions{
		Autoretry:   tea.Bool(false), // 禁用 SDK 自动重试
		MaxAttempts: tea.Int(1),      // 只尝试一次
	}

	// 发送短信
	response, err := p.client.SendSmsWithOptions(sendRequest, runtime)

	if err != nil {
		return nil, p.handleSendError(err)
	}

	if response.Body == nil {
		return nil, NewSMSError("RESPONSE_ERROR", "响应体为空", true, nil)
	}

	code := tea.StringValue(response.Body.Code)
	message := tea.StringValue(response.Body.Message)

	// 如果发送失败
	if code != "OK" {
		errType := p.getErrorType(code)
		return &SendResponse{
			MsgID:     tea.StringValue(response.Body.BizId),
			Success:   false,
			ErrorCode: code,
			ErrorMsg:  message,
		}, NewSMSError(code, message, ShouldRetry(errType), nil)
	}

	// 如果发送的是验证码，存储到 Redis
	if codeValue, ok := req.Params["code"]; ok {
		key := getCodeKey(req.BizID, req.Phone)
		_ = p.redis.Set(ctx, key, codeValue, p.codeExpiry).Err()
	}

	// 返回成功响应
	return &SendResponse{
		MsgID:     tea.StringValue(response.Body.BizId),
		Success:   true,
		ErrorCode: "",
		ErrorMsg:  "",
	}, nil
}

// Verify 验证短信验证码
func (p *AliyunProvider) Verify(ctx context.Context, req *VerifyRequest) (*VerifyResponse, error) {
	key := getCodeKey(req.BizID, req.Phone)

	storedCode, err := p.redis.Get(ctx, key).Result()
	if err == redis.Nil {
		return &VerifyResponse{
			Success: false,
			ErrMsg:  "验证码不存在",
		}, nil
	}
	if err != nil {
		return nil, err
	}

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

// QueryStatus 查询短信发送状态（注意：阿里云需要手机号+BizId 才能查询）
func (p *AliyunProvider) QueryStatus(ctx context.Context, msgID string) (*StatusResponse, error) {
	// 阿里云的查询接口需要手机号，无法仅通过 msgID 查询
	return nil, NewSMSError(
		"NOT_SUPPORTED",
		"阿里云查询需要手机号，请使用 QueryStatusByPhone 方法",
		false,
		nil,
	)
}

// QueryStatusByPhone 通过手机号查询短信状态
func (p *AliyunProvider) QueryStatusByPhone(ctx context.Context, phone string) ([]*StatusResponse, error) {
	// 查询今天和昨天的短信记录（阿里云只支持查询最近30天的记录）
	now := time.Now()
	dates := []string{
		now.Format("20060102"),                      // 今天
		now.Add(-24 * time.Hour).Format("20060102"), // 昨天
	}

	var allResults []*StatusResponse

	for _, sendDate := range dates {
		// 创建查询请求
		queryRequest := &dysmsapi.QuerySendDetailsRequest{
			PhoneNumber: tea.String(phone),
			SendDate:    tea.String(sendDate),
			PageSize:    tea.Int64(10),
			CurrentPage: tea.Int64(1),
		}

		response, err := p.client.QuerySendDetails(queryRequest)
		if err != nil {
			return nil, p.handleQueryError(err)
		}

		// 检查响应
		if response.Body == nil {
			continue
		}

		code := tea.StringValue(response.Body.Code)
		if code != "OK" {
			// 如果查询失败，记录错误但继续查询下一个日期
			continue
		}

		// 解析结果
		if response.Body.SmsSendDetailDTOs != nil &&
			response.Body.SmsSendDetailDTOs.SmsSendDetailDTO != nil {

			for _, detail := range response.Body.SmsSendDetailDTOs.SmsSendDetailDTO {
				status := &StatusResponse{
					MsgID:       tea.StringValue(detail.OutId),
					Phone:       tea.StringValue(detail.PhoneNum),
					Status:      p.parseStatus(tea.Int64Value(detail.SendStatus)),
					SentTime:    p.parseSendDate(tea.StringValue(detail.SendDate)),
					ReceiveTime: p.parseReceiveDate(tea.StringValue(detail.ReceiveDate)),
					ErrorMsg:    tea.StringValue(detail.ErrCode),
				}
				allResults = append(allResults, status)
			}
		}
	}

	return allResults, nil
}

// getCodeKey 获取验证码存储的key
func getCodeKey(bizID, phone string) string {
	return "sms:code:" + bizID + ":" + phone
}

// ========== 辅助方法 ==========

// handleSendError 处理发送错误
func (p *AliyunProvider) handleSendError(err error) error {
	if sdkErr, ok := err.(*tea.SDKError); ok {
		code := tea.StringValue(sdkErr.Code)
		message := tea.StringValue(sdkErr.Message)

		// 尝试从 Data 中获取更多错误信息
		if sdkErr.Data != nil {
			var data map[string]interface{}
			dataStr := tea.StringValue(sdkErr.Data)
			if dataStr != "" {
				_ = json.Unmarshal([]byte(dataStr), &data)
				if recommend, ok := data["Recommend"]; ok {
					message = fmt.Sprintf("%s (建议: %v)", message, recommend)
				}
			}
		}

		errType := p.getErrorType(code)
		return NewSMSError(code, message, ShouldRetry(errType), err)
	}

	// 未知错误，可以重试
	return NewSMSError("UNKNOWN_ERROR", err.Error(), true, err)
}

// handleQueryError 处理查询错误
func (p *AliyunProvider) handleQueryError(err error) error {
	if sdkErr, ok := err.(*tea.SDKError); ok {
		code := tea.StringValue(sdkErr.Code)
		message := tea.StringValue(sdkErr.Message)
		return NewSMSError(code, message, false, err)
	}
	return NewSMSError("QUERY_ERROR", err.Error(), true, err)
}

// getErrorType 获取错误类型
func (p *AliyunProvider) getErrorType(code string) ErrorType {
	// 阿里云错误码映射
	errorMapping := map[string]ErrorType{
		"isv.BUSINESS_LIMIT_CONTROL":      ErrorTypeRateLimit,    // 业务限流
		"isv.OUT_OF_SERVICE":              ErrorTypeCircuitBreak, // 停机
		"isv.AMOUNT_NOT_ENOUGH":           ErrorTypeBalance,      // 余额不足
		"isv.INVALID_PARAMETERS":          ErrorTypeFormat,       // 参数错误
		"isv.MOBILE_NUMBER_ILLEGAL":       ErrorTypeInvalidPhone, // 手机号非法
		"isv.MOBILE_COUNT_OVER_LIMIT":     ErrorTypeFormat,       // 手机号数量超限
		"isp.RAM_PERMISSION_DENY":         ErrorTypeFormat,       // 权限不足
		"isv.TEMPLATE_MISSING_PARAMETERS": ErrorTypeFormat,       // 模板参数缺失
		"isv.SMS_TEMPLATE_ILLEGAL":        ErrorTypeFormat,       // 模板非法
		"isv.SMS_SIGNATURE_ILLEGAL":       ErrorTypeFormat,       // 签名非法
		"isv.SIGN_NAME_ILLEGAL":           ErrorTypeFormat,       // 签名名称非法
		"RequestTimeout":                  ErrorTypeTimeout,      // 请求超时
		"Throttling.User":                 ErrorTypeRateLimit,    // 用户级别限流
		"Throttling":                      ErrorTypeRateLimit,    // 限流
	}

	if errType, ok := errorMapping[code]; ok {
		return errType
	}

	// 默认为其他错误
	return ErrorTypeOther
}

// parseStatus 解析短信状态
// 阿里云状态码：1-等待回执 2-发送失败 3-发送成功
func (p *AliyunProvider) parseStatus(status int64) MessageStatus {
	switch status {
	case 1:
		return StatusPending // 等待回执
	case 2:
		return StatusFailed // 发送失败
	case 3:
		return StatusDelivered // 发送成功
	default:
		return StatusUnknown
	}
}

// parseSendDate 解析发送日期
// 格式示例：2024-12-08 10:30:00
func (p *AliyunProvider) parseSendDate(dateStr string) int64 {
	if dateStr == "" {
		return 0
	}

	// 尝试解析日期
	t, err := time.Parse("2006-01-02 15:04:05", dateStr)
	if err != nil {
		return 0
	}

	return t.Unix()
}

// parseReceiveDate 解析接收日期
func (p *AliyunProvider) parseReceiveDate(dateStr string) int64 {
	return p.parseSendDate(dateStr)
}
