package models

import (
	"fmt"
	"strings"

	"github.com/pocketix/pocketix-go/src/services"
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
	switch aTyped := a.(type) {
	case float64:
		bTyped, ok := b.(float64)
		if !ok {
			return false, fmt.Errorf("type mismatch: %T and %T", a, b)
		}
		return aTyped + bTyped, nil
	case int:
		bTyped, ok := b.(int)
		if !ok {
			return false, fmt.Errorf("type mismatch: %T and %T", a, b)
		}
		return float64(aTyped) + float64(bTyped), nil
	case string:
		bTyped, ok := b.(string)
		if !ok {
			return false, fmt.Errorf("type mismatch: %T and %T", a, b)
		}
		return aTyped + bTyped, nil
	default:
		return false, fmt.Errorf("unsupported type for + operator: %T", a)
	}
}

func CompareValues(a, b any, comparator func(x, y float64) bool) (bool, error) {
	switch a := a.(type) {
	case bool:
		b, ok := b.(bool)
		if !ok {
			return false, fmt.Errorf("type mismatch: %T and %T", a, b)
		}
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
		b, ok := b.(float64)
		if !ok {
			return false, fmt.Errorf("type mismatch: %T and %T", a, b)
		}
		return comparator(a, b), nil
	case int:
		b, ok := b.(int)
		if !ok {
			return false, fmt.Errorf("type mismatch: %T and %T", a, b)
		}
		return comparator(float64(a), float64(b)), nil
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

func (o *OperatorFactory) EvaluateOperator(operator string, child TreeNode, variableStore *VariableStore) (any, error) {
	if len(child.Children) == 0 {
		if child.Type == "variable" {
			// return child.ResultValue, -1, nil
			variable, err := variableStore.GetVariable(child.Value.(string))
			if err != nil {
				return nil, err
			}
			if variable.Value.ResultValue == nil {
				return variable.Value.Value, nil
			}
			return variable.Value.ResultValue, nil
		}
		return child.Value, nil
	}

	opFunc, exists := o.operator[operator]
	if !exists {
		return nil, fmt.Errorf("operator not supported: %s", operator)
	}

	expressionResult := child.Children[0].ResultValue
	var err error

	for i := range len(child.Children) - 1 {
		b := child.Children[i+1].ResultValue

		services.Logger.Println("Comparing", expressionResult, b, "with operator", operator)

		expressionResult, err = opFunc(expressionResult, b)
		if err != nil {
			return nil, err
		}
	}

	return expressionResult, nil
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
