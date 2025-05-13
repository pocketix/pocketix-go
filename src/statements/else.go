package statements

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type Else struct {
	Id    string
	Block []Statement
}

func (e *Else) Execute(
	variableStore *models.VariableStore,
	referencedValueStore *models.ReferencedValueStore,
	deviceCommands []models.SDInformationFromBackend,
	callback func(deviceCommand models.SDCommandInvocation),
) (bool, error) {
	services.Logger.Println("Executing else")
	return ExecuteStatements(e.Block, variableStore, referencedValueStore, deviceCommands, callback)
}

func (e *Else) GetId() string {
	return e.Id
}

func (e *Else) GetBody() []Statement {
	return e.Block
}

func (e *Else) Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
