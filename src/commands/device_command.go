package commands

import (
	"fmt"
	"strings"

	"github.com/pocketix/pocketix-go/src/models"
)

type DeviceCommand struct {
	Id        string
	Arguments *models.TreeNode
}

func (d *DeviceCommand) Execute(variableStore *models.VariableStore, commandHandlingStore *models.CommandsHandlingStore) (bool, error) {
	return true, nil
}

func (d *DeviceCommand) GetId() string {
	return d.Id
}

func (d *DeviceCommand) GetBody() []Command {
	return nil
}

func (d *DeviceCommand) GetArguments() *models.TreeNode {
	return d.Arguments
}

func (d *DeviceCommand) Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}

func (d *DeviceCommand) DeviceCommand2ModelsDeviceCommand(commandHandlingStore *models.CommandsHandlingStore) (models.DeviceCommand, error) {
	splittedDeviceId := strings.Split(d.Id, ".")
	if len(splittedDeviceId) != 2 {
		return models.DeviceCommand{}, fmt.Errorf("invalid device command id %s", d.Id)
	}

	argumentsList := make([]models.TypeValue, 0)
	for _, arg := range d.Arguments.Children {
		argumentsList = append(argumentsList, models.TypeValue{
			Type:  arg.Type,
			Value: arg.Value,
		})
	}

	return models.DeviceCommand{
		DeviceID:  splittedDeviceId[0],
		Command:   splittedDeviceId[1],
		Arguments: argumentsList,
	}, nil
}
