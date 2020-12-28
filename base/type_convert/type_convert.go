package type_convert

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"time"
)

var (
	ErrOverflow            = errors.New("Number is overflowed")
	ErrUnderflow           = errors.New("Number is underflowed")
	ErrBadFormat           = errors.New("Input is bad format, cann't convert")
	ErrUnsupportedPlatform = fmt.Errorf("Library is not support on %d bitsize platform yet!", strconv.IntSize)
)

//将任意类型转换为int8,兼容指针
//对 int, uint 的处理仅兼容32位
func ToInt8(in interface{}) (i8 int8, err error) {
	switch val := in.(type) {
	case int8:
		return val, nil
	case int16:
		return int16ToInt8(val) //出于少转型的考虑，int16ToInt8 / int16ToInt8 / int16ToInt8 分开，实际上, int16ToInt8　可以替代所有，同理下面的 uint*
	case int32:
		return int32ToInt8(val)
	case int64:
		return int64ToInt8(val)
	case int:
		if strconv.IntSize == 32 {
			return int32ToInt8(int32(val))
		} else if strconv.IntSize == 64 {
			return int64ToInt8(int64(val))
		}
		return 0, ErrUnsupportedPlatform
	case uint8:
		return uint8ToInt8(val)
	case uint16:
		return uint16ToInt8(val)
	case uint32:
		return uint32ToInt8(val)
	case uint64:
		return uint64ToInt8(val)
	case uint:
		if strconv.IntSize == 32 {
			return uint32ToInt8(uint32(val))
		} else if strconv.IntSize == 64 {
			return uint64ToInt8(uint64(val))
		}
		return 0, ErrUnsupportedPlatform
	case float32:
		return float32ToInt8(val)
	case float64:
		return float64ToInt8(val)
	case bool:
		if val {
			i8 = 1
		}
		return
	case time.Time:
		return 0, ErrOverflow //time.Time 类型转化为 int8 注定会越界
	case string:
		i64, terr := StringToInt64(val)
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		i8, err = int64ToInt8(i64)
		if err == nil {
			err = terr
		}
		return
	case []byte:
		return bytesToInt8(val)
	default:
		i64, terr := ReflectToInt(reflect.ValueOf(in))
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		i8, err = int64ToInt8(i64)
		if err == nil {
			err = terr
		}
		return
	}
	return
}

/**
使用ToInt8将一个interface{} 数组转化成 int8 数组, 兼容指针，会丢弃一切错误
*/
func ToInt8s(in []interface{}) []int8 {
	argSize := len(in)
	retval := make([]int8, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToInt8(in)
	}
	return retval
}

//将任意类型转换为int16,兼容指针
//对 int, uint 的处理仅兼容32位
func ToInt16(in interface{}) (i16 int16, err error) {
	switch val := in.(type) {
	case int16:
		return val, nil
	case int8:
		return int16(val), nil
	case int32:
		return int32ToInt16(val)
	case int64:
		return int64ToInt16(val)
	case int:
		if strconv.IntSize == 32 {
			return int32ToInt16(int32(val))
		} else if strconv.IntSize == 64 {
			return int64ToInt16(int64(val))
		}
		return 0, ErrUnsupportedPlatform
	case uint8:
		return int16(val), nil //uint8一定在int16范围内
	case uint16:
		return uint16ToInt16(val)
	case uint32:
		return uint32ToInt16(val)
	case uint64:
		return uint64ToInt16(val)
	case uint:
		if strconv.IntSize == 32 {
			return uint32ToInt16(uint32(val))
		} else if strconv.IntSize == 64 {
			return uint64ToInt16(uint64(val))
		}
		return 0, ErrUnsupportedPlatform
	case float32:
		return float32ToInt16(val)
	case float64:
		return float64ToInt16(val)
	case bool:
		if val {
			i16 = 1
		}
		return
	case time.Time:
		return 0, ErrOverflow //time.Time 类型转化为 int16 注定会越界
	case string:
		i64, terr := StringToInt64(val)
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		i16, err = int64ToInt16(i64)
		if err == nil {
			err = terr
		}
		return
	case []byte:
		return bytesToInt16(val)
	default:
		i64, terr := ReflectToInt(reflect.ValueOf(in))
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		i16, err = int64ToInt16(i64)
		if err == nil {
			err = terr
		}
		return
	}
	return
}

