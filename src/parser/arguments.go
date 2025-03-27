package parser

import (
	"encoding/json"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/types"
)

func ParseArguments(rawArguments []types.Argument, argumentTree []*models.TreeNode, variableStore *models.VariableStore) error {
	var args any

	for i, arg := range rawArguments {
		if err := json.Unmarshal(arg.Value, &args); err != nil {
			return err
		}
		argumentTree[i] = models.InitTree(arg.Type, arg.Value, args, variableStore)
	}
	return nil
}
