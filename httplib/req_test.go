package httplib

import "testing"

func Test_Parse(t *testing.T) {
	var testList Addr
	testList = "list://10.188.40.227:9946,10.188.40.218:9290,10.188.40.137:8888"
	m := make(map[string]struct{})
	m["http://10.188.40.227:9946"] = struct{}{}
	m["http://10.188.40.218:9290"] = struct{}{}
	m["http://10.188.40.137:8888"] = struct{}{}
	for i := 0; i < 1000; i++ {
		temp, err := testList.Parse()
		if err != nil {
			t.Error("parse error")
		}
		if _, ok := m[temp]; !ok {
			t.Error("not exists")
		}
	}
}
