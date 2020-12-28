package util

import "math"

func Abs(val interface{}) interface{} {
	switch val.(type) {
	case int:
		if val.(int) < 0 {
			return -val.(int)
		}
	case int8:
		if val.(int8) < 0 {
			return -val.(int8)
		}
	case int16:
		if val.(int16) < 0 {
			return -val.(int16)
		}
	case int32:
		if val.(int32) < 0 {
			return -val.(int32)
		}
	case int64:
		if val.(int64) < 0 {
			return -val.(int64)
		}
	case float32:
		if val.(float32) < 0 {
			return -val.(float32)
		}
	case float64:
		if val.(float64) < 0 {
			return -val.(float64)
		}
	default:
		return val
	}
	return val
}

func Round(val float64, precision int) (retval float64) {
	p := math.Pow10(precision)
	f := val * p
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return val
	}
	if f >= 0.0 {
		retval = math.Ceil(f)
		if (retval - f) > 0.50000000001 {
			retval -= 1.0
		}
	} else {
		retval = math.Ceil(-f)
		if (retval + f) > 0.50000000001 {
			retval -= 1.0
		}
		retval = -retval
	}
	return retval / p
}
