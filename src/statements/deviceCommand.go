package statements

import (
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

func (d *DeviceCommand) DeviceCommand2ModelsDeviceCommand() (models.DeviceCommand, bool) {
	lastDot := strings.LastIndex(d.Id, ".")
	if lastDot == -1 {
		return models.DeviceCommand{}, false
	}

	prefix := d.Id[:lastDot]
	last := d.Id[lastDot+1:]

	if prefix == "" || last == "" {
		return models.DeviceCommand{}, false
	}

	return models.DeviceCommand{
		DeviceUID:         prefix,
		CommandDenotation: last,
		Arguments: models.TypeValue{
			Type:  d.Arguments.Type,
			Value: d.Arguments.Value,
		},
	}, true
}
