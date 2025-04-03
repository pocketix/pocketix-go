package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
)

type Command interface {
	Execute(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore) (bool, error)
	GetId() string
	GetBody() []Command
	GetArguments() *models.TreeNode
	Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error
}
