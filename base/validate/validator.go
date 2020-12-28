package validate

import (
	"context"
	"fmt"
	"reflect"

	"toolbox/base/type_convert"
	"toolbox/log"
)

var (
	registeredValidators map[string]reflect.Type = make(map[string]reflect.Type)
)

func init() {
	registValidator(&Require{}, `required`)
	registValidator(&Range{}, `range`)
	registValidator(&Rangef{}, `rangef`)
	registValidator(&Min{}, `min`)
	registValidator(&Minf{}, `minf`)
	registValidator(&Max{}, `max`)
	registValidator(&Maxf{}, `maxf`)
	registValidator(&Length{}, `length`, `size`)
	registValidator(&MinLength{}, `min_length`, `min_size`)
	registValidator(&MaxLength{}, `max_length`, `max_size`)
	registValidator(&Alpha{}, `alpha`)
	registValidator(&Digit{}, `digit`)
	registValidator(&AlphaDigit{}, `alpha_digit`)
	registValidator(&AlphaDash{}, `alpha_dash`)
	registValidator(&Match{}, `match`, `regex`)
	registValidator(&Email{}, `email`)
	registValidator(&IPv4{}, `ipv4`)
	registValidator(&IPv6{}, `ipv6`)
	registValidator(&IP{}, `ip`)
	registValidator(&Base64{}, `base64`)
	registValidator(&Mobile{}, `mobile`)
	registValidator(&Tel{}, `tel`)
	registValidator(&Phone{}, `phone`)
	registValidator(&ZipCode{}, `zipcode`)
	registValidator(&Enum{}, `enum`)
}

func registValidator(v Validator, name ...string) {
	if len(name) == 0 {
		panic("Bad program, xtype.validate.registValidator should have at least one name") //错误的编程
	}
	reflectVal := reflect.ValueOf(v)
	typ := reflect.Indirect(reflectVal).Type()
	for _, n := range name {
		registeredValidators[n] = typ
	}
}

type Validator interface {
	Validate(ctx context.Context, d interface{}) bool
	Name() string
	Message(fieldName string, d interface{}) string
}

type Validators struct {
	Optional         bool
	ValidateHandlers []Validator
}

func (v Validators) Validate(ctx context.Context, fieldName string, d interface{}) *FieldError {
	for _, v := range v.ValidateHandlers {
		if !v.Validate(ctx, d) {
			return newBindError(fieldName, v.Name(), v.Message(fieldName, d))
		}
	}
	return nil
}

func CreateValidators(ctx context.Context, options []expression) *Validators {
	retval := Validators{}
	var (
		optional bool = false
		nonZero  bool = false
	)
	for _, vo := range options {
		if vo.Name == `optional` {
			optional = true
			continue
		}
		if vo.Name == `required` {
			nonZero = true
		}
		validator := createValidator(ctx, vo)
		if validator != nil {
			retval.ValidateHandlers = append(retval.ValidateHandlers, validator)
		}
	}
	if nonZero && optional {
		log.Warnf(ctx, "When creating Validators, both nonZero and optional are set, optional will be ignored")
		optional = false
	}
	retval.Optional = optional
	return &retval
}

func createValidator(ctx context.Context, opt expression) Validator {
	v, ok := registeredValidators[opt.Name]
	if !ok {
		log.Warnf(ctx, "Cann't create Validator with Name [%s], please see registeredValidators in validate/validator.go", opt.Name)
		return nil
	}
	newValue := reflect.New(v)
	argLen := len(opt.Args)
	if argLen > 0 { //需要绑定参数
		if v.Kind() == reflect.Struct {
			//顺序绑定参数
			elmVal := newValue.Elem()
			fn := elmVal.NumField()
			for i := 0; i < fn; i++ {
				fld := elmVal.Field(i)
				switch fld.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					v, err := type_convert.ToInt64(opt.Args[i])
					if err != nil {
						log.Warnf(ctx, "Error %s when bind %d 'th int argument to Validator with Name [%s]", err.Error(), i, opt.Name)
						continue
					}
					if fld.OverflowInt(v) {
						log.Warnf(ctx, "Overflow when bind %d 'th int argument to Validator with Name [%s]", i, opt.Name)
						continue
					}
					fld.SetInt(v)
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					v, err := type_convert.ToUint64(opt.Args[i])
					if err != nil {
						log.Warnf(ctx, "Error %s when bind %d 'th uint argument to Validator with Name [%s]", err.Error(), i, opt.Name)
						continue
					}
					if fld.OverflowUint(v) {
						log.Warnf(ctx, "Overflow when bind %d 'th uint argument to Validator with Name [%s]", i, opt.Name)
					}
					fld.SetUint(v)
				case reflect.Float64, reflect.Float32:
					v, err := type_convert.ToFloat64(opt.Args[i])
					if err != nil {
						log.Warnf(ctx, "Error %s when bind %d 'th float argument to Validator with Name [%s]", err.Error(), i, opt.Name)
					}
					if fld.OverflowFloat(v) {
						log.Warnf(ctx, "Overflow when bind %d 'th float argument to Validator with Name [%s]", i, opt.Name)
					}
					fld.SetFloat(v)
				case reflect.Bool:
					v, err := type_convert.ToBool(opt.Args[i])
					if err != nil {
						log.Warnf(ctx, "Error %s when bind %d 'th bool argument to Validator with Name [%s]", err.Error(), i, opt.Name)
					}
					fld.SetBool(v)
				case reflect.String:
					fld.SetString(opt.Args[i])
				case reflect.Slice:
					if i == 0 && fn == 1 { //只有一个参数，且为第一个参数为slice时才理解为[]string
						fld.Set(reflect.ValueOf(opt.Args))
						break
					} else {
						log.Warnf(ctx, "Only Can bind slice value to Validator with Name [%s] have only 1 arguments with []string type, current index is %d", opt.Name, i)
					}
				default:
					log.Warnf(ctx, "Can't bind %d 'th bool argument to Validator with Name [%s], argument type is %s", i, opt.Name, fld.Kind().String())
				}
			}
		} else {
			log.Warnf(ctx, "Validator for Name [%s] is not a struct, cann't bind args to it, it's type is %s", opt.Name, v.Kind().String())
		}
	}
	val := newValue.Interface()
	retval, ok := val.(Validator)
	if ok {
		return retval
	} else {
		log.Warnf(ctx, "Created Interface Is not a valid Validator, it's type is %T", val)
		return nil
	}
}

