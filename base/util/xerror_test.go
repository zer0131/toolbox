package util

import (
	"fmt"
	"testing"
)

func TestError(t *testing.T) {
	err := NewError("this is a test,%d", 100)

	str := fmt.Sprintf("%v", err)
	fmt.Println(str)
}
