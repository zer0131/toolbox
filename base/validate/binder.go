package validate

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/zer0131/toolbox/base/type_convert"
	"github.com/zer0131/toolbox/log"
)

func BindFromHttpRequest(ctx context.Context, dest interface{}, req *http.Request) BindingResult {
	return BindFromValues(ctx, dest, parseValuesFromHttpRequest(ctx, req))
}

func BindFromJson(ctx context.Context, dest interface{}, json []byte) BindingResult {
	val, err := parseValuesFromJson(ctx, json)
	if err != nil {
		return BindingResult{OK: false, Err: errors.New(err.Error())}
	}
	return BindFromValues(ctx, dest, val)
}

func BindFromValues(ctx context.Context, dest interface{}, values url.Values) BindingResult {
	refPtr := reflect.ValueOf(dest)
	strct := refPtr.Elem()
	if refPtr.Kind() != reflect.Ptr {
		//编程错误,第一个参数必须是指针
		log.Errorf(ctx, "You must be program incorrrect, the 1st parameter dest of BindFromController must be pointer")
		return BindingResult{OK: false, Err: errors.New("parameter dest must be pointer")}
	}
	strctTyp := strct.Type()
	if strctTyp.Kind() != reflect.Struct {
		//编程错误,第一个参数必须是结构体
		log.Errorf(ctx, "You must be program incorrrect, the 1st parameter dest of BindFromController must be pointer of Struct")
		return BindingResult{OK: false, Err: errors.New("parameter dest must be pointer of struct")}
	}
	fieldNum := strctTyp.NumField()
	retval := BindingResult{OK: true}
	for i := 0; i < fieldNum; i++ {
		fieldDef := strctTyp.Field(i)
		field := strct.Field(i)
		err := bindValue(ctx, field, values, fieldDef)
		if err != nil {
			retval.OK = false
			retval.FieldErrors = append(retval.FieldErrors, *err)
		}
	}
	return retval
}

