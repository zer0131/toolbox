package validate

import (
	"context"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/zer0131/toolbox/base/type_convert"
	"github.com/zer0131/toolbox/base/util"
	"github.com/zer0131/toolbox/log"
)

var (
	regexpEmail *regexp.Regexp
)

func init() {
	var err error
	regexpEmail, err = regexp.Compile("(?i)^([0-9a-z._\\-]{2,64})@([a-z0-9\\-]{1,63})(\\.[a-z]{2,6})+$")
	if err != nil {
		panic(err)
	}
}

func FnRequired(val interface{}) bool {
	if val == nil {
		return false
	}

	switch v := val.(type) {
	case int:
		return v != 0
	case int8:
		return v != 0
	case int16:
		return v != 0
	case int32:
		return v != 0
	case int64:
		return v != 0
	case uint8:
		return v != 0
	case uint16:
		return v != 0
	case uint32:
		return v != 0
	case uint64:
		return v != 0
	case string:
		return v != ``
	case bool:
		return true //任何布尔值实际上都选择了
	case float32:
		return !isZeroFloat(float64(v))
	case float64:
		return !isZeroFloat(v)
	default:
		refVal := reflect.ValueOf(val)
		return fnRequiredReflect(refVal)
	}
}

func FnRange(min, max, val int64) bool {
	return val >= min && val <= max
}

func FnRangef(min, max, val float64) bool {
	return val >= min && val <= max
}

func FnMin(min, val int64) bool {
	return val >= min
}

func FnMax(max, val int64) bool {
	return val <= max
}

func FnMinf(min, val float64) bool {
	return val >= min
}

func FnMaxf(max, val float64) bool {
	return val <= max
}

func FnSize(ctx context.Context, length int, val interface{}) bool {
	return len4Interface(ctx, val) == length
}

func FnMinSize(ctx context.Context, length int, val interface{}) bool {
	return len4Interface(ctx, val) >= length
}

func FnMaxSize(ctx context.Context, length int, val interface{}) bool {
	return len4Interface(ctx, val) <= length
}

func FnLength(ctx context.Context, length int, val interface{}) bool {
	return len4Interface(ctx, val) == length
}

func len4Interface(ctx context.Context, val interface{}) int {
	refV := reflect.ValueOf(val)
	switch refV.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map:
		return refV.Len()
	case reflect.String:
		return util.StrLen(refV.String())
	default:
		log.Warnf(ctx, "Cann't evaluate length of %#v, this will return 0", val)
		return 0
	}
}

func FnIsAlpha(val string) bool {
	if val == `` {
		return false
	}
	runes := []rune(val)
	for _, r := range runes {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			continue
		}
		return false
	}
	return true
}

func FnIsDigit(val string) bool {
	if val == `` {
		return false
	}
	runes := []rune(val)
	for _, r := range runes {
		if r >= '0' && r <= '9' {
			continue
		}
		return false
	}
	return true
}

func FnIsAlphaDigit(val string) bool {
	if val == `` {
		return false
	}
	runes := []rune(val)
	for _, r := range runes {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') {
			continue
		}
		return false
	}
	return true
}

func FnIsAscii(val string) bool {
	if val == `` {
		return false
	}
	runes := []rune(val)
	for _, r := range runes {
		if r < 32 || r >= 127 {
			return false
		}
	}
	return true
}

func FnMatch(ctx context.Context, pattern, val string) bool {
	rex, err := regexp.Compile(pattern)
	if err != nil {
		log.Warnf(ctx, "Cann't Compile string %s to regexp in call validate.FnMatch(%#v, %#v), with error %s, this will make return false",
			pattern, pattern, val, err.Error())
		return false
	}
	return rex.MatchString(val)
}

func FnIsAlphaDash(val string) bool {
	if val == `` {
		return false
	}
	runes := []rune(val)
	for _, r := range runes {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r == '_' {
			continue
		}
		return false
	}
	return true
}

func FnIsEmail(val string) bool {
	if val == `` {
		return false
	}
	return regexpEmail.MatchString(val)
}

