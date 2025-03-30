package utils

import (
	"fmt"
	"reflect"
)

// TODO: Return error if the value is not convertible to bool
func ToBool(value any) (bool, error) {
	switch v := value.(type) {
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(v).Int() != 0, nil
	case float32, float64:
		return reflect.ValueOf(v).Float() != 0, nil
	case string:
		if v == "true" {
			return true, nil
		}
		if v == "false" {
			return false, nil
		}
		return false, fmt.Errorf("cannot convert string %s to bool", v)
	case bool:
		return v, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

func BoolToInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
