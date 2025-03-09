package parser

import (
	"encoding/json"

	"github.com/pocketix/pocketix-go/src/models"
)

func ParseArguments(rawArguments models.Argument) (any, error) {
	var args any

	if err := json.Unmarshal(rawArguments.Value, &args); err != nil {
		return nil, err
	}
	return args, nil
}
