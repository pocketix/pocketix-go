package models

import (
	"fmt"
	"slices"
	"strings"

	"github.com/pocketix/pocketix-go/src/utils"
)

type OperatorFunction func(child TreeNode, variableStore *VariableStore) (any, error)

type OperatorFactory struct {
	operator map[string]func(a, b any) (any, error)
}

func NewOperatorFactory() *OperatorFactory {
	return &OperatorFactory{
		operator: map[string]func(a, b any) (any, error){
			"===": func(a, b any) (any, error) {
				return CompareValues(a, b, func(x, y float64) bool { return x == y })
			},
			"!==": func(a, b any) (any, error) {
				return CompareValues(a, b, func(x, y float64) bool { return x != y })
			},
			"<": func(a, b any) (any, error) {
				if err := ForbidBoolean("<", a, b); err != nil {
					return nil, err
				}
				return CompareValues(a, b, func(x, y float64) bool { return x < y })
			},
			"<=": func(a, b any) (any, error) {
				if err := ForbidBoolean("<=", a, b); err != nil {
					return nil, err
				}
				return CompareValues(a, b, func(x, y float64) bool { return x <= y })
			},
			">": func(a, b any) (any, error) {
				if err := ForbidBoolean(">", a, b); err != nil {
					return nil, err
				}
				return CompareValues(a, b, func(x, y float64) bool { return x > y })
			},
			">=": func(a, b any) (any, error) {
				if err := ForbidBoolean(">", a, b); err != nil {
					return nil, err
				}
				return CompareValues(a, b, func(x, y float64) bool { return x >= y })
			},
			"+": func(a, b any) (any, error) {
				return AddValues(a, b)
			},
			"-": func(a, b any) (any, error) {
				if aFloat, bFloat, err := AllowNumber("-", a, b); err != nil {
					return nil, err
				} else {
					return aFloat - bFloat, nil
				}
			},
			"*": func(a, b any) (any, error) {
				if aFloat, bFloat, err := AllowNumber("*", a, b); err != nil {
					return nil, err
				} else {
					return aFloat * bFloat, nil
				}
			},
			"/": func(a, b any) (any, error) {
				if aFloat, bFloat, err := AllowNumber("/", a, b); err != nil {
					return nil, err
				} else if bFloat == 0 {
					return nil, fmt.Errorf("division by zero")
				} else {
					return aFloat / bFloat, nil
				}
			},
			"%": func(a, b any) (any, error) {
				if aFloat, bFloat, err := AllowNumber("%", a, b); err != nil {
					return nil, err
				} else if bFloat == 0 {
					return nil, fmt.Errorf("division by zero")
				} else {
					return int(aFloat) % int(bFloat), nil
				}
			},
			"&&": func(a, b any) (any, error) {
				aBool, aErr := utils.ToBool(a)
				if aErr != nil {
					return nil, aErr
				}
				bBool, bErr := utils.ToBool(b)
				if bErr != nil {
					return nil, bErr
				}
				return aBool && bBool, nil
			},
			"||": func(a, b any) (any, error) {
				aBool, aErr := utils.ToBool(a)
				if aErr != nil {
					return nil, aErr
				}
				bBool, bErr := utils.ToBool(b)
				if bErr != nil {
					return nil, bErr
				}
				return aBool || bBool, nil
			},
			"!": func(a, b any) (any, error) {
				aBool, aErr := utils.ToBool(a)
				if aErr != nil {
					return nil, aErr
				}
				return !aBool, nil
			},
		},
	}
}

func AddValues(a, b any) (any, error) {
	convertToFloat := func(v any) (float64, bool) {
		switch val := v.(type) {
		case float64:
			return val, true
		case int:
			return float64(val), true
		default:
			return 0, false
		}
	}

	if aFloat, okA := convertToFloat(a); okA {
		if bFloat, okB := convertToFloat(b); okB {
			return aFloat + bFloat, nil
		}
	}

	if aStr, okA := a.(string); okA {
		if bStr, okB := b.(string); okB {
			return aStr + bStr, nil
		}
	}

	return nil, fmt.Errorf("type mismatch: %T and %T", a, b)
}

func CompareValues(a, b any, comparator func(x, y float64) bool) (bool, error) {
	switch a := a.(type) {
	case bool:
		b, err := utils.ToBool(b)
		if err != nil {
			return false, err
		}
		// b, ok := b.(bool)
		// if !ok {
		// return false, fmt.Errorf("type mismatch: %T and %T", a, b)
		// }
		// "===" operator
		if comparator(1, 1) {
			return a == b, nil
		}
		// "!==" operator
		if comparator(1, 0) {
			return a != b, nil
		}
		return false, fmt.Errorf("unsupported operator for boolean: %T", a)
	case float64:
		switch b := b.(type) {
		case bool:
			if b {
				return comparator(a, 1), nil
			} else {
				return comparator(a, 0), nil
			}
		case float64:
			return comparator(a, b), nil
		case int:
			return comparator(a, float64(b)), nil
			// For now, forbid string comparison
		case string:
			return false, fmt.Errorf("type mismatch: %T and %T", a, b)
		default:
			return false, fmt.Errorf("unsupported type: %T", b)
		}
	case int:
		switch b := b.(type) {
		case bool:
			if b {
				return comparator(float64(a), 1), nil
			} else {
				return comparator(float64(a), 0), nil
			}
		case float64:
			return comparator(float64(a), b), nil
		case int:
			return comparator(float64(a), float64(b)), nil
			// For now, forbid string comparison
		case string:
			return false, fmt.Errorf("type mismatch: %T and %T", a, b)
		default:
			return false, fmt.Errorf("unsupported type: %T", b)
		}
	case string:
		b, ok := b.(string)
		if !ok {
			return false, fmt.Errorf("type mismatch: %T and %T", a, b)
		}
		return comparator(float64(strings.Compare(a, b)), 0), nil
	case nil:
		if b == nil {
			return false, nil
		}
		return false, fmt.Errorf("type mismatch: %T and %T", a, b)
	default:
		return false, fmt.Errorf("unsupported type: %T", a)
	}
}