/**
将 interface{} 数组转化为 int16 数组
*/
func ToInt16s(in []interface{}) []int16 {
	argSize := len(in)
	retval := make([]int16, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToInt16(in[i])
	}
	return retval
}

/**
将任意类型转化为int32
*/
func ToInt32(in interface{}) (i32 int32, err error) {
	switch val := in.(type) {
	case int8:
		return int32(val), nil
	case int16:
		return int32(val), nil
	case int32:
		return val, nil
	case int64:
		return int64ToInt32(val)
	case int:
		if strconv.IntSize == 32 {
			return int32(val), nil
		} else if strconv.IntSize == 64 {
			return int64ToInt32(int64(val))
		}
		return 0, ErrUnsupportedPlatform
	case uint8:
		return int32(val), nil
	case uint16:
		return int32(val), nil
	case uint32:
		return uint32ToInt32(val)
	case uint64:
		return uint64ToInt32(val)
	case uint:
		if strconv.IntSize == 32 {
			return uint32ToInt32(uint32(val))
		} else if strconv.IntSize == 64 {
			return uint64ToInt32(uint64(val))
		}
		return 0, ErrUnsupportedPlatform
	case float32:
		return int32(val), nil
	case float64:
		return float64ToInt32(val)
	case bool:
		if val {
			i32 = 1
		}
		return
	case time.Time:
		return int64ToInt32(val.Unix())
	case string:
		i64, terr := StringToInt64(val)
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		i32, err = int64ToInt32(i64)
		if err == nil {
			err = terr
		}
		return
	case []byte:
		return bytesToInt32(val)
	default:
		i64, terr := ReflectToInt(reflect.ValueOf(in))
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		i32, err = int64ToInt32(i64)
		if err == nil {
			err = terr
		}
		return
	}
	return
}

/**
将一个 interface{} 的数组转化为一个 int32 的数组
*/
func ToInt32s(in []interface{}) []int32 {
	argSize := len(in)
	retval := make([]int32, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToInt32(in[i])
	}
	return retval
}

/**
将任意类型转为 int64
*/
func ToInt64(in interface{}) (i64 int64, err error) {
	switch val := in.(type) {
	case int8:
		return int64(val), nil
	case int16:
		return int64(val), nil
	case int32:
		return int64(val), nil
	case int64:
		return val, nil
	case int:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return int64(val), nil
		}
		return 0, ErrUnsupportedPlatform
	case uint8:
		return int64(val), nil
	case uint16:
		return int64(val), nil
	case uint32:
		return int64(val), nil
	case uint64:
		return uint64ToInt64(val)
	case uint:
		if strconv.IntSize == 32 {
			return int64(val), nil
		} else if strconv.IntSize == 64 {
			return uint64ToInt64(uint64(val))
		}
		return 0, ErrUnsupportedPlatform
	case float32:
		return int64(val), nil
	case float64:
		return int64(val), nil
	case bool:
		if val {
			i64 = 1
		}
		return
	case time.Time:
		return val.Unix(), nil
	case string:
		return StringToInt64(val)
	case []byte:
		return bytesToInt64(val)
	default:
		return ReflectToInt(reflect.ValueOf(in))
	}
	return
}

/**
将任意类型的数组转为int64类型的数组
*/
func ToInt64s(in []interface{}) []int64 {
	argSize := len(in)
	retval := make([]int64, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToInt64(in[i])
	}
	return retval
}

