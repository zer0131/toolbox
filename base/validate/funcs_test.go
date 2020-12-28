package validate

import "testing"

type Tests struct {
	Input          interface{}
	ExpectedResult bool
}

func TestFnRequired(t *testing.T) {
	tsts := []Tests{
		Tests{12, true},
		Tests{``, false},
		Tests{nil, false},
		Tests{false, true},
		Tests{`ok`, true},
		Tests{` `, true},
		Tests{make([]string, 0), false},
		Tests{make(map[string]string, 0), false},
	}
	for _, tst := range tsts {
		expv := FnRequired(tst.Input)
		if expv != tst.ExpectedResult {
			t.Errorf("FnRequired(%#+v) failed", tst.Input)
			t.Fail()
		}
	}
}

func TestFnIsEmail(t *testing.T) {
	tsts := []Tests{
		Tests{`zhangsan@sf-express.com`, true},
		Tests{`somebody@nobody.com`, true},
		Tests{`AdbUdfn@Cddl.com.cn`, true},
		Tests{` AdbUdfn@Cddl.com.cn`, false},
		Tests{`(dbUdfn@Cddl.com.cn`, false},
		Tests{`AdBTc$Udfn@Cddl.com.cn`, false},
		Tests{`*f^s71@dfy.com`, false},
	}
	for _, tst := range tsts {
		expv := FnIsEmail(tst.Input.(string))
		if expv != tst.ExpectedResult {
			t.Errorf("FnIsEmail(%#+v) failed", tst.Input)
			t.Fail()
		}
	}
}

func TestFnIsIPv6(t *testing.T) {
	tsts := []Tests{
		Tests{`2000:0000:0000:0000:0001:2345:6789:abcd`, true},
		Tests{`2000:0:0:0:1:2345:6789:abcd`, true},
		Tests{`2000::1:2345:6789:abcd`, true},
		Tests{`::1`, true},
		Tests{`fe80::fab1:56ff:fece:43f1`, true},
	}
	for _, tst := range tsts {
		expv := FnIsIPv6(tst.Input.(string))
		if expv != tst.ExpectedResult {
			t.Errorf("FnIsIPv6(%#+v) failed", tst.Input)
			t.Fail()
		}
	}
}
