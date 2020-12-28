package util

import (
	"strings"
	"testing"
)

type StringTest struct {
	Input string
	Ouput string
}

func TestUCFirst(t *testing.T) {
	tests := []StringTest{
		StringTest{"name is ok", "Name is ok"},
		StringTest{"我的名字", "我的名字"},
		StringTest{"123456789a", "123456789a"},
		StringTest{"a 12345678", "A 12345678"},
		StringTest{"12345678", "12345678"},
		StringTest{"", ""},
		StringTest{"$#1Opq", "$#1Opq"},
		StringTest{"_OqpNfqd", "_OqpNfqd"},
	}
	for _, ts := range tests {
		result := UCFirst(ts.Input)
		if result != ts.Ouput {
			t.Errorf("UCFirst(%s) should return %s, but real got %s \n", ts.Input, ts.Ouput, result)
		}
	}
}

func TestStr2Bytes(t *testing.T) {
	var s = "abcd.e"
	var b = []byte(s)
	tmp := Str2Bytes(s)
	for i, _ := range tmp {
		if tmp[i] != b[i] {
			t.Errorf("string to bytes error")
		}
	}
}

func TestBytes2Str(t *testing.T) {
	var s = "abcd.e"
	var b = []byte(s)
	tmp := Bytes2Str(b)
	if tmp != s {
		t.Errorf("bytes to string error")
	}
}

var s = strings.Repeat("a", 1024)

func BenchmarkTestString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs := []byte(s)
		_ = string(bs)
	}
}

func BenchmarkTestStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs := Str2Bytes(s)
		_ = Bytes2Str(bs)
	}
}
