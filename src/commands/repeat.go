package commands

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/tree"
)

type Repeat struct {
	Id    string
	Count int
	Block []Command
}

func (r *Repeat) Execute(variableStore *models.VariableStore) (bool, error) {
	services.Logger.Println("Executing repeat")

	if r.Count < 0 {
		return false, fmt.Errorf("count cannot be negative")
	}

	for range r.Count {
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

func (r *Repeat) GetArguments() *tree.TreeNode {
	return nil
}

func (r *Repeat) GetCount() int {
	return r.Count
}
