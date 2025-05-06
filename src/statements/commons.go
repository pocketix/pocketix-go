package statements

import "github.com/pocketix/pocketix-go/src/models"

func ExecuteStatements(statements []Statement, variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore) (bool, error) {
	for _, statement := range statements {
		if _, err := statement.Execute(variableStore, referencedValueStore); err != nil {
			return false, err
		}
	}
	return true, nil
}
