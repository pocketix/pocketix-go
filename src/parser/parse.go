package parser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/pocketix/pocketix-go/src/models"
)

func Parse(data []byte) (*models.Program, error) {
	fmt.Println("Parsing JSON...")
	var program models.Program

	var raw struct {
		Blocks []json.RawMessage `json:"block"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		fmt.Println("Error parsing JSON at UnmarshalJSON function: ", err)
		os.Exit(1)
	}

	program.Blocks = ParseBlocks(raw.Blocks)
	return &program, nil
}