/**
将任意类型转为int(int32/int64， 具体看平台)
*/
func ToInt(in interface{}) (i int, err error) {
	if strconv.IntSize == 32 {
		i32, err := ToInt32(in)
		if err != nil && err != ErrBadFormat {
			return 0, err
		}
		return int(i32), err
	} else if strconv.IntSize == 64 {
		i64, err := ToInt64(in)
		if err != nil && err != ErrBadFormat {
			return 0, err
		}
		return int(i64), err
	}
	return 0, ErrUnsupportedPlatform
}

/**
将任意类型的数组转为int 类型的数组
*/
func ToInts(in []interface{}) []int {
	argSize := len(in)
	retval := make([]int, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToInt(in[i])
	}
	return retval
}

func ToUint8(in interface{}) (ui8 uint8, err error) {
	switch val := in.(type) {
	case int8:
		return int64ToUint8(int64(val))
	case int16:
		return int64ToUint8(int64(val))
	case int32:
		return int64ToUint8(int64(val))
	case int64:
		return int64ToUint8(val)
	case int:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return int64ToUint8(int64(val))
		}
		return 0, ErrUnsupportedPlatform
	case uint8:
		return val, nil
	case uint16:
		return uint64ToUint8(uint64(val))
	case uint32:
		return uint64ToUint8(uint64(val))
	case uint64:
		return uint64ToUint8(val)
	case uint:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return uint64ToUint8(uint64(val))
		}
		return 0, ErrOverflow
	case float32:
		return float64ToUint8(float64(val))
	case float64:
		return float64ToUint8(val)
	case bool:
		if val {
			ui8 = 1
		}
		return
	case time.Time:
		return 0, ErrOverflow
	case string:
		i64, terr := StringToInt64(val)
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		ui8, err = int64ToUint8(i64)
		if err == nil {
			err = terr
		}
		return
	case []byte:
		return bytesToUint8(val)
	default:
		i64, terr := ReflectToInt(reflect.ValueOf(in))
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		ui8, err = int64ToUint8(i64)
		if err == nil {
			err = terr
		}
		return
	}
	return
}

func ToUint8s(in []interface{}) []uint8 {
	argSize := len(in)
	retval := make([]uint8, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToUint8(in[i])
	}
	return retval
}

func ToUint16(in interface{}) (u16 uint16, err error) {
	switch val := in.(type) {
	case int8:
		return int64ToUint16(int64(val))
	case int16:
		return int64ToUint16(int64(val))
	case int32:
		return int64ToUint16(int64(val))
	case int64:
		return int64ToUint16(val)
	case int:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return int64ToUint16(int64(val))
		}
		return 0, ErrUnsupportedPlatform
	case uint8:
		return uint16(val), nil
	case uint16:
		return val, nil
	case uint32:
		return uint64ToUint16(uint64(val))
	case uint64:
		return uint64ToUint16(val)
	case uint:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return uint64ToUint16(uint64(val))
		}
		return 0, ErrUnsupportedPlatform
	case float32:
		return float64ToUint16(float64(val))
	case float64:
		return float64ToUint16(val)
	case bool:
		if val {
			u16 = 1
		}
		return
	case time.Time:
		return 0, ErrOverflow
	case string:
		ui64, terr := StringToUint64(val)
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		u16, err = uint64ToUint16(ui64)
		if err == nil {
			err = terr
		}
		return
	case []byte:
		return bytesToUint16(val)
	default:
		ui64, terr := ReflectToUint(reflect.ValueOf(in))
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		u16, err = uint64ToUint16(ui64)
		if err == nil {
			err = terr
		}
		return
	}
	return
}

func ToUint16s(in []interface{}) []uint16 {
	argSize := len(in)
	retval := make([]uint16, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToUint16(in[i])
	}
	return retval
}

