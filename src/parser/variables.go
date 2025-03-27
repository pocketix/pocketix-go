package parser

import (
	"encoding/json"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

func ParseVariables(data json.RawMessage, variableStore *models.VariableStore) {
	var variables map[string]any

	if err := json.Unmarshal(data, &variables); err != nil {
		services.Logger.Println("Error parsing variables", err)
	}

	for varName, varData := range variables {
		varType, varValue := varData.(map[string]any)["type"], varData.(map[string]any)["value"]

		// if varList, ok := varValue.([]any); ok {
		// 	services.Logger.Println("Parsing list of variables", varList)
		// }

		variableStore.AddVariable(models.Variable{
			Name:  varName,
			Type:  varType.(string),
			Value: models.InitTree(varType.(string), varType, varValue, variableStore),
		})
	}
}
