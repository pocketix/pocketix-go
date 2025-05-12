package statements

import "github.com/pocketix/pocketix-go/src/models"

func ExecuteStatements(statements []Statement, variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, deviceCommands []models.SDInformationFromBackend) (any, bool, error) {
	for _, statement := range statements {
		if _, _, err := statement.Execute(variableStore, referencedValueStore, deviceCommands); err != nil {
			return nil, false, err
		}
	}
	return statements, true, nil
}
