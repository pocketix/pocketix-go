package tree

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/utils"
)

type TreeNode struct {
	Type        string      // Type of the argument
	Value       any         // Value of the argument
	Children    []*TreeNode // Children of the argument
	ResultValue any         // Result of the expression
}

func InitTree(argumentType string, argumentValue any, args any, variableStore *models.VariableStore) *TreeNode {
	t := TreeNode{}
	t.Type = argumentType
	// t.Value = argumentValue
	t.Children = t.ParseChildren(args, variableStore)
	return &t
}

func (a *TreeNode) ParseChildren(args any, variableStore *models.VariableStore) []*TreeNode {
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
			children = append(children, &TreeNode{Value: argValue, Type: argType, ResultValue: argValue})
		}
	}
	return children
}

func (a *TreeNode) AddChild(child *TreeNode) {
	a.Children = append(a.Children, child)
}

func (a *TreeNode) Evaluate(variableStore *models.VariableStore) (bool, error) {
	operatorFactory := NewOperatorFactory()
	result, err, _ := a.EvaluateNode(operatorFactory, variableStore)
	return utils.ToBool(result), err
}

func (a *TreeNode) EvaluateNode(factory *OperatorFactory, variableStore *models.VariableStore) (any, error, bool) {
	if len(a.Children) == 0 {
		return a.Value, nil, true
	}

	if len(a.Children) == 1 {
		return a.Children[0].EvaluateNode(factory, variableStore)
	}

	evaluatedChildren := make([]any, 0, len(a.Children))

	for _, child := range a.Children {
		services.Logger.Println("Evaluating child", child.Value)
		result, err, ok := child.EvaluateNode(factory, variableStore)
		if err != nil {
			services.Logger.Println("Error executing argument", a.Value)
			return nil, err, false
		}
		if ok {
			evaluatedChildren = append(evaluatedChildren, result)
		}
	}

	if len(evaluatedChildren) > 0 {
		factoryResult, factoryErr := factory.EvaluateOperator(a.Value.(string), *a, variableStore)
		a.ResultValue = factoryResult
		return factoryResult, factoryErr, factoryErr == nil
	}

	return nil, fmt.Errorf("error executing argument"), false
}

func GetValue(arg any) any {
	return arg.(map[string]any)["value"]
}

func GetType(arg any) string {
	return arg.(map[string]any)["type"].(string)
}
