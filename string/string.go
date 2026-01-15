package g_string

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

// ToInt64 字符串转int64
func ToInt64(s string) int64 {
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

// ToInt 字符串转int
func ToInt(s string) int {
	result, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return result
}

// ToInt32 字符串转int32
func ToInt32(s string) int32 {
	result, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int32(result)
}

// ToFloat64 字符串转float64
func ToFloat64(s string) float64 {
	result, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0.0
	}
	return result
}

// Join 用分隔符拼接多个字符串
func Join(sep string, strs ...string) string {
	return strings.Join(strs, sep)
}

// TrimLeft 从左边删除n个字符 hello 2 --> llo
func TrimLeft(s string, n int) string {
	runes := []rune(s)
	if n >= len(runes) {
		return ""
	}
	return string(runes[n:])
}

// TrimRight 从右边删除n个字符 hello 2 --> hel
func TrimRight(s string, n int) string {
	runes := []rune(s)
	if n >= len(runes) {
		return ""
	}
	return string(runes[:len(runes)-n])
}

// Left 从左边取n个字符 hello 2 --> he
func Left(s string, n int) string {
	runes := []rune(s)
	if n >= len(runes) {
		return s
	}
	return string(runes[:n])
}

// Right 从右边取n个字符 hello 2 --> lo
func Right(s string, n int) string {
	runes := []rune(s)
	if n >= len(runes) {
		return s
	}
	return string(runes[len(runes)-n:])
}

// Trim 去除首尾空格
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// HasPrefix 判断前缀
func HasPrefix(s, prefix string) bool {
	return strings.HasPrefix(s, prefix)
}

// HasSuffix 判断后缀
func HasSuffix(s, suffix string) bool {
	return strings.HasSuffix(s, suffix)
}

// Contains 判断是否包含
func Contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// IsEmpty 判断是否为空或只有空格
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// ToUpper 转大写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// ToLower 转小写
func ToLower(s string) string {
	return strings.ToLower(s)
}

// Replace 替换所有匹配
func Replace(s, old, new string) string {
	return strings.ReplaceAll(s, old, new)
}

// Split 按分隔符分割，返回数组
func Split(s, sep string) []string {
	return strings.Split(s, sep)
}

// Substring 保留多少位
func Substring(s string, start, length int) string {
	runes := []rune(s)
	if start < 0 {
		start = 0
	}
	if start >= len(runes) {
		return ""
	}
	if length <= 0 {
		return string(runes[start:])
	}
	end := start + length
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end])
}

// EqualFold 不区分大小写比较
func EqualFold(s1, s2 string) bool {
	return strings.EqualFold(s1, s2)
}

// PadLeft 左边补充   "hello", '0', 10 ---> "00000hello"
func PadLeft(s string, padChar rune, totalWidth int) string {
	runes := []rune(s)
	if len(runes) >= totalWidth {
		return s
	}
	padding := make([]rune, totalWidth-len(runes))
	for i := range padding {
		padding[i] = padChar
	}
	return string(padding) + s
}

// PadRight 右边补充  "hello", '0', 10  ---> "hello00000"
func PadRight(s string, padChar rune, totalWidth int) string {
	runes := []rune(s)
	if len(runes) >= totalWidth {
		return s
	}
	padding := make([]rune, totalWidth-len(runes))
	for i := range padding {
		padding[i] = padChar
	}
	return s + string(padding)
}

// Count "hello hello hello", "hello" ---> 3
func Count(s, substr string) int {
	return strings.Count(s, substr)
}

// FirstUpper 首字母大写
func FirstUpper(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// FirstLower 首字母小写
func FirstLower(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// ToString 任意类型转字符串
func ToString(v any) string {
	if v == nil {
		return ""
	}

	// 对常见类型进行快速转换，避免反射开销
	switch val := v.(type) {
	case string:
		return val
	case int:
		return strconv.Itoa(val)
	case int64:
		return strconv.FormatInt(val, 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case uint16:
		return strconv.FormatUint(uint64(val), 10)
	case uint8:
		return strconv.FormatUint(uint64(val), 10)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case bool:
		return strconv.FormatBool(val)
	case []byte:
		return string(val)
	}

	// 使用反射处理指针和复杂类型
	val := reflect.ValueOf(v)

	// 如果是指针，递归解引用直到非指针类型或nil
	for val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return ""
		}
		val = val.Elem()
	}

	// 对解引用后的值再次尝试类型断言
	iface := val.Interface()
	switch vi := iface.(type) {
	case string:
		return vi
	case int:
		return strconv.Itoa(vi)
	case int64:
		return strconv.FormatInt(vi, 10)
	case int32:
		return strconv.FormatInt(int64(vi), 10)
	case float64:
		return strconv.FormatFloat(vi, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(vi), 'f', -1, 32)
	case bool:
		return strconv.FormatBool(vi)
	}

	// 其他复杂类型使用 fmt.Sprintf
	return fmt.Sprintf("%v", iface)
}
