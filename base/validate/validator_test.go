package validate

import (
	"context"
	"testing"

	"toolbox/log"
)

func TestCreateValidators(t *testing.T) {
	input := `enum('male', 'female'),optional`
	optons, err := parseExpressions(input)
	if err != nil {
		t.Error(err)
		return
	}
	CreateValidators(log.NewContextWithLogID(context.Background()), optons)
}
