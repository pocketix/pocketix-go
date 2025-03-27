package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
)

type Command interface {
	Execute(variableStore *models.VariableStore) (bool, error)
	GetId() string
	GetBody() []Command
	GetArguments() *models.TreeNode
}
