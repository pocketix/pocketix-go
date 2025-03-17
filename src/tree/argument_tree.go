package tree

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/utils"
)

type TreeNode struct {
	Type     string      // Type of the argument
	Value    any         // Value of the argument
	Children []*TreeNode // Children of the argument
}

func InitTree(argumentType string, args any, variableStore *models.VariableStore) *TreeNode {
	t := TreeNode{}
	t.Value = argumentType
	t.Children = t.ParseChildren(args, variableStore)
	return &t
}

func (a *TreeNode) ParseChildren(args any, variableStore *models.VariableStore) []*TreeNode {
	services.Logger.Println("Parsing children", args)

	argList, ok := args.([]any)
	if !ok {
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
			children = append(children, &TreeNode{Value: argValue, Type: argType})
		}
	}
	return children
}

func (a *TreeNode) AddChild(child *TreeNode) {
	a.Children = append(a.Children, child)
}

func (a *TreeNode) Evaluate(variableStore *models.VariableStore) (bool, error) {
	operatorFactory := NewOperatorFactory()
	result, err, _ := a.EvaluateWithFactory(operatorFactory, variableStore)
	return utils.ToBool(result), err
}

func (a *TreeNode) EvaluateWithFactory(factory *OperatorFactory, variableStore *models.VariableStore) (any, error, bool) {
	for _, child := range a.Children {
		services.Logger.Println("Argument executing", child.Value)
		if result, err, ok := child.CheckGrandChildren(child, factory, variableStore); ok {
			if err != nil {
				services.Logger.Println("Error executing argument ", a.Value)
				return nil, err, false
			}
			if len(child.Children) == 0 {
				return result, nil, true
			}
			// factoryResult, factoryErr := factory.EvaluateOperator(child.Value.(string), *child, variableStore)
			return result, err, true
		}
	}
	return false, fmt.Errorf("error executing argument"), false
}

func (a *TreeNode) CheckGrandChildren(child *TreeNode, factory *OperatorFactory, variableStore *models.VariableStore) (any, error, bool) {
	for _, grandChild := range child.Children {
		if len(grandChild.Children) != 0 {
			if _, err, ok := grandChild.CheckGrandChildren(grandChild, factory, variableStore); ok {
				return grandChild.Value, err, true
			} else if err != nil {
				services.Logger.Println("Error executing argument ", a.Value)
				return nil, err, false
			}
		}
	}
	services.Logger.Println("This subtree can be evaluated")

	if len(child.Children) == 0 {
		return child.Value, nil, true
	}

	factoryResult, factoryErr := factory.EvaluateOperator(child.Value.(string), *child, variableStore)
	return factoryResult, factoryErr, true
}

func GetValue(arg any) any {
	return arg.(map[string]any)["value"]
}

func GetType(arg any) string {
	return arg.(map[string]any)["type"].(string)
}
