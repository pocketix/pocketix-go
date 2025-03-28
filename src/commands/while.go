package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/utils"
)

type While struct {
	Id        string
	Block     []Command
	Arguments *models.TreeNode
}

func (w *While) Execute(variableStore *models.VariableStore) (bool, error) {
	services.Logger.Println("Executing while")

	for {
		result, err := w.Arguments.Evaluate(variableStore)
		if err != nil {
			services.Logger.Println("Error executing while arguments", err)
			return false, err
		}
		if utils.ToBool(result) {
			services.Logger.Println("While is true, can execute body")
			if success, err := ExecuteCommands(w.Block, variableStore); err != nil {
				return success, err
			} else if !success {
				return success, nil
			}
			// variableStore.SetVariable("foo", "a") // Test setting of a variable
		} else {
			services.Logger.Println("While is false, breaking")
			break
		}
	}

	return true, nil
}

func (w *While) GetId() string {
	return w.Id
}

func (w *While) GetBody() []Command {
	return w.Block
}

func (w *While) GetArguments() *models.TreeNode {
	return w.Arguments
}
