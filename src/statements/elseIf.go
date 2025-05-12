package statements

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/utils"
)

type ElseIf struct {
	Id        string
	Block     []Statement
	Arguments *models.TreeNode
}

func (e *ElseIf) Execute(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, deviceCommands []models.SDInformationFromBackend) (any, bool, error) {
	services.Logger.Println("Executing else if")
	if result, err := e.Arguments.Evaluate(variableStore, referencedValueStore); err != nil {
		services.Logger.Println("Error executing else if arguments", err)
	} else {
		if boolResult, boolErr := utils.ToBool(result); boolErr != nil {
			services.Logger.Println("Error converting else if result to bool", boolErr)
			return nil, false, boolErr
		} else if boolResult {
			services.Logger.Println("Else if is true, can execute body")
			return ExecuteStatements(e.Block, variableStore, referencedValueStore, deviceCommands)
		}
	}
	return e, false, nil
}

func (e *ElseIf) GetId() string {
	return e.Id
}

func (e *ElseIf) GetBody() []Statement {
	return e.Block
}

func (e *ElseIf) GetArguments() *models.TreeNode {
	return e.Arguments
}

func (e *ElseIf) Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