type Require struct{}
type Range struct {
	Min int64
	Max int64
}
type Rangef struct {
	Min float64
	Max float64
}
type Min struct {
	Min int64
}
type Minf struct {
	Min float64
}
type Max struct {
	Max int64
}
type Maxf struct {
	Max float64
}
type Length struct {
	Len int
}
type MinLength struct {
	Len int
}
type MaxLength struct {
	Len int
}
type Alpha struct{}
type Digit struct{}
type AlphaDigit struct{}
type Match struct {
	Pattern string
}
type AlphaDash struct{}
type Email struct{}
type IPv4 struct{}
type IPv6 struct{}
type IP struct{}
type Base64 struct{}
type Mobile struct{}
type Tel struct{}
type Phone struct{}
type ZipCode struct{}
type Enum struct {
	Values []string
}

func (v *Require) Validate(ctx context.Context, d interface{}) bool {
	return FnRequired(d)
}
func (v *Require) Name() string {
	return `required`
}
func (v *Require) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s is required, input %#v is not valid", fieldName, d)
}

func (v *Range) Validate(ctx context.Context, d interface{}) bool {
	i64, err := type_convert.ToInt64(d)
	if err != nil {
		log.Warnf(ctx, "type_convert.ToInt64(%#v) failed, in Range.Validate, error is %#v", d, err)
	}
	return FnRange(v.Min, v.Max, i64)
}
func (v *Range) Name() string {
	return `range`
}
func (v *Range) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s must between %d and %d", fieldName, v.Min, v.Max)
}

func (v *Rangef) Validate(ctx context.Context, d interface{}) bool {
	f64, err := type_convert.ToFloat64(d)
	if err != nil {
		log.Warnf(ctx, "type_convert.ToFloat64(%#v) failed, in Rangef.Validate, error is %#v", d, err)
	}
	return FnRangef(v.Min, v.Max, f64)
}
func (v *Rangef) Name() string {
	return `range`
}
func (v *Rangef) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s must between %f and %f", fieldName, v.Min, v.Max)
}

func (v *Min) Validate(ctx context.Context, d interface{}) bool {
	i64, err := type_convert.ToInt64(d)
	if err != nil {
		log.Warnf(ctx, "type_convert.ToInt64(%#v) failed, in Min.Validate, error is %#v", d, err)
	}
	return FnMin(v.Min, i64)
}
func (v *Min) Name() string {
	return `min`
}
func (v *Min) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s must greater than %d", fieldName, v.Min)
}

func (v *Minf) Validate(ctx context.Context, d interface{}) bool {
	f64, err := type_convert.ToFloat64(d)
	if err != nil {
		log.Warnf(ctx, "type_convert.ToFloat64(%#v) failed, in Minf.Validate, error is %#v", d, err)
	}
	return FnMinf(v.Min, f64)
}
func (v *Minf) Name() string {
	return `min`
}
func (v *Minf) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s must greater than %f", fieldName, v.Min)
}

func (v *Max) Validate(ctx context.Context, d interface{}) bool {
	i64, err := type_convert.ToInt64(d)
	if err != nil {
		log.Warnf(ctx, "type_convert.ToInt64(%#v) failed, in Max.Validate, error is %#v", d, err)
	}
	return FnMax(v.Max, i64)
}
func (v *Max) Name() string {
	return `max`
}
func (v *Max) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s must less than %d", fieldName, v.Max)
}

