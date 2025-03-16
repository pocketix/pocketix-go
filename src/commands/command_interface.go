package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/tree"
)

type Command interface {
	Execute(variableStore *models.VariableStore) (bool, error)
	GetId() string
	GetBody() []Command
	GetArguments() *tree.TreeNode
}
