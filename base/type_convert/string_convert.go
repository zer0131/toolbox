package type_convert

import (
	"fmt"
	"reflect"
	"strconv"
)

//尽可能的将一个字符串转换为int64
func StringToInt64(str string) (int int64, err error) {
	int, err = strconv.ParseInt(str, 0, 64)
	if err == nil {
		return
	}
	return tryParseInt64(str)
}

//请注意，如果，函数签名与上 StringToInt64，　不一致，当 str 为负数时会返回 ERROR_NUMBER_UNDERFLOW　的错误，以使上层ＡＰＩ统一
func StringToUint64(str string) (ui uint64, err error) {
	ui, err = strconv.ParseUint(str, 0, 64)
	if err == nil {
		return
	}
	return tryParseUint64(str)
}

/**
将一个interface 转成一个 string
*/
func ToString(d interface{}) string {
	switch v := d.(type) {
	case string:
		return v
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 64)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		return ReflectToString(reflect.ValueOf(d))
	}
}

/**
将一个 interface{} slice 转成一个 string slice
*/
func ToStrings(d []interface{}) []string {
	argSize := len(d)
	retval := make([]string, argSize)
	for i := 0; i < argSize; i++ {
		retval[i] = ToString(d[i])
	}
	return retval
}

func ReflectToString(val reflect.Value) (s string) {
	if val.Kind() == reflect.Ptr {
		return ReflectToString(val.Elem())
	}
	switch val.Kind() {
	case reflect.String:
		return val.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(val.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(val.Uint(), 10)
	case reflect.Bool:
		return strconv.FormatBool(val.Bool())
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(val.Float(), 'f', -1, 64)
	case reflect.Array:
		return fmt.Sprintf("%+v", val.Interface())
	case reflect.Struct:
		return fmt.Sprintf("%+v", val.Interface())
	default:
		return fmt.Sprintf("%s", val.Interface())
	}
}

//尽力的将一个字符串解析为数字，当遇到第一个非数字字符时停止解析
func tryParseInt64(str string) (int int64, err error) {
	runes := []rune(str)
	var nagtive bool = false
	for i, r := range runes {
		if i == 0 {
			if r == '-' {
				nagtive = true
				continue
			}
		}
		if r >= '0' && r <= '9' {
			int = int*10 + int64(r-'0')
		} else { //遇到第一个非数字字符后结束
			err = ErrBadFormat
			break
		}
	}
	if nagtive {
		int = -int
	}
	return
}

func tryParseUint64(str string) (ui uint64, err error) {
	runes := []rune(str)
	if runes[0] == '-' {
		return 0, ErrUnderflow
	}
	for _, r := range runes {
		if r >= '0' && r <= '9' {
			ui = ui*10 + uint64(r-'0')
		} else {
			err = ErrBadFormat
			break
		}
	}
	return
}
