package models

import (
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/tree"
)

type If struct {
	Id           string
	Block        []Command
	Arguments    *tree.TreeNode
	IfElseBlocks []ElseIf
	ElseBlock    Else
}

func (i *If) Execute() (bool, error) {
	services.Logger.Println("Executing if")
	result, err := i.Arguments.Evaluate()
	if err != nil {
		services.Logger.Println("Error executing if arguments", err)
		return false, err
	}
	if result {
		services.Logger.Println("If is true, can execute body")
		return ExecuteCommands(i.Block)
	}

	for i, elseIfBlock := range i.IfElseBlocks {
		if success, err := elseIfBlock.Execute(); err != nil {
			return success, err
		} else if success {
			services.Logger.Println("Else if block", i, "executed successfully")
			return success, nil
		}
	}

	if len(i.ElseBlock.Block) != 0 {
		return i.ElseBlock.Execute()
	}
	return false, nil
}

func ExecuteCommands(commands []Command) (bool, error) {
	for _, cmd := range commands {
		if success, err := cmd.Execute(); err != nil {
			return false, err
		} else if success {
			return success, nil
		}
	}
	return true, nil
}

func (i *If) GetId() string {
	return i.Id
}

func (i *If) GetBody() []Command {
	return i.Block
}

func (i *If) GetArguments() *tree.TreeNode {
	return i.Arguments
}

func (i *If) AddElseBlock(block Command) {
	if elseBlock, ok := block.(*Else); ok {
		i.ElseBlock = *elseBlock
	}
}

func (i *If) AddElseIfBlock(block Command) {
	if elseIfBlock, ok := block.(*ElseIf); ok {
		i.IfElseBlocks = append(i.IfElseBlocks, *elseIfBlock)
	}
}
