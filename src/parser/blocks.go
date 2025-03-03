package parser

import (
	"encoding/json"
	"fmt"

	"github.com/pocketix/pocketix-go/src/interfaces"
)

func ParseBlocks(data []json.RawMessage) []interfaces.Command {
	fmt.Println("Parsing blocks...")
	var commands []interfaces.Command
	for _, block := range data {
		cmd := ParseCommand(block)
		commands = append(commands, cmd)
	}
	return commands
}
