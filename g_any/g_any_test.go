package g_any

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ==================== Int 系列测试 ====================

func TestInt(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want int
		err  error
	}{
		{
			name: "正常情况：int类型",
			val:  int(123),
			want: int(123),
			err:  nil,
		},
		{
			name: "错误情况：类型不匹配",
			val:  "123",
			err:  fmt.Errorf("参数错误:%v", "123"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Int(tt.val)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAsInt(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want int
		err  error
	}{
		{
			name: "正常情况：string转int",
			val:  "123",
			want: 123,
		},
		{
			name: "正常情况：int类型",
			val:  int(456),
			want: 456,
		},
		{
			name: "错误情况：无法转换的类型",
			val:  []int{1},
			err:  fmt.Errorf("参数错误:%v", []int{1}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := AsInt(tt.val)
			assert.Equal(t, tt.want, val)
			if tt.err != nil {
				assert.Equal(t, tt.err, err)
			}
		})
	}
}

// ==================== String 系列测试 ====================

func TestString(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want string
		err  error
	}{
		{
			name: "正常情况：string类型",
			val:  "hello",
			want: "hello",
			err:  nil,
		},
		{
			name: "错误情况：类型不匹配",
			val:  123,
			err:  fmt.Errorf("参数错误:%v", 123),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := String(tt.val)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAsString(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want string
		err  error
	}{
		{
			name: "正常情况：string类型",
			val:  "hello g_any",
			want: "hello g_any",
		},
		{
			name: "正常情况：uint转换",
			val:  uint16(1231),
			want: "1231",
		},
		{
			name: "正常情况：int转换",
			val:  1,
			want: "1",
		},
		{
			name: "正常情况：float转换",
			val:  1e2,
			want: "100",
		},
		{
			name: "正常情况：[]byte转换",
			val:  []byte{72, 101, 108, 108, 111, 44, 32, 87, 111, 114, 108, 100, 33},
			want: "Hello, World!",
		},
		{
			name: "错误情况：不支持的类型",
			val: map[string]any{
				"test": 1,
				"hhh":  "sss",
			},
			want: "map[hhh:sss test:1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := AsString(tt.val)
			assert.Equal(t, tt.want, val)
		})
	}
}

// ==================== Bytes 系列测试 ====================

func TestBytes(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want []byte
		err  error
	}{
		{
			name: "正常情况：[]byte类型",
			val:  []byte("hello"),
			want: []byte("hello"),
			err:  nil,
		},
		{
			name: "错误情况：类型不匹配",
			val:  "hello",
			err:  fmt.Errorf("参数错误:%v", "hello"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Bytes(tt.val)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAsBytes(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want []byte
		err  error
	}{
		{
			name: "正常情况：string转[]byte",
			val:  "hello",
			want: []byte("hello"),
		},
		{
			name: "正常情况：[]byte类型",
			val:  []byte{1, 2, 3},
			want: []byte{1, 2, 3},
		},
		{
			name: "错误情况：无法转换的类型",
			val:  []int{1},
			want: []byte{},
			err:  fmt.Errorf("参数错误:%v", []byte{}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := AsBytes(tt.val)
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}

// ==================== Bool 系列测试 ====================

func TestBool(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want bool
		err  error
	}{
		{
			name: "正常情况：bool类型",
			val:  true,
			want: true,
			err:  nil,
		},
		{
			name: "错误情况：类型不匹配",
			val:  1,
			err:  fmt.Errorf("参数错误:%v", 1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Bool(tt.val)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

// ==================== JSON 系列测试 ====================

func TestJSONScan(t *testing.T) {
	type User struct {
		Name string `json:"name"`
	}

	testCases := []struct {
		name string
		val  any

		wantUser User
		wantErr  error
	}{
		{
			name: "正常情况：JSON字符串",
			val:  `{"name": "Tom"}`,
			wantUser: User{
				Name: "Tom",
			},
		},
		{
			name: "正常情况：[]byte",
			val:  []byte(`{"name": "Jerry"}`),
			wantUser: User{
				Name: "Jerry",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var u User
			err := JSONScan(tc.val, &u)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)
		})
	}
}

// ==================== 其他数字类型测试（示例）====================

func TestInt64(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want int64
		err  error
	}{
		{
			name: "正常情况：int64类型",
			val:  int64(123456789),
			want: int64(123456789),
			err:  nil,
		},
		{
			name: "错误情况：类型不匹配",
			val:  "123",
			err:  fmt.Errorf("参数错误:%v", "123"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Int64(tt.val)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAsInt64(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want int64
		err  error
	}{
		{
			name: "正常情况：string转int64",
			val:  "123456789",
			want: 123456789,
		},
		{
			name: "正常情况：int64类型",
			val:  int64(987654321),
			want: 987654321,
		},
		{
			name: "错误情况：无法转换的类型",
			val:  []int{1},
			err:  fmt.Errorf("参数错误:%v", []int{1}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := AsInt64(tt.val)
			assert.Equal(t, tt.want, val)
			if tt.err != nil {
				assert.Equal(t, tt.err, err)
			}
		})
	}
}

func TestFloat64(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want float64
		err  error
	}{
		{
			name: "正常情况：float64类型",
			val:  float64(123.456),
			want: float64(123.456),
			err:  nil,
		},
		{
			name: "错误情况：类型不匹配",
			val:  "123.456",
			err:  fmt.Errorf("参数错误:%v", "123.456"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Float64(tt.val)
			assert.Equal(t, tt.err, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestAsFloat64(t *testing.T) {
	tests := []struct {
		name string
		val  any
		want float64
		err  error
	}{
		{
			name: "正常情况：string转float64",
			val:  "100.0000000000",
			want: 1e2,
		},
		{
			name: "正常情况：float64类型",
			val:  float64(2.44),
			want: 2.44,
		},
		{
			name: "错误情况：无法转换的类型",
			val:  []int{1},
			err:  fmt.Errorf("参数错误:%v", []int{1}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := AsFloat64(tt.val)
			assert.Equal(t, tt.want, val)
			assert.Equal(t, tt.err, err)
		})
	}
}
