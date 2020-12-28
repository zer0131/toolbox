package type_convert

import (
	"math"
	"testing"
	"time"
)

type intConvertTest struct {
	Input         interface{} //用于转型的原值
	ExpectedValue int64       //期望的结果值
	ExpectedError error       //期望的错误
}

type strConvertTest struct {
	Input         string
	ExpectedValue int64
	ExpectedError error
}

type timeConvertTest struct {
	Input         interface{}
	ExpectedValue int64
	ExpectedError error
}

func TestToTime(t *testing.T) {
	args1 := time.Now()
	args2 := &args1
	args3 := 45178114014
	args4 := -45178114014
	args5 := "2016-01-01 23:59:61"
	args6 := "2017-05-03"
	args7 := "2015-02-29 23:05:08"
	args8 := "2016/01/01 23:59:59"
	args9 := "05/29/2015 08:00:00"
	tests := []timeConvertTest{
		timeConvertTest{args1, args1.Unix(), nil},
		timeConvertTest{45178114078, 45178114078, nil},
		timeConvertTest{"2016-01-01 23:59:61", 1451664001, nil}, //存疑，闰秒的处理和其它库并不完全兼容
		timeConvertTest{"2017-05-03", 1493740800, nil},
		timeConvertTest{"2015-02-29 23:05:08", 0, ERROR_TIME_NOTVALID},
		timeConvertTest{"2016/01/01 23:59:59", 1451663999, nil},
		timeConvertTest{"05/29/2015 08:00:00", 1432857600, nil},
		timeConvertTest{&args1, args1.Unix(), nil},
		timeConvertTest{&args2, args1.Unix(), nil}, //测试指针的指针
		timeConvertTest{&args3, 45178114014, nil},
		timeConvertTest{&args4, -45178114014, nil}, //负数的时间是否视为合法?
		timeConvertTest{&args5, 1451664001, nil},
		timeConvertTest{&args6, 1493740800, nil},
		timeConvertTest{&args7, 0, ERROR_TIME_NOTVALID},
		timeConvertTest{&args8, 1451663999, nil},
		timeConvertTest{&args9, 1432857600, nil},
	}
	for _, tst := range tests {
		tv, err := ToTime(tst.Input)
		if err != tst.ExpectedError || (err == nil && tst.ExpectedValue != tv.Unix()) {
			t.Errorf("Input %v cast to time failed, expected result %d, error: %v, but real get %d, error %v",
				tst.Input, tst.ExpectedValue, tst.ExpectedError, tv.Unix(), err)
			t.Fail()
		}
	}
}

func TestStringToInt64(t *testing.T) {
	test := []strConvertTest{
		strConvertTest{"-19.71", -19, nil},
		strConvertTest{"78421", 78421, nil},
		strConvertTest{"0x12", 18, nil},
		strConvertTest{"071", 57, nil},
		strConvertTest{"75412111", 75412111, nil},
		strConvertTest{"-75412111", -75412111, nil},
	}
	for _, tst := range test {
		val, err := StringToInt64(tst.Input)
		if val != tst.ExpectedValue && err != tst.ExpectedError {
			t.Errorf("StringToUint64(%s) should return %d, but real return %d \n", tst.Input, tst.ExpectedValue, val)
			t.Fail()
		}
	}
}

func TestStringToUint642(t *testing.T) {
	test := []strConvertTest{
		strConvertTest{"-19.71", 0, ErrUnderflow},
		strConvertTest{"78421", 78421, nil},
		strConvertTest{"0x12", 18, nil},
		strConvertTest{"071", 57, nil},
		strConvertTest{"75412111", 75412111, nil},
		strConvertTest{"-75412111", 0, ErrUnderflow},
	}
	for _, tst := range test {
		val, err := StringToUint64(tst.Input)
		if int64(val) != tst.ExpectedValue || err != tst.ExpectedError {
			t.Errorf("StringToUint64(%s) should return %d, error %v, but real return %d, error %v \n", tst.Input, tst.ExpectedValue, tst.ExpectedError, val, err)
			t.Fail()
		}
	}
}

