package type_convert

import "testing"

type IPTest struct {
	StringValue string
	IP          uint32
	ExpectedErr error
}

func TestStringToIPv4(t *testing.T) {
	tests := []IPTest{
		IPTest{"192.168.16.78", 3232239694, nil},
		IPTest{"147.98.78.65", 2472693313, nil},
		IPTest{"175.98.78.65", 2942455361, nil},
		IPTest{"1298.45.78.5", 0, ErrBadFormat},
		IPTest{"154.78", 0, ErrBadFormat},
	}
	for _, tst := range tests {
		tp, err := StringToIPv4(tst.StringValue)
		if tp != tst.IP || err != tst.ExpectedErr {
			t.Errorf("StringToIPv4(%+v) failed, expected result %d, but result is %d, error %+v", tst, tst.IP, tp, err)
			t.Fail()
			return
		}
		if err == nil {
			ipstr := IPv4ToString(tp)
			if ipstr != tst.StringValue {
				t.Errorf("IPv4ToString failed, expected result %s, but result is %s", tst.StringValue, ipstr)
				t.Fail()
				return
			}
		}
	}
}
