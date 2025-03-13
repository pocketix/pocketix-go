package utils

import "reflect"

func ToBool(value any) bool {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() != 0
	case float32, float64:
		return reflect.ValueOf(v).Float() != 0
	case string:
		return v != ""
	case bool:
		return v
	default:
		return false
	}
}
