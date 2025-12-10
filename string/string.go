package g_string

import (
	"math/rand"
	"time"
)

// GenerateCode 生成几位随机数
func GenerateCode(codeLength int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "123456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = chars[r.Intn(len(chars))]
	}
	return string(code)
}

// GenerateIntCode 生成数字验证码
func GenerateIntCode(codeLength int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	if codeLength <= 0 {
		codeLength = 6
	}

	chars := "0123456789"
	code := make([]byte, codeLength)
	for i := range code {
		code[i] = chars[r.Intn(len(chars))]
	}
	return string(code)
}
