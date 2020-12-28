package stat

import (
	"fmt"
	"testing"
)

func Test_GetRawPath(t *testing.T) {

	fmt.Println(GetRawPath("/foo/bar?a=b"))
}