func ToUint32(in interface{}) (u32 uint32, err error) {
	switch val := in.(type) {
	case int8:
		return int64ToUint32(int64(val))
	case int16:
		return int64ToUint32(int64(val))
	case int32:
		return int64ToUint32(int64(val))
	case int64:
		return int64ToUint32(val)
	case int:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return int64ToUint32(int64(val))
		}
		return 0, ErrUnsupportedPlatform
	case uint8:
		return uint32(val), nil
	case uint16:
		return uint32(val), nil
	case uint32:
		return val, nil
	case uint64:
		return uint64ToUint32(uint64(val))
	case uint:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return uint64ToUint32(uint64(val))
		}
		return 0, ErrUnsupportedPlatform
	case float32:
		return float64ToUint32(float64(val))
	case float64:
		return float64ToUint32(val)
	case bool:
		if val {
			u32 = 1
		}
		return
	case time.Time:
		return int64ToUint32(val.Unix())
	case string:
		ui64, terr := StringToUint64(val)
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		u32, err = uint64ToUint32(ui64)
		if err == nil {
			err = terr
		}
		return
	case []byte:
		return bytesToUint32(val)
	default:
		ui64, terr := ReflectToUint(reflect.ValueOf(in))
		if terr != nil && terr != ErrBadFormat {
			return 0, terr
		}
		u32, err = uint64ToUint32(ui64)
		if err == nil {
			err = terr
		}
		return
	}
	return
}

func ToUint32s(in []interface{}) []uint32 {
	argSize := len(in)
	retval := make([]uint32, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToUint32(in[i])
	}
	return retval
}

func ToUint64(in interface{}) (ui64 uint64, err error) {
	switch val := in.(type) {
	case int8:
		return int64ToUint64(int64(val))
	case int16:
		return int64ToUint64(int64(val))
	case int32:
		return int64ToUint64(int64(val))
	case int64:
		return int64ToUint64(val)
	case int:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return int64ToUint64(int64(val))
		}
		return 0, ErrUnsupportedPlatform
	case uint8:
		return uint64(val), nil
	case uint16:
		return uint64(val), nil
	case uint32:
		return uint64(val), nil
	case uint64:
		return val, nil
	case uint:
		if strconv.IntSize == 32 || strconv.IntSize == 64 {
			return uint64(val), nil
		}
		return 0, ErrUnsupportedPlatform
	case float32:
		return float64ToUint64(float64(val))
	case float64:
		return float64ToUint64(val)
	case bool:
		if val {
			ui64 = 1
		}
		return
	case time.Time:
		return uint64(val.Unix()), nil
	case string:
		return StringToUint64(val)
	case []byte:
		return bytesToUint64(val)
	default:
		return ReflectToUint(reflect.ValueOf(in))
	}
	return
}

func ToUint64s(in []interface{}) []uint64 {
	argSize := len(in)
	retval := make([]uint64, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToUint64(in[i])
	}
	return retval
}

func ToUint(in interface{}) (ui uint, err error) {
	if strconv.IntSize == 32 {
		u32, err := ToUint32(in)
		if err != nil && err != ErrBadFormat {
			return 0, err
		}
		return uint(u32), err
	} else if strconv.IntSize == 64 {
		u64, err := ToUint64(in)
		if err != nil && err != ErrBadFormat {
			return 0, err
		}
		return uint(u64), err
	}
	return 0, ErrUnsupportedPlatform
}

func ToUints(in []interface{}) []uint {
	argSize := len(in)
	retval := make([]uint, argSize)
	for i := 0; i < argSize; i++ {
		retval[i], _ = ToUint(in[i])
	}
	return retval
}

