package commands

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type Case struct {
	Id    string
	Block []Command
	Type  string
	Value any
}

func (c *Case) Execute(variableStore *models.VariableStore, commandHandlingStore *models.CommandsHandlingStore) (bool, error) {
	services.Logger.Println("Executing case", c.Value)
	return ExecuteCommands(c.Block, variableStore, commandHandlingStore)
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

func (c *Case) Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error {
	if len(args) == 0 {
		if c.Type == "boolean_expression" {
			return fmt.Errorf("case value should be constant")
		} else {
			return nil
		}
	}
	selectorType := args[0].(string)

	if c.Type == "variable" {
		variable, err := variableStore.GetVariable(c.Value.(string))
		if err != nil {
			return err
		}
		if variable.Type != selectorType {
			return fmt.Errorf("case value type %s does not match selector type %s", variable.Type, selectorType)
		}
	}
	if c.Type != selectorType {
		return fmt.Errorf("case value type %s does not match selector type %s", c.Type, selectorType)
	}
	return nil
}
