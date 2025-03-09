package models

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

func (a *TreeNode) Execute() (bool, error) {
	for _, child := range a.Children {
		services.Logger.Println("Argument executing", child.Value)

		if len(child.Children) == 0 {
			services.Logger.Println("Value is final argument")
			return true, nil
		}

		if result, err := child.Execute(); result {
			if child.Value == "===" {
				for i := range len(child.Children) - 1 {
					services.Logger.Println("Comparing", child.Children[i].Value, "with", child.Children[i+1].Value)
					if child.Children[i].Value != child.Children[i+1].Value {
						services.Logger.Println("value", child.Children[i].Value, "is not equal to value", child.Children[i+1].Value)
						return false, nil
					}
				}
				// child.Value = "def"
				return true, nil
			}
		} else if err != nil {
			return false, err
		}
	}
	return false, fmt.Errorf("error executing argument")
}

func (a *TreeNode) String() string {
	for _, child := range a.Children {
		services.Logger.Println(child)
	}
	return a.Value.(string)
}

func GetValue(arg any) any {
	return arg.(map[string]any)["value"]
}

func GetType(arg any) string {
	return arg.(map[string]any)["type"].(string)
}
