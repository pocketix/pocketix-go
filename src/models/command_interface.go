package models

import "github.com/pocketix/pocketix-go/src/tree"

type Command interface {
	Execute() error
	GetId() string
	GetBody() []Command
	GetArguments() *tree.TreeNode
}
