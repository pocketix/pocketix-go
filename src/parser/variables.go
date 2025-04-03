package parser

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

func ParseVariables(data json.RawMessage, variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore) error {
	var variables map[string]any
	argTypes := []string{"string", "number", "boolean", "variable", "boolean_expression", "str_opt"}

	if err := json.Unmarshal(data, &variables); err != nil {
		services.Logger.Println("Error parsing variables", err)
	}

	type ExpressionVariable struct {
		Name  string
		Type  string
		Value any
	}

	expressionVariables := make([]ExpressionVariable, 0)

	for varName, varData := range variables {
		varType, varValue := varData.(map[string]any)["type"], varData.(map[string]any)["value"]

		if !slices.Contains(argTypes, varType.(string)) {
			return fmt.Errorf("argument type %s is not supported", varType)
		}

		if varType == "boolean_expression" || varType == "variable" {
			expressionVariables = append(expressionVariables, ExpressionVariable{
				Name:  varName,
				Type:  varType.(string),
				Value: varValue,
			})
			continue
		}
		tree, err := models.InitTree(varType.(string), varValue, varValue, variableStore, referenceValueStore)
		if err != nil {
			return err
		}

		variableStore.AddVariable(models.Variable{
			Name:  varName,
			Type:  varType.(string),
			Value: tree,
		})
	}

	for _, expressionVariable := range expressionVariables {
		tree, err := models.InitTree(expressionVariable.Type, expressionVariable.Value, expressionVariable.Value, variableStore, referenceValueStore)
		if err != nil {
			return err
		}

		variableStore.AddVariable(models.Variable{
			Name:  expressionVariable.Name,
			Type:  expressionVariable.Type,
			Value: tree,
		})
	}

	return nil
}
