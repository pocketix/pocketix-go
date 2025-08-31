package utils

import (
	"fmt"
)

func ToBool[T any](value T) (bool, error) {
	switch v := any(value).(type) {
	case int32, int64:
		return v != 0, nil
	case float32, float64:
		return v != 0.0, nil
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
