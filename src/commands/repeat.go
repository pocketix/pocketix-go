package commands

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type Repeat struct {
	Id        string
	Count     any
	CountType string
	Block     []Command
}

func (r *Repeat) Execute(variableStore *models.VariableStore) (bool, error) {
	services.Logger.Println("Executing repeat")

	count := r.Count
	if r.CountType == "variable" {
		// variable, err := variableStore.GetVariable(count.(string))
		// if err != nil {
		// 	return false, err
		// }
		// count = variable.Value.(int)
	}

	if count.(int) < 0 {
		return false, fmt.Errorf("count cannot be negative")
	}

	for range count.(int) {
		if result, err := ExecuteCommands(r.Block, variableStore); err != nil {
			return result, err
		}
	}
	return true, nil
}

func (r *Repeat) GetId() string {
	return r.Id
}

func (r *Repeat) GetBody() []Command {
	return r.Block
}

func (r *Repeat) GetArguments() *models.TreeNode {
	return nil
}

func (r *Repeat) GetCount() any {
	return r.Count
}

func (r *Repeat) GetCountType() string {
	return r.CountType
}
