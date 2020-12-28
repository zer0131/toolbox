package type_convert

import (
	"math"
	"reflect"
	"strconv"
)

func ToFloat32(in interface{}) (f32 float32, err error) {
	switch val := in.(type) {
	case float32:
		return val, nil
	case float64:
		return float32(val), nil
	case int8:
		return float32(val), nil
	case int16:
		return float32(val), nil
	case int32:
		return float32(val), nil
	case int64:
		return float32(val), nil
	case int:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return float32(val), nil
		}
		return 0, ErrUnsupportedPlatform
	case uint8:
		return float32(val), nil
	case uint16:
		return float32(val), nil
	case uint32:
		return float32(val), nil
	case uint64:
		return float32(val), nil
	case uint:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return float32(val), nil
		}
		return 0, ErrUnsupportedPlatform
	case bool:
		if val {
			f32 = 1.0
		}
		return
	case string:
		var f64 float64
		f64, err = StringToFloat64(val)
		if err != nil && err != ErrBadFormat { //ERROR_NUMBER_BADFORMAT 也有可能部分解析出值
			return 0, err
		}
		return float32(f64), err
	default:
		var f64 float64
		f64, err = ReflectToFloat64(reflect.ValueOf(in))
		if err != nil && err != ErrBadFormat {
			return 0, err
		}
		return float32(f64), err
	}
}

func ToFloat32s(in []interface{}) []float32 {
	argSize := len(in)
	retval := make([]float32, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToFloat32(in[i])
	}
	return retval
}

func ToFloat64(in interface{}) (f64 float64, err error) {
	switch val := in.(type) {
	case float64:
		return val, nil
	case float32:
		return float64(val), nil
	case int8:
		return float64(val), nil
	case int16:
		return float64(val), nil
	case int32:
		return float64(val), nil
	case int64:
		return float64(val), nil
	case int:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return float64(val), nil
		}
		return 0, ErrUnsupportedPlatform
	case uint8:
		return float64(val), nil
	case uint16:
		return float64(val), nil
	case uint32:
		return float64(val), nil
	case uint64:
		return float64(val), nil
	case uint:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return float64(val), nil
		}
		return 0, ErrUnsupportedPlatform
	case bool:
		if val {
			f64 = 1.0
		}
		return
	case string:
		return StringToFloat64(val)
	default:
		return ReflectToFloat64(reflect.ValueOf(in))
	}
}

func ToFloat64s(in []interface{}) []float64 {
	argSize := len(in)
	retval := make([]float64, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToFloat64(in[i])
	}
	return retval
}

func ReflectToFloat64(val reflect.Value) (f64 float64, err error) {
	if val.Kind() == reflect.Ptr {
		return ReflectToFloat64(val.Elem())
	}
	switch val.Kind() {
	case reflect.Float32, reflect.Float64:
		return val.Float(), nil
	case reflect.Int, reflect.Int16, reflect.Int8, reflect.Int32, reflect.Int64:
		return float64(val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		u64 := val.Uint()
		return float64(u64), nil
	case reflect.String:
		return StringToFloat64(val.String())
	case reflect.Bool:
		if val.Bool() {
			f64 = 1
		}
		return
	default:
		return 0, ErrBadFormat
	}
}

func StringToFloat64(str string) (f64 float64, err error) {
	f64, err = strconv.ParseFloat(str, 64)
	if err == nil {
		return
	}
	return tryParseFloat64(str)
}

func tryParseFloat64(str string) (f64 float64, err error) {
	var negative, decimal bool
	var precision int
	runes := []rune(str)
	for _, r := range runes {
		if r == '-' {
			negative = true
			continue
		}
		if r == '.' {
			if !decimal { //小数开关未开启
				decimal = true
				continue
			} else {
				err = ErrBadFormat
				break //符号位依旧要处理，跳出去处理符号
			}
		}
		if r >= '0' && r <= '9' {
			i8 := int8(r - '0')
			if decimal { //当前处于解析小数状态
				precision++
				f64 += float64(i8) / math.Pow10(precision)
			} else { //前面的整数部分
				f64 = f64*10 + float64(i8)
			}
		} else {
			err = ErrBadFormat
			break
		}
	}
	if negative {
		f64 = -f64
	}
	return
}
