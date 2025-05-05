package statements

import (
	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

// DeviceType represents a device type statement
type DeviceType struct {
	Id       string
	Type     string
	TypeType string // The type of the Type field (e.g., "string", "variable")
}

// Execute implements the Statement interface
func (d *DeviceType) Execute(variableStore *models.VariableStore, commandHandlingStore *models.CommandsHandlingStore) (bool, error) {
	services.Logger.Println("Executing deviceType with type:", d.Type)

	return true, nil
}

// GetId implements the Statement interface
func (d *DeviceType) GetId() string {
	return d.Id
}

// GetType returns the device type
func (d *DeviceType) GetType() string {
	return d.Type
}

// Validate implements the Statement interface
func (d *DeviceType) Validate(variableStore *models.VariableStore, referenceValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
