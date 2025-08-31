package parser

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/types"
)

func ParseVariables(data json.RawMessage, variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore) error {
	var variables map[string]any
	argTypes := []string{"string", "number", "boolean", "variable", "boolean_expression", "str_opt"}

	if err := json.Unmarshal(data, &variables); err != nil {
		services.Logger.Println("Error parsing variables", err)
	}

	// type Argument struct {
	// 	Name  string
	// 	Type  string
	// 	Value any
	// }

	expressionVariables := make([]types.Argument, 0)

	for varName, varData := range variables {
		varType, varValue := varData.(map[string]any)["type"], varData.(map[string]any)["value"]

		if !slices.Contains(argTypes, varType.(string)) {
			return fmt.Errorf("argument type %s is not supported", varType)
		}

		if varType == "boolean_expression" || varType == "variable" {
			expressionVariables = append(expressionVariables, types.Argument{
				Reference: varName,
				Type:      varType.(string),
				Value:     varValue.(json.RawMessage),
			})
			continue
		}

		varData := types.Argument{
			Reference: varName,
			Type:      varType.(string),
			Value:     varValue.(json.RawMessage),
		}
		tree, err := models.InitTree(varData, variableStore, referencedValueStore)
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
		tree, err := models.InitTree(expressionVariable, variableStore, referencedValueStore)
		if err != nil {
			return err
		}

		variableStore.AddVariable(models.Variable{
			Name:  expressionVariable.Reference,
			Type:  expressionVariable.Type,
			Value: tree,
		})
	}

	return nil
}