func TestToInt8(t *testing.T) {
	args1 := 4 //用于指针测试
	args2 := 654
	args3 := "0x12"
	args4 := "12.78"
	args5 := "-19.71"
	test := []intConvertTest{
		intConvertTest{1, 1, nil},
		intConvertTest{10234, 0, ErrOverflow},
		intConvertTest{-17231, 0, ErrUnderflow},
		intConvertTest{1.2, 1, nil},
		intConvertTest{983117512241, 0, ErrOverflow},
		intConvertTest{78124824, 0, ErrOverflow},
		intConvertTest{-12901, 0, ErrUnderflow},
		intConvertTest{-2, -2, nil},
		intConvertTest{"0x4", 4, nil},
		intConvertTest{"0x12", 18, nil},
		intConvertTest{"017", 15, nil},
		intConvertTest{"234", 0, ErrOverflow},
		intConvertTest{"12.78", 12, ErrBadFormat},
		intConvertTest{"7841.12", 0, ErrOverflow},
		intConvertTest{"-7841.12", 0, ErrUnderflow},
		intConvertTest{"12A", 12, ErrBadFormat},
		intConvertTest{"-12A", -12, ErrBadFormat},
		intConvertTest{'0', 48, nil}, //byte　to int
		intConvertTest{'H', 72, nil},
		intConvertTest{'张', 0, ErrOverflow}, //中文转型
		intConvertTest{&args1, 4, nil},
		intConvertTest{&args2, 0, ErrOverflow},
		intConvertTest{&args3, 18, nil},
		intConvertTest{&args4, 12, ErrBadFormat},
		intConvertTest{&args5, -19, ErrBadFormat}}
	for _, v := range test {
		i8, e := ToInt8(v.Input)
		if int64(i8) != v.ExpectedValue || e != v.ExpectedError {
			t.Errorf("Input %v toInt8 failed, expceted %d, error %v, but get %d, error %v \n", v.Input, v.ExpectedValue, v.ExpectedError, i8, e)
			t.Fail()
		}
	}
}

func TestToInt16(t *testing.T) {
	args1 := 4 //用于指针测试
	args2 := 654
	args3 := "0x12"
	args4 := "12.78"
	args5 := "-19.71"
	args6 := math.MaxInt16 + 1
	args7 := math.MinInt16 - 1
	test := []intConvertTest{
		intConvertTest{1, 1, nil},
		intConvertTest{math.MaxInt16 + 12, 0, ErrOverflow},
		intConvertTest{math.MinInt16 - 1, 0, ErrUnderflow},
		intConvertTest{1.2, 1, nil},
		intConvertTest{983117512241, 0, ErrOverflow},
		intConvertTest{78124824, 0, ErrOverflow},
		intConvertTest{-12901, -12901, nil},
		intConvertTest{-2, -2, nil},
		intConvertTest{"0x4", 4, nil},
		intConvertTest{"0x12", 18, nil},
		intConvertTest{"017", 15, nil},
		intConvertTest{"234", 234, nil},
		intConvertTest{"12.78", 12, ErrBadFormat},
		intConvertTest{"7841.12", 7841, ErrBadFormat},
		intConvertTest{"-7841.12", -7841, ErrBadFormat},
		intConvertTest{"12A", 12, ErrBadFormat},
		intConvertTest{"-12A", -12, ErrBadFormat},
		intConvertTest{'0', 48, nil}, //byte　to int
		intConvertTest{'H', 72, nil},
		intConvertTest{'张', 24352, nil}, //中文转型
		intConvertTest{&args1, 4, nil},
		intConvertTest{&args2, 654, nil},
		intConvertTest{&args3, 18, nil},
		intConvertTest{&args4, 12, ErrBadFormat},
		intConvertTest{&args5, -19, ErrBadFormat},
		intConvertTest{&args6, 0, ErrOverflow},
		intConvertTest{&args7, 0, ErrUnderflow}}
	for _, v := range test {
		i8, e := ToInt16(v.Input)
		if int64(i8) != v.ExpectedValue || e != v.ExpectedError {
			t.Errorf("Input %v toInt16 failed, expceted %d, error %v, but get %d, error %v \n", v.Input, v.ExpectedValue, v.ExpectedError, i8, e)
			t.Fail()
		}
	}
}

