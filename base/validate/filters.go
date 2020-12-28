package validate

import (
	"context"
	"reflect"

	"toolbox/base/type_convert"
	"toolbox/log"
)

var (
	registedFilters map[string]reflect.Type = make(map[string]reflect.Type)
)

func init() {
	registerFilter(&leftTrim{}, `lefttrim`, `left_trim`, `leftTrim`)
	registerFilter(&rightTrim{}, `righttrim`, `right_trim`, `rightTrim`)
	registerFilter(&trimSpace{}, `trim`, `trimSpace`, `trim_space`, `trimspace`)
	registerFilter(&ucfirst{}, `ucfirst`)
	registerFilter(&sbc2dbc{}, `sbc2dbc`)
}

type ValidateFilter interface {
	Filter(src string) string
}

type ValidateFilters struct {
	FilterHandlers []ValidateFilter
}

func (filters ValidateFilters) Filter(src string) string {
	for _, filter := range filters.FilterHandlers {
		src = filter.Filter(src)
	}
	return src
}

type leftTrim struct{}
type rightTrim struct{}
type trimSpace struct{}
type ucfirst struct{}
type sbc2dbc struct{}

func (o *leftTrim) Filter(src string) string {
	return filterLeftTrim(src)
}

func (o *rightTrim) Filter(src string) string {
	return filterRightTrim(src)
}

func (o *trimSpace) Filter(src string) string {
	return filterTrimSpace(src)
}

func (o *ucfirst) Filter(src string) string {
	return filterUCFirst(src)
}

func (o *sbc2dbc) Filter(src string) string {
	return filterSBC2DBC(src)
}

func createValidateFilters(ctx context.Context, options []expression) *ValidateFilters {
	retval := &ValidateFilters{}
	for _, opt := range options {
		filter := createValidateFilter(ctx, opt)
		if filter != nil {
			retval.FilterHandlers = append(retval.FilterHandlers, filter)
		}
	}
	return retval
}

func createValidateFilter(ctx context.Context, exp expression) ValidateFilter {
	v, ok := registedFilters[exp.Name]
	if !ok {
		log.Warnf(ctx, "Cann't create ValidatorFilter with Name [%s], please see registedFilters in validate/filters.go", exp.Name)
		return nil
	}
	newValue := reflect.New(v)
	argLen := len(exp.Args)
	if argLen > 0 { //需要绑定参数
		if v.Kind() == reflect.Struct {
			//顺序绑定参数
			elmVal := newValue.Elem()
			fn := elmVal.NumField()
			for i := 0; i < fn; i++ {
				fld := elmVal.Field(i)
				switch fld.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v, err := type_convert.ToInt64(exp.Args[i])
					if err != nil {
						log.Warnf(ctx, "Error %s when bind %d 'th int argument to ValidatorFilter with Name [%s]", err.Error(), i, exp.Name)
						continue
					}
					if fld.OverflowInt(v) {
						log.Warnf(ctx, "Overflow when bind %d 'th int argument to ValidatorFilter with Name [%s]", i, exp.Name)
						continue
					}
					fld.SetInt(v)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					v, err := type_convert.ToUint64(exp.Args[i])
					if err != nil {
						log.Warnf(ctx, "Error %s when bind %d 'th uint argument to ValidatorFilter with Name [%s]", err.Error(), i, exp.Name)
						continue
					}
					if fld.OverflowUint(v) {
						log.Warnf(ctx, "Overflow when bind %d 'th uint argument to ValidatorFilter with Name [%s]", i, exp.Name)
					}
					fld.SetUint(v)
				case reflect.Float64, reflect.Float32:
					v, err := type_convert.ToFloat64(exp.Args[i])
					if err != nil {
						log.Warnf(ctx, "Error %s when bind %d 'th float argument to ValidatorFilter with Name [%s]", err.Error(), i, exp.Name)
					}
					if fld.OverflowFloat(v) {
						log.Warnf(ctx, "Overflow when bind %d 'th float argument to ValidatorFilter with Name [%s]", i, exp.Name)
					}
					fld.SetFloat(v)
				case reflect.Bool:
					v, err := type_convert.ToBool(exp.Args[i])
					if err != nil {
						log.Warnf(ctx, "Error %s when bind %d 'th bool argument to ValidatorFilter with Name [%s]", err.Error(), i, exp.Name)
					}
					fld.SetBool(v)
				case reflect.String:
					fld.SetString(exp.Args[i])
				case reflect.Slice:
					if i == 0 && fn == 1 { //只有一个参数，且为第一个参数为slice时才理解为[]string
						fld.Set(reflect.ValueOf(exp.Args))
						break
					} else {
						log.Warnf(ctx, "Only Can bind slice value to ValidatorFilter with Name [%s] have only 1 arguments with []string type, current index is %d", exp.Name, i)
					}
				default:
					log.Warnf(ctx, "Can't bind %d 'th bool argument to ValidatorFilter with Name [%s], argument type is %s", i, exp.Name, fld.Kind().String())
				}
			}
		} else {
			log.Warnf(ctx, "ValidatorFilter for Name [%s] is not a struct, cann't bind args to it, it's type is %s", exp.Name, v.Kind().String())
		}
	}
	val := newValue.Interface()
	retval, ok := val.(ValidateFilter)
	if ok {
		return retval
	} else {
		log.Warnf(ctx, "Created Interface Is not a valid ValidateFilter, it's type is %T", val)
		return nil
	}
}

func registerFilter(filter ValidateFilter, name ...string) {
	if len(name) == 0 {
		panic("Bad program, xtype.validate.registerFilter should have at least one name") //错误的编程
	}
	reflectVal := reflect.ValueOf(filter)
	typ := reflect.Indirect(reflectVal).Type()
	for _, n := range name {
		registedFilters[n] = typ
	}
}