func bindValue(ctx context.Context, dest reflect.Value, values url.Values, strctField reflect.StructField) *FieldError {
	tag := strctField.Tag
	name := tag.Get(`name`)
	if name == `` {
		name = strings.ToLower(strctField.Name)
	}
	validateRules := tag.Get(`validate`)
	filterRules := tag.Get(`filter`)
	var (
		input      []string
		validators *Validators
		found      bool
	)
	if validateRules != `` {
		options, err := parseExpressions(validateRules)
		if err != nil {
			log.Warnf(ctx, "Error when parseExpressions(%#v) for ValidateRules, error is %s, Rule is [%s]", validateRules, err.Error(), validateRules)
		} else {
			validators = CreateValidators(ctx, options)
		}
	}
	input, found = values[name]
	if !found {
		if validators != nil && !validators.Optional {
			return newBindError(name, `require`, `Field %s is required, cann't be null`, name)
		}
		input = nil
	} else {
		if len(input) < 1 {
			log.Warnf(ctx, "Empty slice treat as empty, so cann't pass require rule")
			return newBindError(name, `require`, `Field %s is required, cann't be empty`, name)
		}
		if validators != nil && !validators.Optional {
			knd := strctField.Type.Kind()
			if knd == reflect.String ||
				knd == reflect.Bool ||
				knd == reflect.Float32 ||
				knd == reflect.Float64 ||
				knd == reflect.Int ||
				knd == reflect.Int8 ||
				knd == reflect.Int16 ||
				knd == reflect.Int32 ||
				knd == reflect.Int64 ||
				knd == reflect.Uint ||
				knd == reflect.Uint8 ||
				knd == reflect.Uint16 ||
				knd == reflect.Uint32 ||
				knd == reflect.Uint64 {
				if input[0] == `` {
					return newBindError(name, `require`, `Field %s is required, input is empty`, name)
				}
			} else if knd == reflect.Chan ||
				knd == reflect.Interface ||
				knd == reflect.Complex64 ||
				knd == reflect.Complex128 ||
				knd == reflect.Map ||
				knd == reflect.Ptr ||
				knd == reflect.Func ||
				knd == reflect.Struct ||
				knd == reflect.UnsafePointer {
				panic(`Cann't Validate for Unkown types`)
			}
		}
	}
	if input != nil {
		if filterRules != `` {
			options, err := parseExpressions(filterRules)
			if err != nil {
				log.Warnf(ctx, "Error when parseExpressions(%#v) for FilterRules, error is %s, Rule is [%s]", validateRules, err.Error(), filterRules)
			} else {
				filters := createValidateFilters(ctx, options)
				for idx, item := range input {
					input[idx] = filters.Filter(item)
				}
			}
		}
		//Do Bind
		typ := strctField.Type
		switch typ.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			v, err := type_convert.ToInt64(input[0])
			if err != nil {
				log.Warnf(ctx, "Cann't bind input %#v to a %s value", input, typ.Kind().String())
				return nil
			}
			if dest.OverflowInt(v) {
				log.Warnf(ctx, "Cann't bind input %#v to a %s value, data %d is overflowed", input, typ.Kind().String(), v)
				return nil
			}
			//执行验证过程
			if validators != nil {
				berr := validators.Validate(ctx, name, v)
				if berr != nil {
					return berr
				}
			}
			dest.SetInt(v)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			v, err := type_convert.ToUint64(input[0])
			if err != nil {
				log.Warnf(ctx, "Cann't bind input %#v to a %s value", input, typ.Kind().String())
				return nil
			}
			if dest.OverflowUint(v) {
				log.Warnf(ctx, "Cann't bind input %#v to a %s value, data %d is overflowed", input, typ.Kind().String(), v)
				return nil
			}
			//执行验证过程
			if validators != nil {
				berr := validators.Validate(ctx, name, v)
				if berr != nil {
					return berr
				}
			}
			dest.SetUint(v)
		case reflect.String:
			if validators != nil {
				berr := validators.Validate(ctx, name, input[0])
				if berr != nil {
					return berr
				}
			}
			dest.SetString(input[0])
		case reflect.Bool:
			v, err := type_convert.ToBool(input[0])
			if err != nil {
				log.Warnf(ctx, "Cann't bind input %#v to a bool value", input)
				v = false
			}
			dest.SetBool(v)
		case reflect.Float32, reflect.Float64:
			v, err := type_convert.ToFloat64(input[0])
			if err != nil {
				log.Warnf(ctx, "Cann't bind input %#v to a %s value", input, typ.Kind().String())
				return nil
			}
			if dest.OverflowFloat(v) {
				log.Warnf(ctx, "Cann't bind input %#v to a %s value, data %f is overflowed", input, typ.Kind().String(), v)
				return nil
			}
			//执行验证过程
			if validators != nil {
				berr := validators.Validate(ctx, name, v)
				if berr != nil {
					return berr
				}
			}
			dest.SetFloat(v)
		case reflect.Slice:
			elmTyp := typ.Elem()
			switch elmTyp.Kind() {
			case reflect.Int:
				v := toInts(input)
				dest.Set(reflect.ValueOf(v))
			case reflect.Int8:
				v := toInt8s(input)
				dest.Set(reflect.ValueOf(v))
			case reflect.Int16:
				v := toInt16s(input)
				dest.Set(reflect.ValueOf(v))
			case reflect.Int32:
				v := toInt32s(input)
				dest.Set(reflect.ValueOf(v))
			case reflect.Int64:
				v := toInt64s(input)
				dest.Set(reflect.ValueOf(v))
			case reflect.Uint:
				v := toUints(input)
				dest.Set(reflect.ValueOf(v))
			case reflect.Uint8:
				v := toUint8s(input)
				dest.Set(reflect.ValueOf(v))
			case reflect.Uint16:
				v := toUint16s(input)
				dest.Set(reflect.ValueOf(v))
			case reflect.Uint32:
				v := toUint32s(input)
				dest.Set(reflect.ValueOf(v))
			case reflect.Uint64:
				v := toUint64s(input)
				dest.Set(reflect.ValueOf(v))
			case reflect.String:
				dest.Set(reflect.ValueOf(input))
			case reflect.Bool:
				v := toBools(input)
				dest.Set(reflect.ValueOf(v))
			default:
				log.Warnf(ctx, "Cann't bind value %#v to slice type of %s", input, elmTyp.Kind().String())
			}
		}
	}
	return nil
}