func TestToInt32(t *testing.T) {
	args1 := 4 //用于指针测试
	args2 := 654
	args3 := "0x12"
	args4 := "12.78"
	args5 := "-19.71"
	args6 := math.MaxInt32 + 1
	args7 := math.MinInt32 - 1
	test := []intConvertTest{
		intConvertTest{1, 1, nil},
		intConvertTest{math.MaxInt32 + 12, 0, ErrOverflow},
		intConvertTest{math.MinInt32 - 1, 0, ErrUnderflow},
		intConvertTest{1.2, 1, nil},
		intConvertTest{983117512241, 0, ErrOverflow},
		intConvertTest{78124824, 78124824, nil},
		intConvertTest{-12901, -12901, nil},
		intConvertTest{-2, -2, nil},
		intConvertTest{"0x4", 4, nil},
		intConvertTest{"0x12", 18, nil},
		intConvertTest{"017", 15, nil},
		intConvertTest{"234", 234, nil},
		intConvertTest{"12.78", 12, ErrBadFormat},
		intConvertTest{"7841.12", 7841, ErrBadFormat},
		intConvertTest{"-7841.12", -7841, ErrBadFormat},
		intConvertTest{"12A", 12, ErrBadFormat},
		intConvertTest{"-12A", -12, ErrBadFormat},
		intConvertTest{'0', 48, nil}, //byte　to int
		intConvertTest{'H', 72, nil},
		intConvertTest{'张', 24352, nil}, //中文转型
		intConvertTest{&args1, 4, nil},
		intConvertTest{&args2, 654, nil},
		intConvertTest{&args3, 18, nil},
		intConvertTest{&args4, 12, ErrBadFormat},
		intConvertTest{&args5, -19, ErrBadFormat},
		intConvertTest{&args6, 0, ErrOverflow},
		intConvertTest{&args7, 0, ErrUnderflow}}
	for _, v := range test {
		i8, e := ToInt32(v.Input)
		if int64(i8) != v.ExpectedValue || e != v.ExpectedError {
			t.Errorf("Input %v toInt32 failed, expceted %d, error %v, but get %d, error %v \n", v.Input, v.ExpectedValue, v.ExpectedError, i8, e)
			t.Fail()
		}
	}
}

func TestToInt64(t *testing.T) {
	args1 := 4
	args2 := 654
	args3 := "0x12"
	args4 := "12.78"
	args5 := "-19.71"
	args6 := math.MaxInt32 + 1
	args7 := math.MinInt32 - 1
	test := []intConvertTest{
		intConvertTest{1, 1, nil},
		intConvertTest{math.MaxInt32 + 12, math.MaxInt32 + 12, nil},
		intConvertTest{math.MinInt32 - 1, math.MinInt32 - 1, nil},
		intConvertTest{1.2, 1, nil},
		intConvertTest{983117512241, 983117512241, nil},
		intConvertTest{78124824, 78124824, nil},
		intConvertTest{-12901, -12901, nil},
		intConvertTest{-2, -2, nil},
		intConvertTest{"0x4", 4, nil},
		intConvertTest{"0x12", 18, nil},
		intConvertTest{"017", 15, nil},
		intConvertTest{"234", 234, nil},
		intConvertTest{"12.78", 12, ErrBadFormat},
		intConvertTest{"7841.12", 7841, ErrBadFormat},
		intConvertTest{"-7841.12", -7841, ErrBadFormat},
		intConvertTest{"12A", 12, ErrBadFormat},
		intConvertTest{"-12A", -12, ErrBadFormat},
		intConvertTest{'0', 48, nil}, //byte　to int
		intConvertTest{'H', 72, nil},
		intConvertTest{'张', 24352, nil}, //中文转型
		intConvertTest{"983117512241", 983117512241, nil},
		intConvertTest{"-983117512241", -983117512241, nil},
		intConvertTest{&args1, 4, nil},
		intConvertTest{&args2, 654, nil},
		intConvertTest{&args3, 18, nil},
		intConvertTest{&args4, 12, ErrBadFormat},
		intConvertTest{&args5, -19, ErrBadFormat},
		intConvertTest{&args6, math.MaxInt32 + 1, nil},
		intConvertTest{&args7, math.MinInt32 - 1, nil}}
	for _, v := range test {
		i8, e := ToInt64(v.Input)
		if int64(i8) != v.ExpectedValue || e != v.ExpectedError {
			t.Errorf("Input %v toInt64 failed, expceted %d, error %v, but get %d, error %v \n", v.Input, v.ExpectedValue, v.ExpectedError, i8, e)
			t.Fail()
		}
	}
}

