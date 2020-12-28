package util

import (
	"errors"
	"fmt"
	"runtime"
)

func NewError(format string, a ...interface{}) error {

	err := fmt.Sprintf(format, a...)
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return errors.New(err)
	}

	function := runtime.FuncForPC(pc).Name()
	msg := fmt.Sprintf("%s func:%s file:%s line:%d",
		err, function, file, line)
	return errors.New(msg)
}
