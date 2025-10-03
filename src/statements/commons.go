package statements

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/types"
)

func ExecuteStatements(
	statements []Statement,
	variableStore *models.VariableStore,
	referencedValueStore *models.ReferencedValueStore,
	deviceCommands []types.SDInformationFromBackend,
	callback func(deviceCommand types.SDCommandInvocation),
) (bool, error) {
	for _, statement := range statements {
		if _, err := statement.Execute(variableStore, referencedValueStore, deviceCommands, callback); err != nil {
			return false, err
		}
	}
	return true, nil
}
