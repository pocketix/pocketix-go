package commands

import "github.com/pocketix/pocketix-go/src/models"

type TypeValue struct {
	Type  string
	Value any
}

type DeviceCommand struct {
	DeviceID  string
	Command   string
	Arguments []TypeValue
}

// Execute is a command that sends a command to a device.
func (d *DeviceCommand) Execute(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore) (bool, error) {
	return true, nil
}

func (d *DeviceCommand) GetId() string {
	return d.DeviceID
}

func (d *DeviceCommand) GetBody() []Command {
	return nil
}

func (d *DeviceCommand) GetArguments() *models.TreeNode {
	return nil
}

func (d *DeviceCommand) Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