func TestToUint8(t *testing.T) {
	args1 := 4
	args2 := 654
	args3 := "0x12"
	args4 := "12.78"
	args5 := "-19.71"
	args6 := math.MaxInt32 + 1
	args7 := math.MinInt32 - 1
	test := []intConvertTest{
		intConvertTest{1, 1, nil},
		intConvertTest{math.MaxInt32 + 12, 0, ErrOverflow},
		intConvertTest{math.MinInt32 - 1, 0, ErrUnderflow},
		intConvertTest{1.2, 1, nil},
		intConvertTest{983117512241, 0, ErrOverflow},
		intConvertTest{78124824, 0, ErrOverflow},
		intConvertTest{-12901, 0, ErrUnderflow},
		intConvertTest{-2, 0, ErrUnderflow},
		intConvertTest{"0x4", 4, nil},
		intConvertTest{"0x12", 18, nil},
		intConvertTest{"017", 15, nil},
		intConvertTest{"234", 234, nil},
		intConvertTest{"12.78", 12, ErrBadFormat},
		intConvertTest{"7841.12", 0, ErrOverflow},
		intConvertTest{"-7841.12", 0, ErrUnderflow},
		intConvertTest{"12A", 12, ErrBadFormat},
		intConvertTest{"-12A", 0, ErrUnderflow},
		intConvertTest{'0', 48, nil}, //byte　to int
		intConvertTest{'H', 72, nil},
		intConvertTest{'张', 0, ErrOverflow}, //中文转型
		intConvertTest{"983117512241", 0, ErrOverflow},
		intConvertTest{"-983117512241", 0, ErrUnderflow},
		intConvertTest{&args1, 4, nil},
		intConvertTest{&args2, 0, ErrOverflow},
		intConvertTest{&args3, 18, nil},
		intConvertTest{&args4, 12, ErrBadFormat},
		intConvertTest{&args5, 0, ErrUnderflow},
		intConvertTest{&args6, 0, ErrOverflow},
		intConvertTest{&args7, 0, ErrUnderflow}}
	for _, v := range test {
		u, e := ToUint8(v.Input)
		if int64(u) != v.ExpectedValue || e != v.ExpectedError {
			t.Errorf("Input %v ToUint8 failed, expceted %d, error %v, but get %d, error %v \n", v.Input, v.ExpectedValue, v.ExpectedError, u, e)
			t.Fail()
		}
	}
}

