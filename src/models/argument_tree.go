package models

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/utils"
)

type TreeNode struct {
	Type        string      // Type of the argument
	Value       any         // Value of the argument
	Children    []*TreeNode // Children of the argument
	ResultValue any         // Result of the expression
}

func InitTree(argumentType string, argumentValue any, args any, variableStore *VariableStore) *TreeNode {
	t := TreeNode{}
	t.Type = argumentType
	// t.Value = argumentValue
	t.Children = t.ParseChildren(args, variableStore)
	return &t
}

func (a *TreeNode) ParseChildren(args any, variableStore *VariableStore) []*TreeNode {
	services.Logger.Println("Parsing children", args)

	argList, ok := args.([]any)
	if !ok {
		a.Value = args
		return nil
	}

	children := make([]*TreeNode, 0, len(argList))

	for _, arg := range argList {
		argType := GetType(arg)
		argValue := GetValue(arg)

		if value, ok := argValue.([]any); ok {
			services.Logger.Println("Argument is a list of values:", value)
			child := &TreeNode{Value: argType}
			child.Children = child.ParseChildren(value, variableStore)
			children = append(children, child)
		} else {
			services.Logger.Println("Argument is a single value:", argValue, "of type:", argType)
			if argType == "variable" {
				if variable, err := variableStore.GetVariable(argValue.(string)); err != nil {
					services.Logger.Println("Error getting variable", argValue)
				} else {
					children = append(children, &TreeNode{Value: argValue, Type: argType, ResultValue: variable.Value.Value})
				}
			} else {
				children = append(children, &TreeNode{Value: argValue, Type: argType, ResultValue: argValue})
			}
		}
	}
	return children
}

func (a *TreeNode) AddChild(child *TreeNode) {
	a.Children = append(a.Children, child)
}

func (a *TreeNode) Evaluate(variableStore *VariableStore) (bool, float64, error) {
	operatorFactory := NewOperatorFactory()
	result, numericalResult, err, _ := a.EvaluateNode(operatorFactory, variableStore)
	return utils.ToBool(result), numericalResult, err
}

func (a *TreeNode) EvaluateNode(factory *OperatorFactory, variableStore *VariableStore) (any, float64, error, bool) {
	if len(a.Children) == 0 {
		return EvaluateArgumentsHelper(a, factory, variableStore)
		// return a.Value, nil, true
	}

	if len(a.Children) == 1 {
		return a.Children[0].EvaluateNode(factory, variableStore)
	}

	evaluatedChildren := make([]any, 0, len(a.Children))

	for _, child := range a.Children {
		services.Logger.Println("Evaluating child", child.Value)
		result, _, err, ok := child.EvaluateNode(factory, variableStore)
		if err != nil {
			services.Logger.Println("Error executing argument", a.Value)
			return nil, -1, err, false
		}
		if ok {
			evaluatedChildren = append(evaluatedChildren, result)
		}
	}

	if len(evaluatedChildren) > 0 {
		return EvaluateArgumentsHelper(a, factory, variableStore)
	}

	return nil, -1, fmt.Errorf("error executing argument"), false
}

func EvaluateArgumentsHelper(node *TreeNode, factory *OperatorFactory, variableStore *VariableStore) (any, float64, error, bool) {
	if node.Value == nil || (node.Type != "string" && node.Type != "" && node.Type != "variable") {
		return node.Value, -1, nil, true
	}

	factoryResult, factoryNumericalResult, factoryErr := factory.EvaluateOperator(node.Value.(string), *node, variableStore)
	node.ResultValue = factoryResult
	return factoryResult, factoryNumericalResult, factoryErr, factoryErr == nil
}

func GetValue(arg any) any {
	return arg.(map[string]any)["value"]
}

func GetType(arg any) string {
	return arg.(map[string]any)["type"].(string)
}
