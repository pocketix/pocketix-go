package models

import "fmt"

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

func (vs *VariableStore) GetVariable(name string) (any, error) {
	if variable, ok := vs.Variables[name]; ok {
		return variable.Value, nil
	}
	return nil, fmt.Errorf("variable %s not found", name)
}

func (vs *VariableStore) SetVariable(name string, value any) error {
	variable, ok := vs.Variables[name]
	if !ok {
		return fmt.Errorf("variable %s not found", name)
	}

	variable.Value = value
	vs.Variables[name] = variable
	return nil
}
