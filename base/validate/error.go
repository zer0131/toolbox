package validate

import "fmt"

type BindingResult struct {
	OK          bool
	Err         error
	FieldErrors []FieldError
}

type FieldError struct {
	Field    string
	RuleName string
	Message  string
}

func newBindError(field, ruleName, message string, args ...interface{}) *FieldError {
	var formatedMsg string
	if len(args) == 0 {
		formatedMsg = message
	} else {
		formatedMsg = fmt.Sprintf(message, args...)
	}
	return &FieldError{Field: field, RuleName: ruleName, Message: formatedMsg}
}
