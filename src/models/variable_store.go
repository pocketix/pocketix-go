package models

import "fmt"

type VariableStore struct {
	Variables []Variable
}

func NewVariableStore() *VariableStore {
	return &VariableStore{
		Variables: []Variable{},
	}
}

func (vs *VariableStore) AddVariable(variable Variable) {
	vs.Variables = append(vs.Variables, variable)
}

func (vs *VariableStore) GetVariable(name string) (any, error) {
	for _, variable := range vs.Variables {
		if variable.Name == name {
			return variable.Value, nil
		}
	}
	return nil, fmt.Errorf("variable %s not found", name)
}
