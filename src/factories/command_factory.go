package factories

import (
	"github.com/pocketix/pocketix-go/src/interfaces"
	"github.com/pocketix/pocketix-go/src/models"
)

func CommandFactory(id string, blocks []interfaces.Command, tree *models.TreeNode) (interfaces.Command, error) {
	switch id {
	case "if":
		return &models.If{Id: id, Block: blocks, Arguments: tree}, nil
	}
	return nil, nil
}
