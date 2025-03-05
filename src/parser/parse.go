package parser

import (
	"encoding/json"
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

func CheckMissingBlock(data []byte) error {
	var raw map[string]interface{}
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

	for _, block := range program.Blocks {
		cmd, err := ParseBlocks(block)
		if err != nil {
			return nil, err
		}
		// program.Blocks = append(program.Blocks, cmd)
		// log.Println("Command: ", cmd)
		services.Logger.Println("Command: ", cmd)
	}
	return &program, nil
}
