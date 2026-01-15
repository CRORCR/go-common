package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToPtr(t *testing.T) {
	str := "hello"
	_str := ToPtr[string](str)
	assert.Equal(t, &str, _str)
}

func TestDeref(t *testing.T) {
	str := "hello"
	_str := ToPtr[string](str)
	result := Deref(_str)
	assert.Equal(t, str, result)

	var nilPtr *string
	nilResult := Deref(nilPtr)
	assert.Equal(t, "", nilResult)

	num := 42
	_num := ToPtr[int](num)
	numResult := Deref(_num)
	assert.Equal(t, num, numResult)

	var nilNumPtr *int
	nilNumResult := Deref(nilNumPtr)
	assert.Equal(t, 0, nilNumResult)
}