func (o *OperatorFactory) ValidateOperator(node TreeNode) error {
	opFunc, exists := o.operator[node.Value.(string)]
	if !exists {
		return fmt.Errorf("operator not supported: %s", node.Value)
	}

	comparisonOperators := []string{"<", "<=", ">", ">=", "===", "!=="}
	if slices.Contains(comparisonOperators, node.Value.(string)) {
		_, err := ComparisonOperator(node, opFunc)
		if err != nil {
			return fmt.Errorf("error evaluating comparison operator: %s", err)
		}
		return nil
	}
	_, err := NumericLogicaloperator(node, opFunc)
	return err
}

func (o *OperatorFactory) EvaluateOperator(operator string, child TreeNode, variableStore *VariableStore, referenceValueStore *ReferencedValueStore) (any, error) {
	if len(child.Children) == 0 {
		if child.Type == "variable" {
			referencedValue, err := referenceValueStore.GetReferencedValueFromStore(child.Value.(string))
			if err == nil {
				sdParameterValue, valueType, err := referenceValueStore.ResolveParameterFunction(referencedValue.DeviceID, referencedValue.ParameterName)
				if err != nil {
					return nil, err
				}
				if sdParameterValue == nil {
					return nil, fmt.Errorf("referenced value %s not found", child.Value)
				}
				err = referenceValueStore.SetReferencedValue(child.Value.(string), sdParameterValue, valueType)
				if err != nil {
					return nil, err
				}
				return sdParameterValue, nil
			} else if variable, err := variableStore.GetVariable(child.Value.(string)); err == nil {
				if variable.Value.ResultValue == nil {
					return variable.Value.Value, nil
				}
				return variable.Value.ResultValue, nil
			}
			return nil, fmt.Errorf("referenced value %s or variable not found", child.Value)
		}
		// if child.Type == "device_variable" {}
		return child.Value, nil
	}

	opFunc, exists := o.operator[operator]
	if !exists {
		return nil, fmt.Errorf("operator not supported: %s", operator)
	}
	comparisonOperators := []string{"<", "<=", ">", ">=", "===", "!=="}

	if slices.Contains(comparisonOperators, operator) {
		return ComparisonOperator(child, opFunc)
	}

	return NumericLogicaloperator(child, opFunc)
}

func ComparisonOperator(child TreeNode, opFunc func(a, b any) (any, error)) (any, error) {
	for i := range len(child.Children) - 1 {
		// This part is commented out because it is not needed for now
		// Could be ussed in the future when the device variables will have a specific type
		// if child.Children[i].Type == "device_variable" || child.Children[i+1].Type == "device_variable" {
		// 	continue
		// }
		if (child.Children[i].ResultValue == nil && child.Children[i].Type == "variable") ||
			(child.Children[i+1].ResultValue == nil && child.Children[i+1].Type == "variable") {
			continue
		}

		a := child.Children[i].ResultValue
		b := child.Children[i+1].ResultValue

		result, err := opFunc(a, b)
		if err != nil {
			return nil, err
		}

		boolComp, ok := result.(bool)
		if !ok {
			return nil, fmt.Errorf("comparison operator %s must return boolean", child.Type)
		}
		if !boolComp {
			return false, nil
		}
	}
	return true, nil
}

func NumericLogicaloperator(child TreeNode, opFunc func(a, b any) (any, error)) (any, error) {
	if child.Children[0].Type == "variable" {
		return nil, nil
	}
	result := child.Children[0].ResultValue
	var err error

	for i := 1; i < len(child.Children); i++ {
		// if child.Children[i].Type == "device_variable" {
		// 	continue
		// }
		b := child.Children[i].ResultValue

		result, err = opFunc(result, b)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func ForbidBoolean(operator string, a, b any) error {
	if _, ok := a.(bool); ok {
		return fmt.Errorf("operator %s not supported for boolean", operator)
	}
	if _, ok := b.(bool); ok {
		return fmt.Errorf("operator %s not supported for boolean", operator)
	}
	return nil
}

func AllowNumber(operator string, a, b any) (float64, float64, error) {
	var aFloat, bFloat float64

	convertToFloat := func(v any) (float64, error) {
		switch value := v.(type) {
		case float64:
			return value, nil
		case int:
			return float64(value), nil
		default:
			return 0, fmt.Errorf("operator %s not supported for non-numeric type: %T", operator, v)
		}
	}

	if value, err := convertToFloat(a); err != nil {
		return 0, 0, err
	} else {
		aFloat = value
	}
	if value, err := convertToFloat(b); err != nil {
		return 0, 0, err
	} else {
		bFloat = value
	}
	return aFloat, bFloat, nil
}

func AllowBoolean(operator string, a, b any) error {
	if _, ok := a.(bool); !ok {
		return fmt.Errorf("operator %s not supported for non-boolean type: %T", operator, a)
	}
	if _, ok := b.(bool); !ok {
		return fmt.Errorf("operator %s not supported for non-boolean type: %T", operator, b)
	}
	return nil
}
