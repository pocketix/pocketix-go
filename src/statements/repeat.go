package statements

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type Repeat struct {
	Id        string
	Count     any
	CountType string
	Block     []Statement
}

func (r *Repeat) Execute(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, deviceCommands []models.SDInformationFromBackend) (any, bool, error) {
	services.Logger.Println("Executing repeat")

	var count int
	switch r.Count.(type) {
	case float64:
		count = int(r.Count.(float64))
	case int:
		count = r.Count.(int)
	case string:
		variable, err := variableStore.GetVariable(r.Count.(string))
		if err != nil {
			return nil, false, err
		}
		count = int(variable.Value.Value.(float64))
	}

	if count < 0 {
		return nil, false, fmt.Errorf("count cannot be negative")
	}

	for range count {
		if _, result, err := ExecuteStatements(r.Block, variableStore, referencedValueStore, deviceCommands); err != nil {
			return nil, result, err
		}
	}
	return r, true, nil
}

func (r *Repeat) GetId() string {
	return r.Id
}

func (r *Repeat) GetBody() []Statement {
	return r.Block
}

func (r *Repeat) GetCount() any {
	return r.Count
}

func (r *Repeat) GetCountType() string {
	return r.CountType
}

func (r *Repeat) Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error {
	if r.CountType == "variable" {
		variable, err := variableStore.GetVariable(r.Count.(string))
		if err != nil {
			return err
		}
		if variable.Type != "number" {
			return fmt.Errorf("count variable must be of type number")
		} else {
			return nil
		}
	}
	if r.CountType != "number" {
		return fmt.Errorf("count type must be number or variable")
	}
	return nil
}
