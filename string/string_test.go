package g_string

import (
	"testing"
)

func TestToInt64(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"123", 123},
		{"-456", -456},
		{"0", 0},
		{"invalid", 0},
		{"", 0},
		{"9223372036854775807", 9223372036854775807},
	}

	for _, tt := range tests {
		result := ToInt64(tt.input)
		if result != tt.expected {
			t.Errorf("ToInt64(%q) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"123", 123},
		{"-456", -456},
		{"0", 0},
		{"invalid", 0},
		{"", 0},
	}

	for _, tt := range tests {
		result := ToInt(tt.input)
		if result != tt.expected {
			t.Errorf("ToInt(%q) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func TestToInt32(t *testing.T) {
	tests := []struct {
		input    string
		expected int32
	}{
		{"123", 123},
		{"-456", -456},
		{"0", 0},
		{"invalid", 0},
		{"", 0},
		{"2147483647", 2147483647},
	}

	for _, tt := range tests {
		result := ToInt32(tt.input)
		if result != tt.expected {
			t.Errorf("ToInt32(%q) = %d; want %d", tt.input, result, tt.expected)
		}
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"123.45", 123.45},
		{"-456.78", -456.78},
		{"0", 0.0},
		{"invalid", 0.0},
		{"", 0.0},
		{"3.14159", 3.14159},
	}

	for _, tt := range tests {
		result := ToFloat64(tt.input)
		if result != tt.expected {
			t.Errorf("ToFloat64(%q) = %f; want %f", tt.input, result, tt.expected)
		}
	}
}

func TestJoin(t *testing.T) {
	tests := []struct {
		sep      string
		strs     []string
		expected string
	}{
		{",", []string{"a", "b", "c"}, "a,b,c"},
		{"-", []string{"hello", "world"}, "hello-world"},
		{" ", []string{"one"}, "one"},
		{"", []string{"one", "two"}, "onetwo"},
		{",", []string{}, ""},
	}

	for _, tt := range tests {
		result := Join(tt.sep, tt.strs...)
		if result != tt.expected {
			t.Errorf("Join(%q, %v) = %q; want %q", tt.sep, tt.strs, result, tt.expected)
		}
	}
}

func TestTrimLeft(t *testing.T) {
	tests := []struct {
		input    string
		n        int
		expected string
	}{
		{"hello", 2, "llo"},
		{"world", 0, "world"},
		{"test", 10, ""},
		{"你好世界", 2, "世界"},
		{"", 1, ""},
	}

	for _, tt := range tests {
		result := TrimLeft(tt.input, tt.n)
		if result != tt.expected {
			t.Errorf("TrimLeft(%q, %d) = %q; want %q", tt.input, tt.n, result, tt.expected)
		}
	}
}

func TestTrimRight(t *testing.T) {
	tests := []struct {
		input    string
		n        int
		expected string
	}{
		{"hello", 2, "hel"},
		{"world", 0, "world"},
		{"test", 10, ""},
		{"你好世界", 2, "你好"},
		{"", 1, ""},
	}

	for _, tt := range tests {
		result := TrimRight(tt.input, tt.n)
		if result != tt.expected {
			t.Errorf("TrimRight(%q, %d) = %q; want %q", tt.input, tt.n, result, tt.expected)
		}
	}
}

func TestLeft(t *testing.T) {
	tests := []struct {
		input    string
		n        int
		expected string
	}{
		{"hello", 2, "he"},
		{"world", 0, ""},
		{"test", 10, "test"},
		{"你好世界", 2, "你好"},
		{"", 1, ""},
	}

	for _, tt := range tests {
		result := Left(tt.input, tt.n)
		if result != tt.expected {
			t.Errorf("Left(%q, %d) = %q; want %q", tt.input, tt.n, result, tt.expected)
		}
	}
}

func TestRight(t *testing.T) {
	tests := []struct {
		input    string
		n        int
		expected string
	}{
		{"hello", 2, "lo"},
		{"world", 0, ""},
		{"test", 10, "test"},
		{"你好世界", 2, "世界"},
		{"", 1, ""},
	}

	for _, tt := range tests {
		result := Right(tt.input, tt.n)
		if result != tt.expected {
			t.Errorf("Right(%q, %d) = %q; want %q", tt.input, tt.n, result, tt.expected)
		}
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"  hello  ", "hello"},
		{"\t\nworld\t\n", "world"},
		{"test", "test"},
		{"   ", ""},
		{"", ""},
	}

	for _, tt := range tests {
		result := Trim(tt.input)
		if result != tt.expected {
			t.Errorf("Trim(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestHasPrefix(t *testing.T) {
	tests := []struct {
		s        string
		prefix   string
		expected bool
	}{
		{"hello world", "hello", true},
		{"hello world", "world", false},
		{"test", "test", true},
		{"test", "testing", false},
		{"", "", true},
	}

	for _, tt := range tests {
		result := HasPrefix(tt.s, tt.prefix)
		if result != tt.expected {
			t.Errorf("HasPrefix(%q, %q) = %v; want %v", tt.s, tt.prefix, result, tt.expected)
		}
	}
}

func TestHasSuffix(t *testing.T) {
	tests := []struct {
		s        string
		suffix   string
		expected bool
	}{
		{"hello world", "world", true},
		{"hello world", "hello", false},
		{"test", "test", true},
		{"test", "testing", false},
		{"", "", true},
	}

	for _, tt := range tests {
		result := HasSuffix(tt.s, tt.suffix)
		if result != tt.expected {
			t.Errorf("HasSuffix(%q, %q) = %v; want %v", tt.s, tt.suffix, result, tt.expected)
		}
	}
}

func TestContains(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected bool
	}{
		{"hello world", "lo wo", true},
		{"hello world", "xyz", false},
		{"test", "test", true},
		{"test", "", true},
		{"", "", true},
	}

	for _, tt := range tests {
		result := Contains(tt.s, tt.substr)
		if result != tt.expected {
			t.Errorf("Contains(%q, %q) = %v; want %v", tt.s, tt.substr, result, tt.expected)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"", true},
		{"   ", true},
		{"\t\n", true},
		{"hello", false},
		{" hello ", false},
	}

	for _, tt := range tests {
		result := IsEmpty(tt.input)
		if result != tt.expected {
			t.Errorf("IsEmpty(%q) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}

func TestToUpper(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "HELLO"},
		{"World", "WORLD"},
		{"ALREADY", "ALREADY"},
		{"", ""},
	}

	for _, tt := range tests {
		result := ToUpper(tt.input)
		if result != tt.expected {
			t.Errorf("ToUpper(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestToLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"HELLO", "hello"},
		{"World", "world"},
		{"already", "already"},
		{"", ""},
	}

	for _, tt := range tests {
		result := ToLower(tt.input)
		if result != tt.expected {
			t.Errorf("ToLower(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestReplace(t *testing.T) {
	tests := []struct {
		s        string
		old      string
		new      string
		expected string
	}{
		{"hello world", "world", "Go", "hello Go"},
		{"test test test", "test", "demo", "demo demo demo"},
		{"no match", "xyz", "abc", "no match"},
		{"", "a", "b", ""},
	}

	for _, tt := range tests {
		result := Replace(tt.s, tt.old, tt.new)
		if result != tt.expected {
			t.Errorf("Replace(%q, %q, %q) = %q; want %q", tt.s, tt.old, tt.new, result, tt.expected)
		}
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		s        string
		sep      string
		expected []string
	}{
		{"a,b,c", ",", []string{"a", "b", "c"}},
		{"hello-world", "-", []string{"hello", "world"}},
		{"single", ",", []string{"single"}},
		{"", ",", []string{""}},
	}

	for _, tt := range tests {
		result := Split(tt.s, tt.sep)
		if len(result) != len(tt.expected) {
			t.Errorf("Split(%q, %q) length = %d; want %d", tt.s, tt.sep, len(result), len(tt.expected))
			continue
		}
		for i := range result {
			if result[i] != tt.expected[i] {
				t.Errorf("Split(%q, %q)[%d] = %q; want %q", tt.s, tt.sep, i, result[i], tt.expected[i])
			}
		}
	}
}

func TestSubstring(t *testing.T) {
	tests := []struct {
		input    string
		start    int
		length   int
		expected string
	}{
		{"hello world", 0, 5, "hello"},
		{"hello world", 6, 5, "world"},
		{"hello world", 0, 100, "hello world"},
		{"hello world", 6, 0, "world"},
		{"hello world", -1, 5, "hello"},
		{"hello world", 100, 5, ""},
		{"你好世界", 0, 2, "你好"},
		{"", 0, 5, ""},
	}

	for _, tt := range tests {
		result := Substring(tt.input, tt.start, tt.length)
		if result != tt.expected {
			t.Errorf("Substring(%q, %d, %d) = %q; want %q", tt.input, tt.start, tt.length, result, tt.expected)
		}
	}
}

func TestPadLeft(t *testing.T) {
	tests := []struct {
		input      string
		padChar    rune
		totalWidth int
		expected   string
	}{
		{"hello", '0', 10, "00000hello"},
		{"test", '*', 8, "****test"},
		{"long", ' ', 4, "long"},
		{"", 'a', 5, "aaaaa"},
	}

	for _, tt := range tests {
		result := PadLeft(tt.input, tt.padChar, tt.totalWidth)
		if result != tt.expected {
			t.Errorf("PadLeft(%q, %c, %d) = %q; want %q", tt.input, tt.padChar, tt.totalWidth, result, tt.expected)
		}
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		input      string
		padChar    rune
		totalWidth int
		expected   string
	}{
		{"hello", '0', 10, "hello00000"},
		{"test", '*', 8, "test****"},
		{"long", ' ', 4, "long"},
		{"", 'a', 5, "aaaaa"},
	}

	for _, tt := range tests {
		result := PadRight(tt.input, tt.padChar, tt.totalWidth)
		if result != tt.expected {
			t.Errorf("PadRight(%q, %c, %d) = %q; want %q", tt.input, tt.padChar, tt.totalWidth, result, tt.expected)
		}
	}
}

func TestCount(t *testing.T) {
	tests := []struct {
		s        string
		substr   string
		expected int
	}{
		{"hello hello hello", "hello", 3},
		{"test", "t", 2},
		{"abc", "d", 0},
		{"", "", 1},
		{"aaa", "aa", 1},
	}

	for _, tt := range tests {
		result := Count(tt.s, tt.substr)
		if result != tt.expected {
			t.Errorf("Count(%q, %q) = %d; want %d", tt.s, tt.substr, result, tt.expected)
		}
	}
}

func TestEqualFold(t *testing.T) {
	tests := []struct {
		s1       string
		s2       string
		expected bool
	}{
		{"Hello", "hello", true},
		{"WORLD", "world", true},
		{"test", "TEST", true},
		{"hello", "world", false},
	}

	for _, tt := range tests {
		result := EqualFold(tt.s1, tt.s2)
		if result != tt.expected {
			t.Errorf("EqualFold(%q, %q) = %v; want %v", tt.s1, tt.s2, result, tt.expected)
		}
	}
}

func TestFirstUpper(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "Hello"},
		{"HELLO", "HELLO"},
		{"hELLO", "HELLO"},
		{"", ""},
		{"a", "A"},
	}

	for _, tt := range tests {
		result := FirstUpper(tt.input)
		if result != tt.expected {
			t.Errorf("FirstUpper(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFirstLower(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"Hello", "hello"},
		{"HELLO", "hELLO"},
		{"hello", "hello"},
		{"", ""},
		{"A", "a"},
	}

	for _, tt := range tests {
		result := FirstLower(tt.input)
		if result != tt.expected {
			t.Errorf("FirstLower(%q) = %q; want %q", tt.input, result, tt.expected)
		}
	}
}

func TestToString(t *testing.T) {
	type TestStruct struct {
		Name string
		Age  int
	}

	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"int", 123, "123"},
		{"int64", int64(456), "456"},
		{"float64", 3.14, "3.14"},
		{"bool true", true, "true"},
		{"bool false", false, "false"},
		{"string", "hello", "hello"},
		{"nil", nil, ""},
		{"struct", TestStruct{Name: "张三", Age: 25}, "{张三 25}"},
		{"pointer to int", func() *int { i := 42; return &i }(), "42"},
		{"slice", []int{1, 2, 3}, "[1 2 3]"},
		{"map", map[string]int{"a": 1, "b": 2}, "map[a:1 b:2]"},
		{"empty string", "", ""},
		{"rune", '中', "20013"},
		{"byte", byte(65), "65"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToString(tt.input)
			// 对于 map 类型，顺序可能不确定，需要特殊处理
			if tt.name == "map" {
				if len(result) == 0 {
					t.Errorf("ToString(%v) = %q; want non-empty map string", tt.input, result)
				}
			} else {
				if result != tt.expected {
					t.Errorf("ToString(%v) = %q; want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}
