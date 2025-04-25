package statements

import (
	"github.com/pocketix/pocketix-go/src/models"
)

type Statement interface {
	Execute(variableStore *models.VariableStore, commandHandlingStore *models.CommandsHandlingStore) (bool, error)
	GetId() string
	GetBody() []Statement
	GetArguments() *models.TreeNode
	Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error
}
