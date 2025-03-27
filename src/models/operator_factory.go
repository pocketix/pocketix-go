package models

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/services"
)

type OperatorFunction func(child TreeNode, variableStore *VariableStore) (any, error)

type OperatorFactory struct {
	operator map[string]func(a, b any) (any, error)
}

func NewOperatorFactory() *OperatorFactory {
	return &OperatorFactory{
		operator: map[string]func(a, b any) (any, error){
			"===": func(a, b any) (any, error) {
				return a == b, nil
			},
			"!==": func(a, b any) (any, error) {
				return a != b, nil
			},
			"<": func(a, b any) (any, error) {
				return CompareValues(a, b, func(x, y float64) bool { return x < y })
			},
			"+": func(a, b any) (any, error) {
				return AddValues(a, b)
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
		return comparator(float64(len(a)), float64(len(b))), nil
	default:
		return false, fmt.Errorf("unsupported type: %T", a)
	}
}

// TODO add more operators
// func NewOperatorFactory() *OperatorFactory {
// 	return &OperatorFactory{
// 		operator: map[string]OperatorFunction{
// 			"===": func(child TreeNode, variableStore *VariableStore) (any, error) {
// 				if len(child.Children) == 0 {
// 					return child.Value, nil
// 				}
// 				for i := range len(child.Children) - 1 {
// 					a := child.Children[i].ResultValue
// 					b := child.Children[i+1].ResultValue

// 					services.Logger.Println("Comparing", a, b)
// 					if a != b {
// 						return false, nil
// 					}
// 				}
// 				return true, nil
// 			},
// 			"!==": func(child TreeNode, variableStore *VariableStore) (any, error) {
// 				if len(child.Children) == 0 {
// 					return child.Value, nil
// 				}
// 				for i := range len(child.Children) - 1 {
// 					a := child.Children[i].ResultValue
// 					b := child.Children[i+1].ResultValue

// 					services.Logger.Println("Comparing", a, b)
// 					if a == b {
// 						return false, nil
// 					}
// 				}

// 				return true, nil
// 			},
// 			"<": func(child TreeNode, variableStore *VariableStore) (any, error) {
// 				if len(child.Children) == 0 {
// 					return child.Value, nil
// 				}
// 				for i := range len(child.Children) - 1 {
// 					a := child.Children[i].ResultValue
// 					b := child.Children[i+1].ResultValue

// 					services.Logger.Println("Comparing", a, b)
// 					if a >= b {
// 						return false, nil
// 					}
// 				}
// 				return true, nil
// 			},
// 			// ">=": func(child TreeNode) (any, error) {
// 			// 	if len(child.Children) == 0 {
// 			// 		return child.Value, nil
// 			// 	}
// 			// 	for i := range len(child.Children) - 1 {
// 			// 		services.Logger.Println("Comparing", child.Children[i].Value, child.Children[i+1].Value)
// 			// 		if child.Children[i].Value < child.Children[i+1].Value {
// 			// 			return false, nil
// 			// 		}
// 			// 	}
// 			// 	return true, nil
// 			// },
// 			// "<=": func(child TreeNode) (any, error) {
// 			// 	if len(child.Children) == 0 {
// 			// 		return child.Value, nil
// 			// 	}
// 			// 	for i := range len(child.Children) - 1 {
// 			// 		services.Logger.Println("Comparing", child.Children[i].Value, child.Children[i+1].Value)
// 			// 		if child.Children[i].Value > child.Children[i+1].Value {
// 			// 			return false, nil
// 			// 		}
// 			// 	}
// 			// 	return true, nil
// 			// },
// 			// ">": func(child TreeNode) (any, error) {
// 			// 	if len(child.Children) == 0 {
// 			// 		return child.Value, nil
// 			// 	}
// 			// 	for i := range len(child.Children) - 1 {
// 			// 		services.Logger.Println("Comparing", child.Children[i].Value, child.Children[i+1].Value)
// 			// 		if child.Children[i].Value <= child.Children[i+1].Value {
// 			// 			return false, nil
// 			// 		}
// 			// 	}
// 			// 	return true, nil
// 			// },
// 			// "<": func(child TreeNode) (any, error) {
// 			// 	if len(child.Children) == 0 {
// 			// 		return child.Value, nil
// 			// 	}
// 			// 	for i := range len(child.Children) - 1 {
// 			// 		services.Logger.Println("Comparing", child.Children[i].Value, child.Children[i+1].Value)
// 			// 		if child.Children[i].Value >= child.Children[i+1].Value {
// 			// 			return false, nil
// 			// 		}
// 			// 	}
// 			// 	return true, nil
// 			// },
// 		},
// 	}
// }

// func CheckOperatorsType(leftOperand models.TreeNode, rightOperand models.TreeNode) error {
// 	if leftOperand.Type == "string" && rightOperand.Type == "string" {

// }

func GetValueFromStore(node TreeNode, variableStore *VariableStore) (any, error) {
	if node.Type == "variable" {
		if variable, err := variableStore.GetVariable(node.Value.(string)); err != nil {
			return nil, err
		} else {
			return variable.Value, nil
		}
	}
	return node.ResultValue, nil
}

func (o *OperatorFactory) EvaluateOperator(operator string, child TreeNode, variableStore *VariableStore) (any, float64, error) {
	if len(child.Children) == 0 {
		if child.Type == "variable" {
			// return child.ResultValue, -1, nil
			variable, err := variableStore.GetVariable(child.Value.(string))
			if err != nil {
				return nil, -1, err
			}
			if variable.Value.ResultValue == nil {
				return variable.Value.Value, -1, nil
			}
			return variable.Value.ResultValue, -1, nil
		}
		return child.Value, -1, nil
	}

	opFunc, exists := o.operator[operator]
	if !exists {
		return nil, -1, fmt.Errorf("operator not supported: %s", operator)
	}

	var numericalResult float64
	for i := range len(child.Children) - 1 {
		a := child.Children[i].ResultValue
		b := child.Children[i+1].ResultValue

		// if errA != nil || errB != nil {
		// 	return nil, fmt.Errorf("error getting value from store: %v, %v", errA, errB)
		// }

		services.Logger.Println("Comparing", a, b, "with operator", operator)

		result, err := opFunc(a, b)
		if err != nil {
			return nil, -1, err
		}

		if boolResult, ok := result.(bool); ok {
			if !boolResult {
				return false, -1, nil
			}
		} else if numResult, ok := result.(float64); ok {
			numericalResult += numResult
		}
	}

	return true, numericalResult, nil
}
