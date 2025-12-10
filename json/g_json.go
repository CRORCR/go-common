package g_json

import (
	jsoniter "github.com/json-iterator/go"
)

// Marshal 序列化
func Marshal(data interface{}) ([]byte, error) {
	return jsoniter.Marshal(data)
}

// MarshalToString 序列化成字符串
func MarshalToString(data interface{}) (string, error) {
	return jsoniter.MarshalToString(data)
}

// UnMarshal 反序列化
func UnMarshal(data []byte, v interface{}) error {
	return jsoniter.Unmarshal(data, v)
}

// UnmarshalFromString 字符串反序列化
func UnmarshalFromString(data string, v interface{}) error {
	return jsoniter.UnmarshalFromString(data, v)
}
