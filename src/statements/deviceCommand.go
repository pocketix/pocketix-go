package statements

import (
	"fmt"
	"strings"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type DeviceCommand struct {
	Id        string
	Arguments *models.TreeNode
}

func (d *DeviceCommand) Execute(
	variableStore *models.VariableStore,
	referencedValueStore *models.ReferencedValueStore,
	deviceCommands []models.SDInformationFromBackend,
	callback func(deviceCommand models.SDCommandInvocation),
) (bool, error) {
	// informationFromBackend, err := referencedValueStore.ResolveDeviceInformationFunction(deviceUID, commandDenotation, "sdCommand", deviceCommands)
	services.Logger.Printf("Executing device command: %v", d.Id)
	deviceCommand, ok := d.DeviceCommand2ModelsDeviceCommand()
	if !ok {
		return false, fmt.Errorf("failed to convert DeviceCommand to models.DeviceCommand, probably due to missing procedure %s", d.Id)
	}

	var sdCommandInformation models.SDInformationFromBackend
	if backendInformation, ok := Contains(deviceCommands, deviceCommand); ok {
		sdCommandInformation = backendInformation
	} else {
		sdCommandInfo, err := referencedValueStore.ResolveDeviceInformationFunction(deviceCommand.DeviceUID, deviceCommand.CommandDenotation, "sdCommand", &deviceCommands)
		services.Logger.Printf("Resolved device information for command %s: %v", d.Id, sdCommandInfo)
		if err != nil {
			services.Logger.Printf("Failed to resolve device information for command %s: %v", d.Id, err)
			return false, err
		}
		sdCommandInformation = sdCommandInfo
	}

	sdCommandInvocation, err := deviceCommand.PrepareCommandToSend(sdCommandInformation)
	if err != nil {
		return false, err
	}
	callback(*sdCommandInvocation)
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

	if d.Arguments == nil {
		return models.DeviceCommand{
			DeviceUID:         prefix,
			CommandDenotation: last,
		}, true
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

func Contains(slice []models.SDInformationFromBackend, item models.DeviceCommand) (models.SDInformationFromBackend, bool) {
	for _, v := range slice {
		if v.DeviceUID == item.DeviceUID && v.Command.CommandDenotation == item.CommandDenotation {
			return v, true
		}
	}
	return models.SDInformationFromBackend{}, false
}
