package models

func ExecuteCommands(commands []Command, variableStore *VariableStore) (bool, error) {
	for _, cmd := range commands {
		if success, err := cmd.Execute(variableStore); err != nil {
			return false, err
		} else if success {
			return success, nil
		}
	}
	return true, nil
}
