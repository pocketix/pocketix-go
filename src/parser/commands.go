package parser

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/pocketix/pocketix-go/src/factories"
	"github.com/pocketix/pocketix-go/src/interfaces"
)

func ParseCommand(data json.RawMessage) interfaces.Command {
	fmt.Println("Parsing command...")

	var raw struct {
		Type string `json:"id"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Println("Error parsing JSON at ParseCommand function: ", err)
		os.Exit(1)
	}
	if command, err := factories.CommandFactory(raw.Type, data); err != nil {
		fmt.Println("Error creating command: ", err)
		os.Exit(1)
	} else if command.HasBlock() {
		if statement, ok := command.(interfaces.Statement); ok {
			fmt.Println("Setting statement...")
			statement.SetBlock(ParseBlocks(GetRawJSON(data)))
			return command
		} else {
			log.Fatalf("command %s does not implement interfaces.Statement", raw.Type)
		}
	}
	return nil
}

func GetRawJSON(data json.RawMessage) []json.RawMessage {
	var raw struct {
		Blocks []json.RawMessage `json:"block"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Println("Error parsing JSON at GetRawJSON function: ", err)
		os.Exit(1)
	}
	return raw.Blocks
}
