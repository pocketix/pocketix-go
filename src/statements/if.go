package statements

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/utils"
)

type If struct {
	Id           string
	Block        []Statement
	Arguments    *models.TreeNode
	IfElseBlocks []ElseIf
	ElseBlock    Else
}

func (i *If) Execute(variableStore *models.VariableStore, commandHandlingStore *models.CommandsHandlingStore) (bool, error) {
	services.Logger.Println("Executing if")
	result, err := i.Arguments.Evaluate(variableStore, commandHandlingStore.ReferencedValueStore)
	if err != nil {
		services.Logger.Println("Error executing if arguments", err)
		return false, err
	}
	if boolResult, boolErr := utils.ToBool(result); boolErr != nil {
		services.Logger.Println("Error converting if result to bool", boolErr)
		return false, boolErr
	} else if boolResult {
		services.Logger.Println("If is true, can execute body")
		return ExecuteCommands(i.Block, variableStore, commandHandlingStore)
	}

	for i, elseIfBlock := range i.IfElseBlocks {
		if success, err := elseIfBlock.Execute(variableStore, commandHandlingStore); err != nil {
			return success, err
		} else if success {
			services.Logger.Println("Else if block", i, "executed successfully")
			return success, nil
		}
	}

	if i.ElseBlock.Id != "" {
		return i.ElseBlock.Execute(variableStore, commandHandlingStore)
	}
	return false, nil
}

func (i *If) GetId() string {
	return i.Id
}

func (i *If) GetBody() []Statement {
	return i.Block
}

func (i *If) GetArguments() *models.TreeNode {
	return i.Arguments
}

func (i *If) AddElseBlock(block Statement) {
	if elseBlock, ok := block.(*Else); ok {
		i.ElseBlock = *elseBlock
	}
}

func (i *If) AddElseIfBlock(block Statement) {
	if elseIfBlock, ok := block.(*ElseIf); ok {
		i.IfElseBlocks = append(i.IfElseBlocks, *elseIfBlock)
	}
}

func (i *If) Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
