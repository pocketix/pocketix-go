package models

import (
	"fmt"
	"slices"

	"github.com/pocketix/pocketix-go/src/services"
)

type TreeNode struct {
	Type        string      // Type of the argument
	Value       any         // Value of the argument
	Children    []*TreeNode // Children of the argument
	ResultValue any         // Result of the expression
}

var typeValidators = map[string]func(any) error{
	"string": func(v any) error {
		if _, ok := v.(string); !ok {
			return fmt.Errorf("expected string, got %T", v)
		}
		return nil
	},
	"number": func(v any) error {
		_, ok := v.(float64)
		_, ok2 := v.(int)
		if !ok && !ok2 {
			return fmt.Errorf("expected number, got %T", v)
		}
		return nil
	},
	"boolean": func(v any) error {
		if _, ok := v.(bool); !ok {
			return fmt.Errorf("expected boolean, got %T", v)
		}
		return nil
	},
}

func InitTree(argumentType string, argumentValue any, args any, variableStore *VariableStore, commandHandlingStore *CommandsHandlingStore) (*TreeNode, error) {
	t := TreeNode{}
	t.Type = argumentType

	factory := NewOperatorFactory()
	parsedChildren, err := t.ParseChildren(args, factory, variableStore, commandHandlingStore)
	if err != nil {
		return nil, err
	}

	t.Children = parsedChildren
	return &t, nil
}

func (a *TreeNode) ParseChildren(args any, operatorFactory *OperatorFactory, variableStore *VariableStore, commandHandlingStore *CommandsHandlingStore) ([]*TreeNode, error) {
	services.Logger.Println("Parsing children", args)

	argList, ok := args.([]any)
	if !ok {
		if err := ValidateType(a.Type, args); err != nil {
			return nil, err
		}
		a.Value = args
		return nil, nil
	}

	children := make([]*TreeNode, 0, len(argList))

	for _, arg := range argList {
		argType := GetType(arg)
		argValue := GetValue(arg)

		if value, ok := argValue.([]any); ok {
			services.Logger.Println("Argument is a list of values:", value)

			child := &TreeNode{Value: argType}
			childrenList, err := child.ParseChildren(value, operatorFactory, variableStore, commandHandlingStore)
			if err != nil {
				return nil, err
			}
			child.Children = childrenList

			err = child.ValidateNode(operatorFactory)
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		} else {
			services.Logger.Println("Argument is a single value:", argValue, "of type:", argType)

			argTypes := []string{"string", "number", "boolean", "variable", "boolean_expression", "str_opt"} // insert "device_variable" type later
			if !slices.Contains(argTypes, argType) {
				return nil, fmt.Errorf("argument type %s is not supported", argType)
			}

			if err := ValidateType(argType, argValue); err != nil {
				return nil, err
			}

			if argType == "variable" {
				referencedValue, err := NewReferencedValue(argValue.(string))
				if err == nil {
					commandHandlingStore.ReferencedValueStore.AddReferencedValue(argValue.(string), referencedValue)
					children = append(children, &TreeNode{Value: argValue, Type: argType, ResultValue: referencedValue.Value})
				} else {
					services.Logger.Println("Error creating referenced value:", err)
					if variable, err := variableStore.GetVariable(argValue.(string)); err != nil {
						return nil, err
					} else {
						children = append(children, &TreeNode{Value: argValue, Type: argType, ResultValue: variable.Value.Value})
					}
				}
				// } else if argType == "device_variable" {}
			} else {
				children = append(children, &TreeNode{Value: argValue, Type: argType, ResultValue: argValue})
			}
		}
	}
	return children, nil
}

func ValidateType(argType string, argValue any) error {
	if validator, exists := typeValidators[argType]; exists {
		return validator(argValue)
	}
	return nil
}

func (a *TreeNode) ValidateNode(factory *OperatorFactory) error {
	// Check if children are empty
	for _, child := range a.Children {
		if len(child.Children) > 0 {
			return nil
		}
	}
	// This node's children are empty, therefore this node is leaf node and can validate it
	return factory.ValidateOperator(*a)
}

func (a *TreeNode) AddChild(child *TreeNode) {
	a.Children = append(a.Children, child)
}

func (a *TreeNode) Evaluate(variableStore *VariableStore, referenceValueStore *ReferencedValueStore) (any, error) {
	operatorFactory := NewOperatorFactory()
	result, err, _ := a.EvaluateNode(operatorFactory, variableStore, referenceValueStore)
	return result, err
}

func (a *TreeNode) EvaluateNode(factory *OperatorFactory, variableStore *VariableStore, referenceValueStore *ReferencedValueStore) (any, error, bool) {
	if len(a.Children) == 0 {
		return EvaluateArgumentsHelper(a, factory, variableStore, referenceValueStore)
	}

	if len(a.Children) == 1 {
		return a.Children[0].EvaluateNode(factory, variableStore, referenceValueStore)
	}

	evaluatedChildren := make([]any, 0, len(a.Children))

	for _, child := range a.Children {
		services.Logger.Println("Evaluating child", child.Value)
		result, err, ok := child.EvaluateNode(factory, variableStore, referenceValueStore)
		if err != nil {
			services.Logger.Println("Error executing argument", a.Value)
			return nil, err, false
		}
		if ok {
			evaluatedChildren = append(evaluatedChildren, result)
		}
	}

	if len(evaluatedChildren) > 0 {
		return EvaluateArgumentsHelper(a, factory, variableStore, referenceValueStore)
	}

	return nil, fmt.Errorf("error executing argument"), false
}

func EvaluateArgumentsHelper(node *TreeNode, factory *OperatorFactory, variableStore *VariableStore, referenceValueStore *ReferencedValueStore) (any, error, bool) {
	if node.Value == nil || (node.Type != "string" && node.Type != "" && node.Type != "variable") {
		return node.Value, nil, true
	}

	factoryResult, factoryErr := factory.EvaluateOperator(node.Value.(string), *node, variableStore, referenceValueStore)
	node.ResultValue = factoryResult
	return factoryResult, factoryErr, factoryErr == nil
}

func GetValue(arg any) any {
	return arg.(map[string]any)["value"]
}

func GetType(arg any) string {
	return arg.(map[string]any)["type"].(string)
}