func FnIsIPv4(val string) bool {
	if val == `` {
		return false
	}
	segs := strings.Split(val, ".")
	segsLen := len(segs)
	if segsLen != 4 {
		return false
	}
	for i := 0; i < segsLen; i++ {
		tmp, err := type_convert.StringToInt64(segs[i])
		if err != nil {
			return false
		}
		if tmp < 0 || tmp > 255 {
			return false
		}
	}
	return true
}

func FnIsIPv6(val string) bool {
	if val == `` {
		return false
	}
	segs := strings.Split(val, ":")
	segsLen := len(segs)
	if segsLen < 3 { //形如: ::1 为最短形式
		return false
	}
	isHexChr := func(str string) bool {
		runes := []rune(str)
		for _, c := range runes {
			if (c >= '0' && c <= '9') ||
				(c >= 'A' && c <= 'F') ||
				(c >= 'a' && c <= 'f') {
				continue
			}
			return false
		}
		return true
	}
	if segs[0] == `` && segs[1] == `` { //形如: ::45de:78db:ef09,则从2开始向后不允许出现空的
		for i := 2; i < segsLen; i++ {
			if segs[i] == `` {
				return false
			}
			if !isHexChr(segs[i]) {
				return false
			}
		}
		if segsLen > 6 { //太长了,IPv6一共才８段
			return false
		}
		return true
	}
	var emptyFnd bool //手否找到了 ::
	for i := 0; i < segsLen; i++ {
		if segs[i] == `` {
			if emptyFnd {
				return false
			}
			emptyFnd = true
			continue
		}
		if !isHexChr(segs[i]) {
			return false
		}
	}
	if emptyFnd {
		if segsLen > 6 { //IPv6一共8段,至少２段缩写了
			return false
		}
	} else {
		if segsLen != 8 { //没有找到缩写则必须８段
			return false
		}
	}
	return true
}

func FnIsIP(val string) bool {
	if val == `` {
		return false
	}
	if FnIsIPv4(val) {
		return true
	}
	return FnIsIPv6(val)
}

func FnIsBase64(val string) bool {
	return false
}

func FnIsMobile(val string) bool {
	if strings.HasPrefix(val, "+86") { //考虑中国带国家区号的 +8618621557027 形式
		val = val[3:]
	}
	if len(val) != 11 {
		return false
	}
	if val[0] != '1' {
		return false
	}
	return FnIsDigit(val)
}

func FnIsTel(val string) bool {
	if val == `` {
		return false
	}
	return true
}

func FnIsPhone(val string) bool {
	return FnIsTel(val) || FnIsMobile(val)
}

func FnIsZipCode(val string) bool {
	return false
}

func filterLeftTrim(str string) string {
	return strings.TrimLeftFunc(str, unicode.IsSpace)
}

func filterRightTrim(str string) string {
	return strings.TrimRightFunc(str, unicode.IsSpace)
}

func filterTrimSpace(str string) string {
	return strings.TrimSpace(str)
}

func filterSBC2DBC(str string) string {
	return util.SBC2DBC(str)
}

func filterUCFirst(str string) string {
	return util.UCFirst(str)
}

func isZeroFloat(f float64) bool {
	return f >= -0.0000001 && f <= 0.0000001
}

func fnRequiredReflect(val reflect.Value) bool {
	if val.IsNil() || !val.IsValid() {
		return false
	}
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		return val.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() != 0
	case reflect.Float32, reflect.Float64:
		return !isZeroFloat(val.Float())
	case reflect.Bool:
		return true //布尔值永远视为已设置
	case reflect.String:
		return val.String() == ``
	case reflect.Array, reflect.Map, reflect.Slice:
		return val.Len() > 0
	case reflect.Struct:
		return true //结构体暂永远视为已设置
	case reflect.Invalid:
		return false
	case reflect.Ptr:
		valRef := val.Elem()
		return fnRequiredReflect(valRef)
	}
	return false
}
