package xmap

import (
	"fmt"
)

func GetString(m map[string]interface{}, key string) (res string, err error) {
	val, ok := m[key]
	if !ok {
		return "", fmt.Errorf("no such key: %s", key)
	}
	res, ok1 := val.(string)
	if !ok1 {
		return "", fmt.Errorf("value of %s is not a string", key)
	}
	return res, err
}

func GetStringArray(m map[string]interface{}, key string) (res []string, err error) {
	val, ok := m[key]
	if !ok {
		return []string{}, fmt.Errorf("no such key: %s", key)
	}
	res, ok1 := val.([]string)
	if !ok1 {
		return []string{}, fmt.Errorf("value of %s is not a string array", key)
	}
	return res, err
}
