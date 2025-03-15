package tree

import "github.com/pocketix/pocketix-go/src/services"

type OperatorFunction func(args TreeNode) (any, error)

type OperatorFactory struct {
	operator map[string]OperatorFunction
}

// TODO add more operators
func NewOperatorFactory() *OperatorFactory {
	return &OperatorFactory{
		operator: map[string]OperatorFunction{
			"===": func(child TreeNode) (any, error) {
				if len(child.Children) == 0 {
					return child.Value, nil
				}
				for i := range len(child.Children) - 1 {
					services.Logger.Println("Comparing", child.Children[i].Value, child.Children[i+1].Value)
					if child.Children[i].Value != child.Children[i+1].Value {
						return false, nil
					}
				}
				return true, nil
			},
			"!==": func(child TreeNode) (any, error) {
				if len(child.Children) == 0 {
					return child.Value, nil
				}
				for i := range len(child.Children) - 1 {
					services.Logger.Println("Comparing", child.Children[i].Value, child.Children[i+1].Value)
					if child.Children[i].Value == child.Children[i+1].Value {
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

func (o *OperatorFactory) EvaluateOperator(operator string, child TreeNode) (any, error) {
	if fn, exists := o.operator[operator]; exists {
		return fn(child)
	}
	services.Logger.Println("Operator not found", operator)
	return false, nil
}
