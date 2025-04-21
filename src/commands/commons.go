package commands

import "github.com/pocketix/pocketix-go/src/models"

func ExecuteCommands(commands []Command, variableStore *models.VariableStore, commandHandlingStore *models.CommandsHandlingStore) (bool, error) {
	for _, cmd := range commands {
		if success, err := cmd.Execute(variableStore, commandHandlingStore); err != nil {
			return false, err
		} else if success {
			return success, nil
		}
	}
	return true, nil
}
