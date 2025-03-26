package tree

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type OperatorFunction func(child TreeNode, variableStore *models.VariableStore) (any, error)

type OperatorFactory struct {
	operator map[string]OperatorFunction
}

// TODO add more operators
func NewOperatorFactory() *OperatorFactory {
	return &OperatorFactory{
		operator: map[string]OperatorFunction{
			"===": func(child TreeNode, variableStore *models.VariableStore) (any, error) {
				if len(child.Children) == 0 {
					return child.Value, nil
				}
				for i := range len(child.Children) - 1 {
					a := child.Children[i].ResultValue
					b := child.Children[i+1].ResultValue

					services.Logger.Println("Comparing", a, b)
					if a != b {
						return false, nil
					}
				}
				return true, nil
			},
			"!==": func(child TreeNode, variableStore *models.VariableStore) (any, error) {
				if len(child.Children) == 0 {
					return child.Value, nil
				}
				for i := range len(child.Children) - 1 {
					a := child.Children[i].ResultValue
					b := child.Children[i+1].ResultValue

					services.Logger.Println("Comparing", a, b)
					if a == b {
						return false, nil
					}
				}

				return true, nil
			},
			// ">=": func(child TreeNode) (any, error) {
			// 	if len(child.Children) == 0 {
			// 		return child.Value, nil
			// 	}
			// 	for i := range len(child.Children) - 1 {
			// 		services.Logger.Println("Comparing", child.Children[i].Value, child.Children[i+1].Value)
			// 		if child.Children[i].Value < child.Children[i+1].Value {
			// 			return false, nil
			// 		}
			// 	}
			// 	return true, nil
			// },
			// "<=": func(child TreeNode) (any, error) {
			// 	if len(child.Children) == 0 {
			// 		return child.Value, nil
			// 	}
			// 	for i := range len(child.Children) - 1 {
			// 		services.Logger.Println("Comparing", child.Children[i].Value, child.Children[i+1].Value)
			// 		if child.Children[i].Value > child.Children[i+1].Value {
			// 			return false, nil
			// 		}
			// 	}
			// 	return true, nil
			// },
			// ">": func(child TreeNode) (any, error) {
			// 	if len(child.Children) == 0 {
			// 		return child.Value, nil
			// 	}
			// 	for i := range len(child.Children) - 1 {
			// 		services.Logger.Println("Comparing", child.Children[i].Value, child.Children[i+1].Value)
			// 		if child.Children[i].Value <= child.Children[i+1].Value {
			// 			return false, nil
			// 		}
			// 	}
			// 	return true, nil
			// },
			// "<": func(child TreeNode) (any, error) {
			// 	if len(child.Children) == 0 {
			// 		return child.Value, nil
			// 	}
			// 	for i := range len(child.Children) - 1 {
			// 		services.Logger.Println("Comparing", child.Children[i].Value, child.Children[i+1].Value)
			// 		if child.Children[i].Value >= child.Children[i+1].Value {
			// 			return false, nil
			// 		}
			// 	}
			// 	return true, nil
			// },
		},
	}
}

func GetValueFromStore(node TreeNode, variableStore *models.VariableStore) (any, error) {
	if node.Type == "variable" {
		if variable, err := variableStore.GetVariable(node.ResultValue.(string)); err != nil {
			return nil, err
		} else {
			return variable.Value, nil
		}
	}
	return node.ResultValue, nil
}

func (o *OperatorFactory) EvaluateOperator(operator string, child TreeNode, variableStore *models.VariableStore) (any, error) {
	if len(child.Children) == 0 {
		if child.Type == "variable" {
			if variable, err := variableStore.GetVariable(child.ResultValue.(string)); err != nil {
				return nil, err
			} else {
				return variable.Value, nil
			}
		}
		return child.Value, nil
	}

	if fn, exists := o.operator[operator]; exists {
		return fn(child, variableStore)
	}
	services.Logger.Println("Operator not found", operator)
	return false, fmt.Errorf("operator not found")
}