func ReflectToInt(iv reflect.Value) (int64, error) {
	if iv.Kind() == reflect.Ptr {
		return ReflectToInt(iv.Elem())
	}
	switch iv.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return iv.Int(), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		uint := iv.Uint()
		if uint > math.MaxInt64 {
			return 0, ErrOverflow
		}
		return int64(uint), nil
	case reflect.Float32, reflect.Float64:
		flt := iv.Float()
		if flt < math.MinInt64 {
			return 0, ErrUnderflow
		}
		if flt > math.MaxInt64 {
			return 0, ErrOverflow
		}
		return int64(flt), nil
	case reflect.Bool:
		if iv.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.String:
		str := iv.String()
		return StringToInt64(str)
	}
	return 0, fmt.Errorf("Cann't Cast to Int from type %s", iv.Type().Name())
}

func ReflectToUint(iv reflect.Value) (uint64, error) {
	if iv.Kind() == reflect.Ptr {
		return ReflectToUint(iv.Elem())
	}
	switch iv.Kind() {
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		it := iv.Int()
		if it < 0 {
			return 0, ErrUnderflow
		}
		return uint64(it), nil
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return iv.Uint(), nil
	case reflect.Float32, reflect.Float64:
		flt := iv.Float()
		if flt < 0 {
			return 0, ErrUnderflow
		}
		return uint64(flt), nil
	case reflect.Bool:
		if iv.Bool() {
			return 1, nil
		}
		return 0, nil
	case reflect.String:
		str := iv.String()
		return StringToUint64(str)
	}
	return 0, fmt.Errorf("Cann't Cast to Int from type %s", iv.Type().Name())
}

func float64ToInt8(in float64) (int8, error) {
	if in > math.MaxInt8 {
		return 0, ErrOverflow
	}
	if in < math.MinInt8 {
		return 0, ErrUnderflow
	}
	return int8(in), nil
}

func float64ToInt16(in float64) (int16, error) {
	if in > math.MaxInt16 {
		return 0, ErrOverflow
	}
	if in < math.MinInt16 {
		return 0, ErrUnderflow
	}
	return int16(in), nil
}

func float64ToInt32(in float64) (int32, error) {
	if in > math.MaxInt32 {
		return 0, ErrOverflow
	}
	if in < math.MinInt32 {
		return 0, ErrUnderflow
	}
	return int32(in), nil
}

func float32ToInt8(in float32) (int8, error) {
	if in > math.MaxInt8 {
		return 0, ErrOverflow
	}
	if in < math.MinInt8 {
		return 0, ErrUnderflow
	}
	return int8(in), nil
}

func float32ToInt16(in float32) (int16, error) {
	if in > math.MaxInt16 {
		return 0, ErrOverflow
	}
	if in < math.MinInt16 {
		return 0, ErrUnderflow
	}
	return int16(in), nil
}

func uint64ToInt8(in uint64) (int8, error) {
	if in > math.MaxInt8 {
		return 0, ErrOverflow
	}
	return int8(in), nil
}

func uint64ToInt16(in uint64) (int16, error) {
	if in > math.MaxInt16 {
		return 0, ErrOverflow
	}
	return int16(in), nil
}

func uint64ToInt32(in uint64) (int32, error) {
	if in > math.MaxInt32 {
		return 0, ErrOverflow
	}
	return int32(in), nil
}

func uint64ToInt64(in uint64) (int64, error) {
	if in > math.MaxInt64 {
		return 0, ErrOverflow
	}
	return int64(in), nil
}

func uint32ToInt8(in uint32) (int8, error) {
	if in > math.MaxInt8 {
		return 0, ErrOverflow
	}
	return int8(in), nil
}

func uint32ToInt16(in uint32) (int16, error) {
	if in > math.MaxInt16 {
		return 0, ErrOverflow
	}
	return int16(in), nil
}

func uint32ToInt32(in uint32) (int32, error) {
	if in > math.MaxInt32 {
		return 0, ErrOverflow
	}
	return int32(in), nil
}

func uint16ToInt8(in uint16) (int8, error) {
	if in > math.MaxInt8 {
		return 0, ErrOverflow
	}
	return int8(in), nil
}

func uint16ToInt16(in uint16) (int16, error) {
	if in > math.MaxInt16 {
		return 0, ErrOverflow
	}
	return int16(in), nil
}

