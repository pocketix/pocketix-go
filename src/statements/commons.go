package statements

import "github.com/pocketix/pocketix-go/src/models"

func ExecuteStatements(
	statements []Statement,
	variableStore *models.VariableStore,
	referencedValueStore *models.ReferencedValueStore,
	deviceCommands []models.SDInformationFromBackend,
	callback func(deviceCommand models.SDCommandInvocation),
) (bool, error) {
	for _, statement := range statements {
		if _, err := statement.Execute(variableStore, referencedValueStore, deviceCommands, callback); err != nil {
			return false, err
		}
	}
	return true, nil
}