func parseValuesFromHttpRequest(ctx context.Context, req *http.Request) url.Values {
	if req.Form == nil {
		err := req.ParseForm()
		if err != nil {
			log.Warnf(ctx, "try to parse form from request failed, error is %s", err.Error())
		}
	}
	retval := make(url.Values)
	for k, v := range req.Form {
		if k == `` {
			continue
		}
		retval[k] = v
		normaledKey := normalizeName(k)
		if normaledKey != k {
			retval[normaledKey] = v
		}
	}
	return retval
}

func parseValuesFromJson(ctx context.Context, jsonStr []byte) (retVal url.Values, err error) {
	var mapResult map[string]interface{}
	err = json.Unmarshal(jsonStr, &mapResult)
	if err != nil {
		log.Errorf(ctx, "try to parse json failed, error is %s", err.Error())
		return
	}

	retVal = make(url.Values)
	for k, v := range mapResult {
		if k == `` {
			continue
		}
		normaledKey := normalizeName(k)
		if reflect.ValueOf(v).Kind() == reflect.Slice {
			x := v.([]interface{})
			retVal[normaledKey] = type_convert.ToStrings(x)
		} else {
			retVal[normaledKey] = []string{type_convert.ToString(v)}
		}
	}
	return
}

func normalizeName(str string) string {
	runes := []rune(str)
	var (
		buf       bytes.Buffer
		needUpper bool = false
	)
	uc := func(c rune) rune {
		if c >= 'a' && c <= 'z' {
			return c - 32
		}
		return c
	}
	if runes[0] == '_' || runes[0] == '-' {
		needUpper = true
	} else {
		buf.WriteRune(uc(runes[0])) //避免循环中处理第0个字符
	}
	for i := 1; i < len(runes); i++ {
		if runes[i] == '_' || runes[i] == '-' {
			needUpper = true
			continue
		}
		if needUpper {
			buf.WriteRune(uc(runes[i]))
			needUpper = false
		} else {
			buf.WriteRune(runes[i])
		}
	}
	return buf.String()
}

func toInts(input []string) []int {
	ls := len(input)
	retval := make([]int, ls)
	for i := 0; i < ls; i++ {
		retval[i], _ = type_convert.ToInt(input[i])
	}
	return retval
}

func toInt8s(input []string) []int8 {
	ls := len(input)
	retval := make([]int8, ls)
	for i := 0; i < ls; i++ {
		retval[i], _ = type_convert.ToInt8(input[i])
	}
	return retval
}

func toInt16s(input []string) []int16 {
	ls := len(input)
	retval := make([]int16, ls)
	for i := 0; i < ls; i++ {
		retval[i], _ = type_convert.ToInt16(input[i])
	}
	return retval
}

func toInt32s(input []string) []int32 {
	ls := len(input)
	retval := make([]int32, ls)
	for i := 0; i < ls; i++ {
		retval[i], _ = type_convert.ToInt32(input[i])
	}
	return retval
}

func toInt64s(input []string) []int64 {
	ls := len(input)
	retval := make([]int64, ls)
	for i := 0; i < ls; i++ {
		retval[i], _ = type_convert.ToInt64(input[i])
	}
	return retval
}

func toUints(input []string) []uint {
	ls := len(input)
	retval := make([]uint, ls)
	for i := 0; i < ls; i++ {
		retval[i], _ = type_convert.ToUint(input[i])
	}
	return retval
}

func toUint8s(input []string) []uint8 {
	ls := len(input)
	retval := make([]uint8, ls)
	for i := 0; i < ls; i++ {
		retval[i], _ = type_convert.ToUint8(input[i])
	}
	return retval
}

func toUint16s(input []string) []uint16 {
	ls := len(input)
	retval := make([]uint16, ls)
	for i := 0; i < ls; i++ {
		retval[i], _ = type_convert.ToUint16(input[i])
	}
	return retval
}

func toUint32s(input []string) []uint32 {
	ls := len(input)
	retval := make([]uint32, ls)
	for i := 0; i < ls; i++ {
		retval[i], _ = type_convert.ToUint32(input[i])
	}
	return retval
}

func toUint64s(input []string) []uint64 {
	ls := len(input)
	retval := make([]uint64, ls)
	for i := 0; i < ls; i++ {
		retval[i], _ = type_convert.ToUint64(input[i])
	}
	return retval
}

func toBools(input []string) []bool {
	ls := len(input)
	retval := make([]bool, ls)
	for i := 0; i < ls; i++ {
		retval[i], _ = type_convert.ToBool(input[i])
	}
	return retval
}