func uint8ToInt8(in uint8) (int8, error) {
	if in > math.MaxInt8 {
		return 0, ErrOverflow
	}
	return int8(in), nil
}

func int64ToInt8(in int64) (int8, error) {
	if in < math.MinInt8 {
		return 0, ErrUnderflow
	}
	if in > math.MaxInt8 {
		return 0, ErrOverflow
	}
	return int8(in), nil
}

func int64ToInt16(in int64) (int16, error) {
	if in < math.MinInt16 {
		return 0, ErrUnderflow
	}
	if in > math.MaxInt16 {
		return 0, ErrOverflow
	}
	return int16(in), nil
}

func int64ToInt32(in int64) (int32, error) {
	if in > math.MaxInt32 {
		return 0, ErrOverflow
	}
	if in < math.MinInt32 {
		return 0, ErrUnderflow
	}
	return int32(in), nil
}

func int32ToInt8(in int32) (int8, error) {
	if in < math.MinInt8 {
		return 0, ErrUnderflow
	}
	if in > math.MaxInt8 {
		return 0, ErrOverflow
	}
	return int8(in), nil
}

func int32ToInt16(in int32) (int16, error) {
	if in < math.MinInt16 {
		return 0, ErrUnderflow
	}
	if in > math.MaxInt16 {
		return 0, ErrOverflow
	}
	return int16(in), nil
}

func int16ToInt8(in int16) (int8, error) {
	if in < math.MinInt8 {
		return 0, ErrUnderflow
	}
	if in > math.MaxInt8 {
		return 0, ErrOverflow
	}
	return int8(in), nil
}

func int64ToUint8(in int64) (uint8, error) {
	if in < 0 {
		return 0, ErrUnderflow
	}
	if in > math.MaxUint8 {
		return 0, ErrOverflow
	}
	return uint8(in), nil
}

func int64ToUint16(in int64) (uint16, error) {
	if in < 0 {
		return 0, ErrUnderflow
	}
	if in > math.MaxUint16 {
		return 0, ErrOverflow
	}
	return uint16(in), nil
}

func int64ToUint32(in int64) (uint32, error) {
	if in < 0 {
		return 0, ErrUnderflow
	}
	if in > math.MaxUint32 {
		return 0, ErrOverflow
	}
	return uint32(in), nil
}

func int64ToUint64(in int64) (uint64, error) {
	if in < 0 {
		return 0, ErrUnderflow
	}
	return uint64(in), nil
}

func uint64ToUint8(in uint64) (uint8, error) {
	if in > math.MaxUint8 {
		return 0, ErrOverflow
	}
	return uint8(in), nil
}

func uint64ToUint16(in uint64) (uint16, error) {
	if in > math.MaxUint16 {
		return 0, ErrOverflow
	}
	return uint16(in), nil
}

func uint64ToUint32(in uint64) (uint32, error) {
	if in > math.MaxUint32 {
		return 0, ErrOverflow
	}
	return uint32(in), nil
}

func float64ToUint8(in float64) (uint8, error) {
	if in < 0 {
		return 0, ErrUnderflow
	}
	if in > math.MaxUint8 {
		return 0, ErrOverflow
	}
	return uint8(in), nil
}

func float64ToUint16(in float64) (uint16, error) {
	if in < 0 {
		return 0, ErrUnderflow
	}
	if in > math.MaxUint16 {
		return 0, ErrOverflow
	}
	return uint16(in), nil
}

func float64ToUint32(in float64) (uint32, error) {
	if in < 0 {
		return 0, ErrUnderflow
	}
	if in > math.MaxUint32 {
		return 0, ErrOverflow
	}
	return uint32(in), nil
}

func float64ToUint64(in float64) (uint64, error) {
	if in < 0 {
		return 0, ErrUnderflow
	}
	return uint64(in), nil
}
