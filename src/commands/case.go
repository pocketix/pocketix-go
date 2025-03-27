package commands

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type Case struct {
	Id    string
	Block []Command
	Type  string
	Value any
}

func (c *Case) Execute(variableStore *models.VariableStore) (bool, error) {
	services.Logger.Println("Executing case", c.Value)
	return ExecuteCommands(c.Block, variableStore)
}

func (c *Case) GetId() string {
	return c.Id
}

func (c *Case) GetBody() []Command {
	return c.Block
}

func (c *Case) GetArguments() *models.TreeNode {
	return nil
}

func (c *Case) GetValue() any {
	return c.Value
}