func TestToUint16(t *testing.T) {
	args1 := 4
	args2 := 654
	args3 := "0x12"
	args4 := "12.78"
	args5 := "-19.71"
	args6 := math.MaxInt32 + 1
	args7 := math.MinInt32 - 1
	test := []intConvertTest{
		intConvertTest{1, 1, nil},
		intConvertTest{math.MaxInt32 + 12, 0, ErrOverflow},
		intConvertTest{math.MinInt32 - 1, 0, ErrUnderflow},
		intConvertTest{1.2, 1, nil},
		intConvertTest{983117512241, 0, ErrOverflow},
		intConvertTest{78124824, 0, ErrOverflow},
		intConvertTest{-12901, 0, ErrUnderflow},
		intConvertTest{-2, 0, ErrUnderflow},
		intConvertTest{"0x4", 4, nil},
		intConvertTest{"0x12", 18, nil},
		intConvertTest{"017", 15, nil},
		intConvertTest{"234", 234, nil},
		intConvertTest{"12.78", 12, ErrBadFormat},
		intConvertTest{"7841.12", 7841, ErrBadFormat},
		intConvertTest{"-7841.12", 0, ErrUnderflow},
		intConvertTest{"12A", 12, ErrBadFormat},
		intConvertTest{"-12A", 0, ErrUnderflow},
		intConvertTest{'0', 48, nil}, //byte　to int
		intConvertTest{'H', 72, nil},
		intConvertTest{'张', 24352, nil}, //中文转型
		intConvertTest{"983117512241", 0, ErrOverflow},
		intConvertTest{"-983117512241", 0, ErrUnderflow},
		intConvertTest{&args1, 4, nil},
		intConvertTest{&args2, 654, nil},
		intConvertTest{&args3, 18, nil},
		intConvertTest{&args4, 12, ErrBadFormat},
		intConvertTest{&args5, 0, ErrUnderflow},
		intConvertTest{&args6, 0, ErrOverflow},
		intConvertTest{&args7, 0, ErrUnderflow}}
	for _, v := range test {
		u16, e := ToUint16(v.Input)
		if int64(u16) != v.ExpectedValue || e != v.ExpectedError {
			if str, ok := v.Input.(*string); ok {
				t.Errorf("Input is %s \n", *str)
			}
			if i, ok := v.Input.(*int); ok {
				t.Errorf("Input is %d \n", *i)
			}
			t.Errorf("Input %v ToUint16 failed, expceted %d, error %v, but get %d, error %v \n", v.Input, v.ExpectedValue, v.ExpectedError, u16, e)
			t.Fail()
		}
	}
}

func TestToUint32(t *testing.T) {
	args1 := 4
	args2 := 654
	args3 := "0x12"
	args4 := "12.78"
	args5 := "-19.71"
	args6 := math.MaxInt32 + 1
	args7 := math.MinInt32 - 1
	test := []intConvertTest{
		intConvertTest{1, 1, nil},
		intConvertTest{math.MaxInt32 + 12, math.MaxInt32 + 12, nil},
		intConvertTest{math.MinInt32 - 1, 0, ErrUnderflow},
		intConvertTest{1.2, 1, nil},
		intConvertTest{983117512241, 0, ErrOverflow},
		intConvertTest{78124824, 78124824, nil},
		intConvertTest{-12901, 0, ErrUnderflow},
		intConvertTest{-2, 0, ErrUnderflow},
		intConvertTest{"0x4", 4, nil},
		intConvertTest{"0x12", 18, nil},
		intConvertTest{"017", 15, nil},
		intConvertTest{"234", 234, nil},
		intConvertTest{"12.78", 12, ErrBadFormat},
		intConvertTest{"7841.12", 7841, ErrBadFormat},
		intConvertTest{"-7841.12", 0, ErrUnderflow},
		intConvertTest{"12A", 12, ErrBadFormat},
		intConvertTest{"-12A", 0, ErrUnderflow},
		intConvertTest{'0', 48, nil}, //byte　to int
		intConvertTest{'H', 72, nil},
		intConvertTest{'张', 24352, nil}, //中文转型
		intConvertTest{"983117512241", 0, ErrOverflow},
		intConvertTest{"-983117512241", 0, ErrUnderflow},
		intConvertTest{&args1, 4, nil},
		intConvertTest{&args2, 654, nil},
		intConvertTest{&args3, 18, nil},
		intConvertTest{&args4, 12, ErrBadFormat},
		intConvertTest{&args5, 0, ErrUnderflow},
		intConvertTest{&args6, math.MaxInt32 + 1, nil},
		intConvertTest{&args7, 0, ErrUnderflow}}
	for _, v := range test {
		u32, e := ToUint32(v.Input)
		if int64(u32) != v.ExpectedValue || e != v.ExpectedError {
			if str, ok := v.Input.(*string); ok {
				t.Errorf("Input is %s \n", *str)
			}
			if i, ok := v.Input.(*int); ok {
				t.Errorf("Input is %d \n", *i)
			}
			t.Errorf("Input %v ToUint16 failed, expceted %d, error %v, but get %d, error %v \n", v.Input, v.ExpectedValue, v.ExpectedError, u32, e)
			t.Fail()
		}
	}
}

