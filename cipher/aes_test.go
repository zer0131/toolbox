package main

import (
	"github.com/pkg/errors"
	"testing"
)

var (
	// 1234567890abcdef1234567890abcdef
	validKeyStr = "MTIzNDU2Nzg5MGFiY2RlZjEyMzQ1Njc4OTBhYmNkZWY="

	// 1234567890abcdef
	validIVStr = "MTIzNDU2Nzg5MGFiY2RlZg=="

	// from php
	phpKeyStr = "DF61KwJDpRN3qNe5kLtncTpmrzJaQoROK7ZlxcB9jDE="
	phpIVStr  = "jDiG6tGcm3FI1RjnGNJuFQ=="
)

func Test_Test_AES256CBCEncrypt_keyLengthError(t *testing.T) {

	// decode后16位
	keyStr := "Rm5lVkt1ZGEwQmdvMGNaT0RuaU5ZQT09"

	_, err := AES256CBCEncrypt("", keyStr, "")
	if err == nil || errors.Cause(err) != ErrKeyLength {
		t.Error("err should be ErrKeyLength")
		t.SkipNow()
	}
}

func Test_AES256CBCEncrypt_validKeyAndIV(t *testing.T) {
	var tests = []struct {
		data, expect string
	}{
		{
			data:   "1234567890abcdef1234567890abcdef", // 32个字节
			expect: "9cLmNfDo/PDzc9/HFo1qqckXYLYvTa7uwve0KS/jbVG7Mx/nA38wnIY+vMzJJ5I6",
		},
		{
			data:   "1234567890abcdef1234567890abcdef123", // 32+3个字节
			expect: "9cLmNfDo/PDzc9/HFo1qqckXYLYvTa7uwve0KS/jbVGEGpex4gp8AY3ZqO9LfVJ0",
		},
		{
			data:   "1234567890abcdef", // 16个字节
			expect: "9cLmNfDo/PDzc9/HFo1qqRcsMpa0/U/QHRhx/b77UWM=",
		},
		{
			data:   "1234567890abcdef123", // 16+3个字节
			expect: "9cLmNfDo/PDzc9/HFo1qqfWIflE+CW/UbIhs0EzUuVs=",
		},
	}

	for i, tt := range tests {
		actual, err := AES256CBCEncrypt(tt.data, validKeyStr, validIVStr)
		if err != nil {
			t.Errorf("index %d err=%s", i, err)
			t.SkipNow()
		}
		if actual != tt.expect {
			t.Errorf("index %d expect %+v actual %+v", i, tt.expect, actual)
			t.SkipNow()
		}
	}
}

func Test_AES256CBCEncrypt_phpKeyAndIV(t *testing.T) {
	data := "1234567890abcdef123"
	expect := "arHQQ5abtctSggJwPpdnpTbrM1VG4j5ezCrnHjNT4EE="
	actual, _ := AES256CBCEncrypt(data, phpKeyStr, phpIVStr)
	if actual != expect {
		t.Error("actual should be same as php")
		t.SkipNow()
	}
}

func Test_AES256CBCDecrypt(t *testing.T) {
	var tests = []struct {
		data, expect string
	}{
		{
			expect: "1234567890abcdef1234567890abcdef", // 32个字节
			data:   "9cLmNfDo/PDzc9/HFo1qqckXYLYvTa7uwve0KS/jbVG7Mx/nA38wnIY+vMzJJ5I6",
		},
		{
			expect: "1234567890abcdef1234567890abcdef123", // 32+3个字节
			data:   "9cLmNfDo/PDzc9/HFo1qqckXYLYvTa7uwve0KS/jbVGEGpex4gp8AY3ZqO9LfVJ0",
		},
		{
			expect: "1234567890abcdef", // 16个字节
			data:   "9cLmNfDo/PDzc9/HFo1qqRcsMpa0/U/QHRhx/b77UWM=",
		},
		{
			expect: "1234567890abcdef123", // 16+3个字节
			data:   "9cLmNfDo/PDzc9/HFo1qqfWIflE+CW/UbIhs0EzUuVs=",
		},
	}

	for i, tt := range tests {
		actual, err := AES256CBCDecrypt(tt.data, validKeyStr, validIVStr)
		if err != nil {
			t.Errorf("index %d err=%s", i, err)
			t.SkipNow()
		}
		if actual != tt.expect {
			t.Errorf("index %d expect <%s> actual <%s>", i, tt.expect, actual)
			t.SkipNow()
		}
	}

}
