// 1. 提供两类函数：
//   - 严格类型检查函数（如 Int()）：只在值已经是目标类型时才成功
//   - 类型转换函数（如 AsInt()）：会尝试将值转换为目标类型
//
// 使用示例：
//	示例 1：从 any 类型获取 int
//	var data any = 123
//	num, err := g_any.Int(data)  // 成功，num = 123
//
//	示例 2：类型转换
//	var strData any = "456"
//	num2, err := g_any.AsInt(strData)  // 成功，num2 = 456（字符串被转换为 int）

package g_any

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// ==================== Int 系列函数 ====================

// Int 返回 int 数据 val必须是 int 类型
func Int(val any) (int, error) {
	v, ok := val.(int)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsInt 尝试将值转换为 int，支持类型转换
func AsInt(val any) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case string:
		// 如果是 string 类型，尝试解析为 int
		res, err := strconv.ParseInt(v, 10, 64)
		return int(res), err
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== Uint 系列函数 ====================

// Uint 返回 uint 数据
func Uint(val any) (uint, error) {
	v, ok := val.(uint)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsUint 尝试将值转换为 uint
func AsUint(val any) (uint, error) {
	switch v := val.(type) {
	case uint:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 64)
		return uint(res), err
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== Int8 系列函数 ====================

// Int8 返回 int8 数据
func Int8(val any) (int8, error) {
	v, ok := val.(int8)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsInt8 尝试将值转换为 int8
func AsInt8(val any) (int8, error) {
	switch v := val.(type) {
	case int8:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 8)
		return int8(res), err
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== Uint8 系列函数 ====================

// Uint8 返回 uint8 数据
func Uint8(val any) (uint8, error) {
	v, ok := val.(uint8)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsUint8 尝试将值转换为 uint8
func AsUint8(val any) (uint8, error) {
	switch v := val.(type) {
	case uint8:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 8)
		return uint8(res), err
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== Int16 系列函数 ====================

// Int16 返回 int16 数据
func Int16(val any) (int16, error) {
	v, ok := val.(int16)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsInt16 尝试将值转换为 int16
func AsInt16(val any) (int16, error) {
	switch v := val.(type) {
	case int16:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 16)
		return int16(res), err
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== Uint16 系列函数 ====================

// Uint16 返回 uint16 数据
func Uint16(val any) (uint16, error) {
	v, ok := val.(uint16)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsUint16 尝试将值转换为 uint16
func AsUint16(val any) (uint16, error) {
	switch v := val.(type) {
	case uint16:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 16)
		return uint16(res), err
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== Int32 系列函数 ====================

// Int32 返回 int32 数据
func Int32(val any) (int32, error) {
	v, ok := val.(int32)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsInt32 尝试将值转换为 int32
func AsInt32(val any) (int32, error) {
	switch v := val.(type) {
	case int32:
		return v, nil
	case string:
		res, err := strconv.ParseInt(v, 10, 32)
		return int32(res), err
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== Uint32 系列函数 ====================

// Uint32 返回 uint32 数据
func Uint32(val any) (uint32, error) {
	v, ok := val.(uint32)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsUint32 尝试将值转换为 uint32
func AsUint32(val any) (uint32, error) {
	switch v := val.(type) {
	case uint32:
		return v, nil
	case string:
		res, err := strconv.ParseUint(v, 10, 32)
		return uint32(res), err
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== Int64 系列函数 ====================

// Int64 返回 int64 数据
func Int64(val any) (int64, error) {
	v, ok := val.(int64)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsInt64 尝试将值转换为 int64
func AsInt64(val any) (int64, error) {
	switch v := val.(type) {
	case int64:
		return v, nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== Uint64 系列函数 ====================

// Uint64 返回 uint64 数据
func Uint64(val any) (uint64, error) {
	v, ok := val.(uint64)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsUint64 尝试将值转换为 uint64
func AsUint64(val any) (uint64, error) {
	switch v := val.(type) {
	case uint64:
		return v, nil
	case string:
		return strconv.ParseUint(v, 10, 64)
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== Float32 系列函数 ====================

// Float32 返回 float32 数据
func Float32(val any) (float32, error) {
	v, ok := val.(float32)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsFloat32 尝试将值转换为 float32
func AsFloat32(val any) (float32, error) {
	switch v := val.(type) {
	case float32:
		return v, nil
	case string:
		res, err := strconv.ParseFloat(v, 32)
		return float32(res), err
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== Float64 系列函数 ====================

// Float64 返回 float64 数据
func Float64(val any) (float64, error) {
	v, ok := val.(float64)
	if !ok {
		return 0, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsFloat64 尝试将值转换为 float64
func AsFloat64(val any) (float64, error) {
	switch v := val.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	}
	return 0, fmt.Errorf("参数错误:%v", val)
}

// ==================== String 系列函数 ====================

// String 返回 string 数据
func String(val any) (string, error) {
	v, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsString 将各种类型的值转换为字符串
func AsString(v any) string {
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

// ==================== Bytes 系列函数 ====================

// Bytes 返回 []byte 数据
func Bytes(val any) ([]byte, error) {
	v, ok := val.([]byte)
	if !ok {
		return nil, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// AsBytes 尝试将值转换为 []byte
func AsBytes(val any) ([]byte, error) {
	switch v := val.(type) {
	case []byte:
		return v, nil
	case string:
		return []byte(v), nil
	}

	return []byte{}, fmt.Errorf("参数错误:%v", val)
}

// ==================== Bool 系列函数 ====================

// Bool 返回 bool 数据
func Bool(val any) (bool, error) {
	v, ok := val.(bool)
	if !ok {
		return false, fmt.Errorf("参数错误:%v", val)
	}
	return v, nil
}

// ==================== JSON 相关函数 ====================

// JSONScan 将 val 转化为一个对象
// 例如：var user User; JSONScan(jsonData, &user)
func JSONScan(val any, dest any) error {
	data, err := AsBytes(val)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}
