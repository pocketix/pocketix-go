package models

import (
	"fmt"
	"reflect"
)

type VariableStore struct {
	Variables map[string]Variable
}

func NewVariableStore() *VariableStore {
	return &VariableStore{
		Variables: make(map[string]Variable),
	}
}

func (vs *VariableStore) AddVariable(variable Variable) {
	vs.Variables[variable.Name] = variable
}

func (vs *VariableStore) GetVariable(name string) (*Variable, error) {
	variable, ok := vs.Variables[name]
	if !ok {
		return nil, fmt.Errorf("variable %s not found", name)
	}
	return &variable, nil
}

func (vs *VariableStore) SetVariable(name string, value any, valueType string) error {
	variable, err := vs.GetVariable(name)
	if err != nil {
		return err
	}

	if valueType == "variable" {
		RValVariable, err := vs.GetVariable(value.(string))
		if err != nil {
			return err
		}
		result, err := RValVariable.Value.Evaluate(vs)
		if err != nil {
			return err
		}
		value = result
	}

	if reflect.TypeOf(value).String() != reflect.TypeOf(variable.Value.Value).String() {
		return fmt.Errorf("type mismatch at SetVariable")
	}

	variable.Value.Value = value
	vs.Variables[name] = *variable
	return nil
}
