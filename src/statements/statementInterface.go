package statements

import (
	"github.com/pocketix/pocketix-go/src/models"
)

type Statement interface {
	Execute(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, deviceCommands []models.SDInformationFromBackend) (any, bool, error)
	GetId() string
	Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error
}

type BodyStatement interface {
	GetBody() []Statement
}

type ArgumentsStatement interface {
	GetArguments() *models.TreeNode
}
