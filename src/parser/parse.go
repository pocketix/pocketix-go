package parser

import (
	"encoding/json"
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

func CheckMissingBlock(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if _, ok := raw["block"]; !ok {
		return fmt.Errorf("at least one block is required")
	}
	return nil
}

func Parse(data []byte) (*models.Program, error) {
	var program models.Program
	if err := CheckMissingBlock(data); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &program); err != nil {
		return nil, err
	}

	var previousCommand models.Command
	for _, block := range program.Blocks {
		cmd, err := ParseBlocks(block)

		if cmd.GetId() == "if" {
			previousCommand = cmd
		} else if cmd.GetId() == "else" {
			if previousCommand != nil {
				previousCommand.(*models.If).AddElseBlock(cmd)
				previousCommand.Execute()
				previousCommand = nil
			} else {
				services.Logger.Println("Error: Else without if")
				return nil, fmt.Errorf("else without if")
			}
		} else if cmd.GetId() == "elseif" {
			if previousCommand != nil {
				previousCommand.(*models.If).AddElseIfBlock(cmd)
			} else {
				services.Logger.Println("Error: Elseif without if")
				return nil, fmt.Errorf("elseif without if")
			}
		} else {
			if previousCommand != nil {
				previousCommand.Execute()
				previousCommand = nil
			}

			cmd.Execute()
		}

		if err != nil {
			return nil, err
		}
		// cmd.Execute()
	}

	if previousCommand != nil {
		previousCommand.Execute()
	}

	return &program, nil
}
