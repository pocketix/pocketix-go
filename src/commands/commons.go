package commands

import "github.com/pocketix/pocketix-go/src/models"

func ExecuteCommands(commands []Command, variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore) (bool, error) {
	for _, cmd := range commands {
		if success, err := cmd.Execute(variableStore, referenceValueStore); err != nil {
			return false, err
		} else if success {
			return success, nil
		}
	}
	return true, nil
}