func (v *Maxf) Validate(ctx context.Context, d interface{}) bool {
	f64, err := type_convert.ToFloat64(d)
	if err != nil {
		log.Warnf(ctx, "type_convert.ToFloat64(%#v) failed, in Maxf.Validate, error is %#v", d, err)
	}
	return FnMaxf(v.Max, f64)
}
func (v *Maxf) Name() string {
	return `max`
}
func (v *Maxf) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s must less than %f", fieldName, v.Max)
}

func (v *Length) Validate(ctx context.Context, d interface{}) bool {
	return FnLength(ctx, v.Len, d)
}
func (v *Length) Name() string {
	return `length`
}
func (v *Length) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s length should equal %d", fieldName, v.Len)
}

func (v *MinLength) Validate(ctx context.Context, d interface{}) bool {
	return FnMinSize(ctx, v.Len, d)
}
func (v *MinLength) Name() string {
	return `min_length`
}
func (v *MinLength) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s length should greater than %d", fieldName, v.Len)
}

func (v *MaxLength) Validate(ctx context.Context, d interface{}) bool {
	return FnMaxSize(ctx, v.Len, d)
}
func (v *MaxLength) Name() string {
	return `max_length`
}
func (v *MaxLength) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s length should less than %d", fieldName, v.Len)
}

func (v *Alpha) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsAlpha(str)
	}
	log.Warnf(ctx, "validate.Alpha.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *Alpha) Name() string {
	return `alpha`
}
func (v *Alpha) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s only can contains alpha characters", fieldName)
}

func (v *Digit) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsDigit(str)
	}
	log.Warnf(ctx, "validate.Digit.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *Digit) Name() string {
	return `digit`
}
func (v *Digit) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s only can contains digit characters", fieldName)
}

func (v *AlphaDigit) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsAlphaDigit(str)
	}
	log.Warnf(ctx, "validate.AlphaDigit.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *AlphaDigit) Name() string {
	return `alpha_digit`
}
func (v *AlphaDigit) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s should only can contains alpha or digit characters", fieldName)
}

func (v *Match) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnMatch(ctx, v.Pattern, str)
	}
	log.Warnf(ctx, "validate.Match.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *Match) Name() string {
	return `match`
}
func (v *Match) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s cann't match regex %s", fieldName, v.Pattern)
}

func (v *AlphaDash) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsAlphaDash(str)
	}
	log.Warnf(ctx, "validate.AlphaDash.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *AlphaDash) Name() string {
	return `alpha_dash`
}
func (v *AlphaDash) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s should can only contains alpha or dash characters", fieldName)
}

func (v *Email) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsEmail(str)
	}
	log.Warnf(ctx, "validate.Email.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *Email) Name() string {
	return `email`
}
func (v *Email) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s is not a valid email", fieldName)
}

func (v *IPv4) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsIPv4(str)
	}
	log.Warnf(ctx, "validate.IPv4.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *IPv4) Name() string {
	return `ipv4`
}
func (v *IPv4) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s value is not a valid ipv4 string", fieldName)
}

func (v *IPv6) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsIPv6(str)
	}
	log.Warnf(ctx, "validate.IPv6.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *IPv6) Name() string {
	return `ipv6`
}
func (v *IPv6) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s value is not a valid ipv6 string", fieldName)
}

func (v *IP) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsIP(str)
	}
	log.Warnf(ctx, "validate.IP.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *IP) Name() string {
	return `ip`
}
func (v *IP) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s value is not a valid ipv4 or ipv6 string", fieldName)
}

func (v *Base64) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsBase64(str)
	}
	log.Warnf(ctx, "validate.Base64.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *Base64) Name() string {
	return `base64`
}
func (v *Base64) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s value is not a valid base64 string", fieldName)
}

func (v *Mobile) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsMobile(str)
	}
	log.Warnf(ctx, "validate.Mobile.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *Mobile) Name() string {
	return `mobile`
}
func (v *Mobile) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s value is not a valid mobile", fieldName)
}

func (v *Tel) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsTel(str)
	}
	log.Warnf(ctx, "validate.Tel.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *Tel) Name() string {
	return `tel`
}
func (v *Tel) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s value is not a valid tel", fieldName)
}

func (v *Phone) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsPhone(str)
	}
	log.Warnf(ctx, "validate.Phone.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *Phone) Name() string {
	return `phone`
}
func (v *Phone) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s value is not mobile or tel", fieldName)
}

func (v *ZipCode) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		return FnIsZipCode(str)
	}
	log.Warnf(ctx, "validate.ZipCode.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}
func (v *ZipCode) Name() string {
	return `zipcode`
}
func (v *ZipCode) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s value is not a valid zipcode", fieldName)
}

func (v *Enum) Validate(ctx context.Context, d interface{}) bool {
	if str, ok := d.(string); ok {
		for _, val := range v.Values {
			if val == str {
				return true
			}
		}
		return false
	}
	log.Warnf(ctx, "validate.Enum.Validate Only take string arguments, for other type of arguments, it always return false")
	return false
}

func (v *Enum) Name() string {
	return `enum`
}
func (v *Enum) Message(fieldName string, d interface{}) string {
	return fmt.Sprintf("Field %s value should in (%v)", fieldName, v.Values)
}
