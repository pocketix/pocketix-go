package statements

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/utils"
)

type While struct {
	Id        string
	Block     []Statement
	Arguments *models.TreeNode
}

func (w *While) Execute(variableStore *models.VariableStore, commandHandlingStore *models.CommandsHandlingStore) (bool, error) {
	services.Logger.Println("Executing while")

	for {
		result, err := w.Arguments.Evaluate(variableStore, commandHandlingStore.ReferencedValueStore)
		if err != nil {
			services.Logger.Println("Error executing while arguments", err)
			return false, err
		}
		if boolResult, boolErr := utils.ToBool(result); boolErr != nil {
			services.Logger.Println("Error converting while result to bool", boolErr)
			return false, boolErr
		} else if boolResult {
			services.Logger.Println("While is true, can execute body")
			if success, err := ExecuteCommands(w.Block, variableStore, commandHandlingStore); err != nil {
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

func (w *While) GetBody() []Statement {
	return w.Block
}

func (w *While) GetArguments() *models.TreeNode {
	return w.Arguments
}

func (w *While) Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
