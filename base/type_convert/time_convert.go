package type_convert

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	timeType = reflect.TypeOf(time.Time{})

	ERROR_TIME_NOTVALID = errors.New("It's is not a valid time")

	DEFAULT_TIME_LAYOUT = "2006-01-02 15:04:05"

	DefaultLocation = time.FixedZone("Asia/Shanghai", 3600*8)
)

func ToTime(in interface{}) (t time.Time, err error) {
	switch val := in.(type) {
	case *time.Time:
		return *val, nil
	case time.Time:
		return val, nil
	case string:
		return StringToTime(val, DefaultLocation)
	case int8:
		return int64ToTime(int64(val)), nil
	case int16:
		return int64ToTime(int64(val)), nil
	case int32:
		return int64ToTime(int64(val)), nil
	case int64:
		return int64ToTime(int64(val)), nil
	case uint8:
		return int64ToTime(int64(val)), nil
	case uint16:
		return int64ToTime(int64(val)), nil
	case uint32:
		return int64ToTime(int64(val)), nil
	case uint64:
		i64 := int64(val)
		if i64 < 0 { //uint64转int64可能上溢
			err = ErrOverflow
			return
		}
		return int64ToTime(int64(val)), nil
	case []byte:
		return bytesToTime(val)
	default:
		return ReflectToTime(reflect.ValueOf(in))
	}
	return
}

func ReflectToTime(val reflect.Value) (t time.Time, err error) {
	if val.Kind() == reflect.Ptr {
		if val.Type() == reflect.PtrTo(timeType) { // 这本来就是一个 *time.Time
			v := val.Interface()
			tv, _ := v.(*time.Time)
			return *tv, nil
		}
		return ReflectToTime(val.Elem())
	}
	if val.Type() == timeType {
		v := val.Interface()
		tv, _ := v.(time.Time)
		return tv, nil
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int64ToTime(val.Int()), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i64 := int64(val.Uint())
		if i64 < 0 {
			err = ErrOverflow
			return
		}
		return int64ToTime(i64), nil
	case reflect.String:
		str := val.String()
		return StringToTime(str, DefaultLocation)
	default:
		err = fmt.Errorf("Input Data %#v cann't cast to time.Time, with type %s \n", val.Interface(), val.Type().Name())
		return
	}
	return
}

func int64ToTime(i64 int64) (t time.Time) {
	return time.Unix(i64, 0)
}

func bytesToTime(bs []byte) (t time.Time, err error) {
	err = t.UnmarshalBinary(bs)
	return
}

func StringToTime(str string, loc *time.Location) (t time.Time, err error) {
	t, err = time.ParseInLocation(DEFAULT_TIME_LAYOUT, str, loc)
	if err == nil {
		return
	}
	var year, month, day, hour, minute, secs int = 0, 0, 0, 0, 0, 0
	segs := strings.SplitN(str, " ", 2)
	//处理年月日
	dayPart := strings.TrimSpace(segs[0])
	if strings.Contains(dayPart, "-") { //yyyy-MM-dd style
		parts := strings.SplitN(dayPart, "-", 3)
		year = partToInt(parts[0])
		month = partToInt(parts[1])
		day = partToInt(parts[2])
		if !checkDate(year, month, day) {
			err = ERROR_TIME_NOTVALID
			return
		}
	} else if strings.Contains(dayPart, "/") { //形如 yyyy/MM/dd 或者 MM/dd/yyyy 形式
		parts := strings.SplitN(dayPart, "/", 3)
		year = partToInt(parts[0])
		month = partToInt(parts[1])
		day = partToInt(parts[2])
		if year < 100 && day > 100 { //是形如 MM/dd/yyyy 形式的，则颠倒 year,day
			year, month, day = day, year, month
		}
		if !checkDate(year, month, day) {
			err = ERROR_TIME_NOTVALID
			return
		}
	} else {
		//TODO:其它兼容格式
		err = ERROR_TIME_NOTVALID
		return
	}
	if len(segs) > 1 {
		timePart := strings.TrimSpace(segs[1])
		parts := strings.SplitN(timePart, ":", 3)
		hour = partToInt(parts[0])
		partLen := len(parts)
		if partLen > 1 {
			minute = partToInt(parts[1])
		}
		if partLen > 2 {
			secs = partToInt(parts[2])
		}
		if !checkTime(hour, minute, secs) {
			err = ERROR_TIME_NOTVALID
			return
		}
	}
	t = time.Date(year, time.Month(month), day, hour, minute, secs, 0, DefaultLocation)
	return t, nil
}

func partToInt(str string) int {
	i32, _ := strconv.ParseInt(strings.TrimLeft(strings.TrimSpace(str), "0"), 10, 32)
	return int(i32)
}

func checkTime(hour, minute, sec int) bool {
	if hour > 23 || minute > 59 {
		return false
	}
	if hour == 23 && minute == 59 {
		return sec <= 61 //允许闰秒，闰秒一般发生在23点５９分
	}
	return sec <= 60
}

func checkDate(year, month, day int) bool {
	if month < 1 || month > 12 {
		return false
	}
	if day < 1 || day > 31 {
		return false
	}
	if month == 2 {
		if isLeapYear(year) {
			if day > 29 {
				return false
			}
		}
		return day <= 28
	} else if month == 4 || month == 6 || month == 11 {
		return day <= 30
	}
	return true //大月的规则已经检查过了
}

func isLeapYear(year int) bool {
	if 0 == (year % 400) {
		return true
	}
	if 0 == (year % 100) {
		return false
	}
	if 0 == (year % 4) {
		return true
	}
	return false
}
