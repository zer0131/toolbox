package validate

import (
	"reflect"
	"testing"
)

type Test struct {
	Input   string
	Options []expression
	Error   error
}

func TestParseValidateOptions(t *testing.T) {
	tests := []Test{
		Test{`length(4, 5),required,email,maxSize(12),minSize(6)`,
			[]expression{
				expression{Name: `length`, Args: []string{`4`, `5`}},
				expression{Name: `required`, Args: nil},
				expression{Name: `email`, Args: nil},
				expression{Name: `maxSize`, Args: []string{`12`}},
				expression{Name: `minSize`, Args: []string{`6`}},
			},
			nil,
		},
		Test{`length (4, 5), required, escape('12\'dmg', 12), escape2('123 dfq', 1231)`,
			[]expression{
				expression{Name: `length`, Args: []string{`4`, `5`}},
				expression{Name: `required`, Args: nil},
				expression{Name: `escape`, Args: []string{`12'dmg`, `12`}},
				expression{Name: `escape2`, Args: []string{`123 dfq`, `1231`}},
			}, nil},
		Test{`,,space1(' 12 '),space2(  12 ), space3()`,
			[]expression{
				expression{Name: `space1`, Args: []string{` 12 `}},
				expression{Name: `space2`, Args: []string{`12`}},
				expression{Name: `space3`, Args: nil},
			}, nil},
		Test{`name(,`, nil, ErrParseFailed},
		Test{`,,,,`, nil, nil},
		Test{`(12, 45, 65)`, nil, ErrParseFailed},
		Test{`12)`, nil, ErrParseFailed},
		Test{`sbc2dbc,trim,ucfirst`, []expression{
			expression{Name: `sbc2dbc`, Args: nil},
			expression{Name: `trim`, Args: nil},
			expression{Name: `ucfirst`, Args: nil},
		}, nil},
	}
	for _, tst := range tests {
		output, err := parseExpressions(tst.Input)
		if !reflect.DeepEqual(output, tst.Options) || tst.Error != err {
			t.Errorf("ParseValidateOptions(%#v) failed, expect output %#v, real got %#v, expect error %#v, real got error %#v",
				tst.Input, tst.Options, output, tst.Error, err)
		} else {
			t.Logf("It passed ParseValidateOptions(%#v)", tst.Input)
		}
	}
}
