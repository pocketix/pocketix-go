package factories

import (
	"github.com/pocketix/pocketix-go/src/interfaces"
	"github.com/pocketix/pocketix-go/src/models"
)

func CommandFactory(id string, blocks []interfaces.Command, arguments []models.Argument) (interfaces.Command, error) {
	switch id {
	case "if":
		return &models.If{Id: id, Block: blocks, Arguments: arguments}, nil
	}
	return nil, nil
}
