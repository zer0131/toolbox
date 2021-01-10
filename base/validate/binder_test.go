package validate

import (
	"context"
	"net/url"
	"testing"

	"github.com/zer0131/toolbox/log"
)

type TestCase struct {
	Input  string
	Output string
}

func TestNormalizeName(t *testing.T) {
	tsts := []TestCase{
		{`i_am_ok`, `IAmOk`},
		{` i am ok`, ` i am ok`},
		{`_i-am-ok `, `IAmOk `},
		{`你是我的眼`, `你是我的眼`},
	}
	for _, tst := range tsts {
		output := normalizeName(tst.Input)
		if output != tst.Output {
			t.Errorf("NormalizeName(%#v) failed, expect %#v but real got %#v", tst.Input, tst.Output, output)
		}
	}
}

func TestBindFromValues(t *testing.T) {
	type ValueHolder struct {
		Name   string `filter:"sbc2dbc,trim" validate:"min_length(2),max_length(5)"`
		Age    uint16 `validate:"range(0, 150)"`
		Sex    string `validate:"enum('female', 'male'),required"`
		Email  string `filter:"trim" validate:"email, optional"`
		Data   []string
		NArray []int
	}
	values := url.Values{
		`name`:   []string{`　西门吹水９`},
		`age`:    []string{`32`},
		`sex`:    []string{`male`},
		`data`:   []string{`字符串1`, `字符串2`},
		`narray`: []string{`12`, `23`, `234`},
	}
	var vh ValueHolder
	result := BindFromValues(log.NewContextWithLogID(context.Background()), &vh, values)
	if result.OK {
		t.Logf("Bind OK, data is [%#v]", vh)
	} else {
		t.Errorf("Binding result is %#v", result)
	}
}

func BenchmarkBindFromValues(b *testing.B) {
	type ValueHolder struct {
		Name   string `filter:"sbc2dbc,trim" validate:"min_length(2),max_length(5)"`
		Age    uint16 `validate:"range(0, 150)"`
		Sex    string `validate:"enum('female', 'male'),required"`
		Email  string `filter:"trim" validate:"email, optional"`
		Data   []string
		NArray []int
	}
	values := url.Values{
		`name`:   []string{`　西门吹水９`},
		`age`:    []string{`32`},
		`sex`:    []string{`male`},
		`data`:   []string{`字符串1`, `字符串2`},
		`narray`: []string{`12`, `23`, `234`},
	}
	for i := 0; i < b.N; i++ {
		var vh ValueHolder
		result := BindFromValues(log.NewContextWithLogID(context.Background()), &vh, values)
		if !result.OK {
			b.Errorf("Binding result failed, result is %#v", result)
			return
		}
	}
}
