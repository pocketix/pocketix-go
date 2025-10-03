package statements

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/types"
)

type Statement interface {
	Execute(
		variableStore *models.VariableStore,
		referencedValueStore *models.ReferencedValueStore,
		deviceCommands []types.SDInformationFromBackend,
		callback func(deviceCommand types.SDCommandInvocation),
	) (bool, error)
	GetId() string
	Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error
}

type BodyStatement interface {
	GetBody() []Statement
}

type ArgumentsStatement interface {
	GetArguments() *models.TreeNode
}
