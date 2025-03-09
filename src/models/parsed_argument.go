package models

import (
	"github.com/pocketix/pocketix-go/src/services"
)

type ParsedArgument struct {
	Type  string
	Value any
}

func (p *ParsedArgument) Execute() (bool, error) {

	if argumentValue, ok := p.Value.(ParsedArgument); ok {
		services.Logger.Println("ok Executing parsed argument: ", p.Type, argumentValue)
		argumentValue.Execute()
	} else {
		if value, ok := p.Value.([]ParsedArgument); ok {
			services.Logger.Println("Executing parsed argument: ", p.Type, value)
			value[0].Execute()
		} else {
			services.Logger.Println("Value is final argument")
		}
	}
	return true, nil
}

func GetOperator(arg any) string {
	return arg.(map[string]any)["operator"].(string)
}
