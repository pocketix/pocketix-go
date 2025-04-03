package parser

import (
	"encoding/json"
	"fmt"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/types"
)

func CheckMissingBlock(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if _, ok := raw["header"]; !ok {
		return fmt.Errorf("header is required")
	}

	if _, ok := raw["block"]; !ok {
		return fmt.Errorf("at least one block is required")
	}
	return nil
}

func ParseWithoutExecuting(data []byte, variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore) (*types.Program, error) {
	var program types.Program

	if err := CheckMissingBlock(data); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &program); err != nil {
		return nil, err
	}

	if err := ParseVariables(program.Header.Variables, variableStore, referenceValueStore); err != nil {
		return nil, err
	}

	var previousCommand commands.Command
	for _, block := range program.Blocks {
		cmd, err := ParseBlockWithoutExecuting(block, variableStore, referenceValueStore)
		if err != nil {
			return nil, err
		}
		if cmd.GetId() == "if" {
			previousCommand = cmd
		} else if cmd.GetId() == "else" {
			if previousCommand != nil {
				previousCommand.(*commands.If).AddElseBlock(cmd)
				previousCommand = nil
			} else {
				services.Logger.Println("Error: Else without if")
				return nil, fmt.Errorf("else without if")
			}
		} else if cmd.GetId() == "elseif" {
			if previousCommand != nil {
				previousCommand.(*commands.If).AddElseIfBlock(cmd)
			} else {
				services.Logger.Println("Error: Elseif without if")
				return nil, fmt.Errorf("elseif without if")
			}
		} else {
			if previousCommand != nil {
				previousCommand = nil
			}
		}
	}

	return &program, nil
}

func Parse(data []byte, variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore) (*types.Program, error) {
	var program types.Program

	if err := CheckMissingBlock(data); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &program); err != nil {
		return nil, err
	}

	if err := ParseVariables(program.Header.Variables, variableStore, referenceValueStore); err != nil {
		return nil, err
	}

	var previousCommand commands.Command
	for _, block := range program.Blocks {
		cmd, err := ParseBlocks(block, variableStore, referenceValueStore)
		if err != nil {
			return nil, err
		}

		if cmd.GetId() == "if" {
			previousCommand = cmd
		} else if cmd.GetId() == "else" {
			if previousCommand != nil {
				previousCommand.(*commands.If).AddElseBlock(cmd)
				_, err := previousCommand.Execute(variableStore, referenceValueStore)
				if err != nil {
					return nil, err
				}
				previousCommand = nil
			} else {
				services.Logger.Println("Error: Else without if")
				return nil, fmt.Errorf("else without if")
			}
		} else if cmd.GetId() == "elseif" {
			if previousCommand != nil {
				previousCommand.(*commands.If).AddElseIfBlock(cmd)
			} else {
				services.Logger.Println("Error: Elseif without if")
				return nil, fmt.Errorf("elseif without if")
			}
		} else {
			if previousCommand != nil {
				_, err := previousCommand.Execute(variableStore, referenceValueStore)
				if err != nil {
					return nil, err
				}
				previousCommand = nil
			}

			_, err := cmd.Execute(variableStore, referenceValueStore)
			if err != nil {
				return nil, err
			}
		}
	}

	if previousCommand != nil {
		_, err := previousCommand.Execute(variableStore, referenceValueStore)
		if err != nil {
			return nil, err
		}
	}

	return &program, nil
}