func TestToInt(t *testing.T) {
	args1 := 4
	args2 := 654
	args3 := "0x12"
	args4 := "12.78"
	args5 := "-19.71"
	args6 := math.MaxInt32 + 1
	args7 := math.MinInt32 - 1
	test := []intConvertTest{
		intConvertTest{1, 1, nil},
		intConvertTest{math.MaxInt32 + 12, math.MaxInt32 + 12, nil},
		intConvertTest{math.MinInt32 - 1, math.MinInt32 - 1, nil},
		intConvertTest{1.2, 1, nil},
		intConvertTest{983117512241, 983117512241, nil},
		intConvertTest{78124824, 78124824, nil},
		intConvertTest{-12901, -12901, nil},
		intConvertTest{-2, -2, nil},
		intConvertTest{"0x4", 4, nil},
		intConvertTest{"0x12", 18, nil},
		intConvertTest{"017", 15, nil},
		intConvertTest{"234", 234, nil},
		intConvertTest{"12.78", 12, ErrBadFormat},
		intConvertTest{"7841.12", 7841, ErrBadFormat}, //小数字符串转int，也会报ErrBadFormat
		intConvertTest{"-7841.12", -7841, ErrBadFormat},
		intConvertTest{"12A", 12, ErrBadFormat},
		intConvertTest{"-12A", -12, ErrBadFormat},
		intConvertTest{'0', 48, nil}, //byte　to int
		intConvertTest{'H', 72, nil},
		intConvertTest{'张', 24352, nil}, //中文转型
		intConvertTest{"983117512241", 983117512241, nil},
		intConvertTest{"-983117512241", -983117512241, nil},
		intConvertTest{&args1, 4, nil},
		intConvertTest{&args2, 654, nil},
		intConvertTest{&args3, 18, nil},
		intConvertTest{&args4, 12, ErrBadFormat},
		intConvertTest{&args5, -19, ErrBadFormat},
		intConvertTest{&args6, math.MaxInt32 + 1, nil},
		intConvertTest{&args7, math.MinInt32 - 1, nil}}
	for _, v := range test {
		u32, e := ToInt(v.Input)
		if int64(u32) != v.ExpectedValue || e != v.ExpectedError {
			if str, ok := v.Input.(*string); ok {
				t.Errorf("Input is %s \n", *str)
			}
			if i, ok := v.Input.(*int); ok {
				t.Errorf("Input is %d \n", *i)
			}
			t.Errorf("Input %v ToInt failed, expceted %d, error %v, but get %d, error %v \n", v.Input, v.ExpectedValue, v.ExpectedError, u32, e)
			t.Fail()
		}
	}
}

type boolConvertTest struct {
	Input         interface{}
	ExpectedValue bool
	ExpectedError error
}

func TestToBool(t *testing.T) {
	args1 := true
	args2 := &args1
	args3 := "true"
	args4 := "yes"
	args5 := 12
	args6 := 0
	args7 := "no"
	args8 := "false"
	args9 := "好好学习，天天向上"
	tests := []boolConvertTest{
		boolConvertTest{true, true, nil},
		boolConvertTest{false, false, nil},
		boolConvertTest{1, true, nil},
		boolConvertTest{0, false, nil},
		boolConvertTest{"yes", true, nil},
		boolConvertTest{"no", false, nil},
		boolConvertTest{"1", true, nil},
		boolConvertTest{"0", false, nil},
		boolConvertTest{"Y", true, nil},
		boolConvertTest{"N", false, nil},
		boolConvertTest{"yEs", true, nil},
		boolConvertTest{"y", true, nil},
		boolConvertTest{"这段话应该无法转成Bool", false, ErrBadFormat},
		boolConvertTest{&args1, true, nil}, //指针应当能够正确处理
		boolConvertTest{&args2, true, nil}, //指针的指针而已
		boolConvertTest{&args3, true, nil},
		boolConvertTest{&args4, true, nil},
		boolConvertTest{&args5, true, nil},
		boolConvertTest{&args6, false, nil},
		boolConvertTest{&args7, false, nil},
		boolConvertTest{&args8, false, nil},
		boolConvertTest{&args9, false, ErrBadFormat},
	}
	for _, tst := range tests {
		tv, err := ToBool(tst.Input)
		if tv != tst.ExpectedValue || err != tst.ExpectedError {
			t.Errorf("Input %v ToBool failed, expect %v and error %v, but got %v and %v", tst.Input, tst.ExpectedValue, tst.ExpectedError, tv, err)
			t.Fail()
		}
	}
}

