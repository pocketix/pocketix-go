package parser

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/types"
)

func ParseArguments(rawArguments []types.Argument, argumentTree []*models.TreeNode, variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore) error {
	var args any
	argTypes := []string{"string", "number", "boolean", "variable", "boolean_expression", "str_opt"}

	for i, arg := range rawArguments {
		if !slices.Contains(argTypes, arg.Type) {
			return fmt.Errorf("argument type %s is not supported", arg.Type)
		}
		if err := json.Unmarshal(arg.Value, &args); err != nil {
			return err
		}
		tree, err := models.InitTree(arg.Type, arg.Value, args, variableStore, referencedValueStore)
		if err != nil {
			return err
		}
		argumentTree[i] = tree
	}
	return nil
}
