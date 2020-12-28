package layer

import (
	"testing"
)

func Test_Register(t *testing.T) {
	RegisterModel("foo", "bar")
	RegisterModelWrapper("foo", "bar")
	RegisterService("foo", "bar")
	RegisterServiceWrapper("foo", "bar")

	ModelList()
	ModelWrapperList()
	ServiceList()
	ServiceWrapperList()
}
