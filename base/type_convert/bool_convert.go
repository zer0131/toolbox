package type_convert

import (
	"reflect"
	"strings"
)

/**
任意类型转化为 bool 类型
*/
func ToBool(in interface{}) (bool, error) {
	switch val := in.(type) {
	case bool:
		return val, nil
	case int8:
		return val != 0, nil
	case int16:
		return val != 0, nil
	case int32:
		return val != 0, nil
	case int64:
		return val != 0, nil
	case int:
		return val != 0, nil
	case uint8:
		return val > 0, nil
	case uint16:
		return val > 0, nil
	case uint32:
		return val > 0, nil
	case uint64:
		return val > 0, nil
	case uint:
		return val > 0, nil
	case []byte:
		return false, ErrBadFormat
	case string:
		return StringToBool(val)
	default:
		return ReflectToBool(reflect.ValueOf(in))
	}
}

/**
将一个 interface slice 转化成 bool slice
*/
func ToBools(d []interface{}) []bool {
	argSize := len(d)
	retval := make([]bool, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToBool(d[i])
	}
	return retval
}

func ReflectToBool(val reflect.Value) (bool, error) {
	if val.Kind() == reflect.Ptr {
		return ReflectToBool(val.Elem())
	}
	switch val.Kind() {
	case reflect.Bool:
		return val.Bool(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() != 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() > 0, nil
	case reflect.String:
		return StringToBool(val.String())
	default:
		return false, ErrBadFormat
	}
}

func StringToBool(str string) (bool, error) {
	if str == "true" || str == "yes" || str == "ok" || str == "y" || str == "1" {
		return true, nil
	}
	str = strings.ToLower(str)
	if str == "true" || str == "yes" || str == "ok" || str == "y" {
		return true, nil
	}
	if str == "no" || str == "false" || str == "0" || str == "n" {
		return false, nil
	}
	return false, ErrBadFormat
}
