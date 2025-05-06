package statements

import (
	"fmt"
	"strings"

	"github.com/pocketix/pocketix-go/src/models"
)

type DeviceCommand struct {
	Id        string
	Arguments *models.TreeNode
}

func (d *DeviceCommand) Execute(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore) (bool, error) {
	return true, nil
}

func (d *DeviceCommand) GetId() string {
	return d.Id
}

func (d *DeviceCommand) GetArguments() *models.TreeNode {
	return d.Arguments
}

func (d *DeviceCommand) Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}

func (d *DeviceCommand) DeviceCommand2ModelsDeviceCommand() (models.DeviceCommand, error) {
	splittedDeviceId := strings.Split(d.Id, ".")
	if len(splittedDeviceId) != 2 {
		return models.DeviceCommand{}, fmt.Errorf("invalid device command id %s", d.Id)
	}

	return models.DeviceCommand{
		DeviceUID:         splittedDeviceId[0],
		CommandDenotation: splittedDeviceId[1],
		Arguments: models.TypeValue{
			Type:  d.Arguments.Type,
			Value: d.Arguments.Value,
		},
	}, nil
}
