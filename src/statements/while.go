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

func (w *While) Execute(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, deviceCommands []models.SDInformationFromBackend) (any, bool, error) {
	services.Logger.Println("Executing while")

	for {
		result, err := w.Arguments.Evaluate(variableStore, referencedValueStore)
		if err != nil {
			services.Logger.Println("Error executing while arguments", err)
			return nil, false, err
		}
		if boolResult, boolErr := utils.ToBool(result); boolErr != nil {
			services.Logger.Println("Error converting while result to bool", boolErr)
			return nil, false, boolErr
		} else if boolResult {
			services.Logger.Println("While is true, can execute body")
			if _, success, err := ExecuteStatements(w.Block, variableStore, referencedValueStore, deviceCommands); err != nil {
				return nil, success, err
			}
		} else {
			services.Logger.Println("While is false, breaking")
			break
		}
	}

	return w, true, nil
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

func (w *While) Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
