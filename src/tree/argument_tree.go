package tree

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/services"
)

type TreeNode struct {
	Value    any
	Children []*TreeNode
}

func InitTree(argumentType string, args any) *TreeNode {
	t := TreeNode{}
	t.Value = argumentType
	t.Children = t.ParseChildren(args)
	return &t
}

func (a *TreeNode) ParseChildren(args any) []*TreeNode {
	services.Logger.Println("Parsing children", args)
	var children []*TreeNode

	for _, arg := range args.([]any) {
		argType := GetType(arg)

		argValue := GetValue(arg)

		if value, ok := argValue.([]any); ok {
			services.Logger.Println("Argument is a list of values:", value)
			children = append(children, &TreeNode{Value: argType})
			children[len(children)-1].Children = children[len(children)-1].ParseChildren(value)
		} else {
			services.Logger.Println("Argument is a single value:", argValue, "of type:", argType)
			children = append(children, &TreeNode{Value: argValue})
		}
	}

	return children
}

func (a *TreeNode) AddChild(child *TreeNode) {
	a.Children = append(a.Children, child)
}

func (a *TreeNode) Evaluate() (any, error, bool) {
	operatorFactory := NewOperatorFactory()
	return a.EvaluateWithFactory(operatorFactory)
}

func (a *TreeNode) EvaluateWithFactory(factory *OperatorFactory) (any, error, bool) {
	for _, child := range a.Children {
		services.Logger.Println("Argument executing", child.Value)
		if result, err, ok := child.CheckGrandChildren(child, factory); ok {
			if err != nil {
				services.Logger.Println("Error executing argument ", a.Value)
				return nil, err, false
			}
			if len(child.Children) == 0 {
				return result, nil, true
			}
			factoryResult, factoryErr := factory.EvaluateOperator(child.Value.(string), *child)
			return factoryResult, factoryErr, true
		}
	}
	return false, fmt.Errorf("error executing argument"), false
}

func (a *TreeNode) CheckGrandChildren(child *TreeNode, factory *OperatorFactory) (any, error, bool) {
	for _, grandChild := range child.Children {
		if len(grandChild.Children) != 0 {
			if _, err, ok := grandChild.CheckGrandChildren(grandChild, factory); ok {
				return grandChild.Value, err, true
			} else if err != nil {
				services.Logger.Println("Error executing argument ", a.Value)
				return nil, err, false
			}
		}
	}
	services.Logger.Println("This subtree can be evaluated")
	factoryResult, factoryErr := factory.EvaluateOperator(child.Value.(string), *child)
	child.Value = factoryResult
	child.RemoveChildren()
	return factoryResult, factoryErr, true
}

func (a *TreeNode) RemoveChildren() {
	a.Children = nil
}

func GetValue(arg any) any {
	return arg.(map[string]any)["value"]
}

func GetType(arg any) string {
	return arg.(map[string]any)["type"].(string)
}