type floatConvertTest struct {
	Input         interface{}
	ExceptedValue float64
	ExceptedError error
}

func TestToFloat32(t *testing.T) {
	args1 := 1.02
	//args2	:= float64(math.MaxFloat32) + float64(10.2)
	args3 := &args1 //测试指针的指针
	args4 := "-127.98.45"
	args5 := "3127.98A"
	args6 := "0x1847.4" //float类型不支持进制表示
	args7 := "abced"
	args8 := true
	tests := []floatConvertTest{
		floatConvertTest{452.798, 452.798, nil},
		//floatConvertTest{float64(math.MaxFloat32) + float64(7841.1), 0, nil},
		floatConvertTest{true, 1, nil},
		floatConvertTest{false, 0, nil},
		floatConvertTest{"-128743.34", -128743.34, nil},
		floatConvertTest{"-128743.34.89", -128743.34, ErrBadFormat},
		floatConvertTest{"128743.34defq", 128743.34, ErrBadFormat},
		floatConvertTest{&args1, 1.02, nil},
		//floatConvertTest{&args2, 0, ERROR_NUMBER_OVERFLOW},
		floatConvertTest{&args3, 1.02, nil},
		floatConvertTest{&args4, -127.98, ErrBadFormat},
		floatConvertTest{&args5, 3127.98, ErrBadFormat},
		floatConvertTest{&args6, 0, ErrBadFormat},
		floatConvertTest{&args7, 0, ErrBadFormat},
		floatConvertTest{&args8, 1, nil},
	}
	for _, tst := range tests {
		fv, err := ToFloat32(tst.Input)
		if !floatEquals(float64(fv), tst.ExceptedValue) || err != tst.ExceptedError {
			t.Errorf("Input %v ToFloat32 failed, expect %v and error %v, but got %v and %v", tst.Input, tst.ExceptedValue, tst.ExceptedError, fv, err)
			t.Fail()
		}
	}
}

func TestToFloat64(t *testing.T) {
	args1 := 1.02
	args2 := math.MaxFloat32 + 10.2
	args3 := &args1 //测试指针的指针
	args4 := "-127.98.45"
	args5 := "3127.98A"
	args6 := "0x1847.4" //float类型不支持进制表示
	args7 := "abced"
	args8 := true
	tests := []floatConvertTest{
		floatConvertTest{452.798, 452.798, nil},
		floatConvertTest{math.MaxFloat32 + 7841.1, math.MaxFloat32 + 7841.1, nil},
		floatConvertTest{true, 1, nil},
		floatConvertTest{false, 0, nil},
		floatConvertTest{"-128743.34", -128743.34, nil},
		floatConvertTest{"-128743.34.89", -128743.34, ErrBadFormat},
		floatConvertTest{"128743.34defq", 128743.34, ErrBadFormat},
		floatConvertTest{&args1, 1.02, nil},
		floatConvertTest{&args2, math.MaxFloat32 + 10.2, nil},
		floatConvertTest{&args3, 1.02, nil},
		floatConvertTest{&args4, -127.98, ErrBadFormat},
		floatConvertTest{&args5, 3127.98, ErrBadFormat},
		floatConvertTest{&args6, 0, ErrBadFormat},
		floatConvertTest{&args7, 0, ErrBadFormat},
		floatConvertTest{&args8, 1, nil},
	}
	for _, tst := range tests {
		fv, err := ToFloat64(tst.Input)
		if !floatEquals(fv, tst.ExceptedValue) || err != tst.ExceptedError {
			t.Errorf("Input %v ToFloat64 failed, expect %v and error %v, but got %v and %v", tst.Input, tst.ExceptedValue, tst.ExceptedError, fv, err)
			t.Fail()
		}
	}
}

const FLOAT_TOLERANCE = 0.01

func floatEquals(f1, f2 float64) bool {
	return math.Abs(f1-f2) < FLOAT_TOLERANCE
}
