package xmap

import "testing"

func TestGetString_OK(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"key": 1,
	}

	v, err := GetString(m, "foo")
	if err != nil || v != "bar" {
		t.Fail()
	}
}

func TestGetString_NoKey(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"key": 1,
	}

	_, err := GetString(m, "ff")
	if err == nil {
		t.Fail()
	}
}

func TestGetString_TypeNotMatch(t *testing.T) {
	m := map[string]interface{}{
		"foo": "bar",
		"key": 1,
	}

	_, err := GetString(m, "key")
	if err == nil {
		t.Fail()
	}
}

func TestGetStringArray_OK(t *testing.T) {
	m := map[string]interface{}{
		"foo": []string{"123", "456"},
		"key": 1,
	}

	v, err := GetStringArray(m, "foo")
	if err != nil || len(v) != 2 || v[0] != "123" || v[1] != "456" {
		t.Fail()
	}
}

func TestGetStringArray_NoKey(t *testing.T) {
	m := map[string]interface{}{
		"foo": []string{"123", "123"},
		"key": 1,
	}

	_, err := GetStringArray(m, "ff")
	if err == nil {
		t.Fail()
	}
}

func TestGetStringArray_TypeNotMatch(t *testing.T) {
	m := map[string]interface{}{
		"foo": []string{"123", "123"},
		"key": 1,
	}

	_, err := GetStringArray(m, "key")
	if err == nil {
		t.Fail()
	}
}
