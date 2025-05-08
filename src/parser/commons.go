package parser

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/pocketix/pocketix-go/src/types"
)

func HandleIfStatement(
	statement statements.Statement,
	previousStatement *statements.Statement,
	collect func(statement statements.Statement),
) error {
	if statement.GetId() == "if" {
		*previousStatement = statement
	} else if statement.GetId() == "else" {
		if *previousStatement != nil {
			(*previousStatement).(*statements.If).AddElseBlock(statement)
			collect(*previousStatement)
			*previousStatement = nil
		} else {
			services.Logger.Println("Error: Else without if")
			return fmt.Errorf("else without if")
		}
	} else if statement.GetId() == "elseif" {
		if *previousStatement != nil {
			(*previousStatement).(*statements.If).AddElseIfBlock(statement)
		} else {
			services.Logger.Println("Error: Elseif without if")
			return fmt.Errorf("elseif without if")
		}
	} else {
		if *previousStatement != nil {
			collect(*previousStatement)
			*previousStatement = nil
		}

		collect(statement)
	}
	return nil
}

// HandleDeviceTypeStatement checks if a statement is a deviceType and replaces it with the corresponding device command.
// It returns the processed statement (either the original or a new device command) and a boolean indicating if the statement was a deviceType.
//
// Parameters:
//   - statement: the statement to check
//   - devices: the list of devices to use for replacement
//   - deviceIndex: pointer to the current device index (will be incremented if a replacement is made)
//
// Returns:
//   - the processed statement (either the original or a new device command)
//   - a boolean indicating if the statement was a deviceType
func HandleDeviceTypeStatement(
	statement statements.Statement,
	devices []types.Device,
	deviceIndex *int,
) (statements.Statement, bool) {
	// Check if the statement is a deviceType
	if _, ok := statement.(*statements.DeviceType); ok {
		// If we have devices and the index is valid, replace the deviceType with the corresponding device
		if len(devices) > 0 && *deviceIndex < len(devices) {
			// Get the device at the current index
			device := devices[*deviceIndex]
			// Create a device command statement from the device
			deviceCommand := &statements.DeviceCommand{
				Id: device.ID,
				Arguments: &models.TreeNode{
					Type:  "str_opt", // Assuming the device values are string options
					Value: device.Values[0],
				},
			}
			// Increment the device index for the next deviceType
			*deviceIndex++
			// Return the device command
			return deviceCommand, true
		}
	}
	// Return the original statement if it's not a deviceType or no replacement was made
	return statement, false
}
