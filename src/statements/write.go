package statements

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type Write struct {
	Id        string
	Arguments *models.TreeNode
}

func (w *Write) Execute(
	variableStore *models.VariableStore,
	referencedValueStore *models.ReferencedValueStore,
	deviceCommands []models.SDInformationFromBackend,
	callback func(deviceCommand models.SDCommandInvocation),
) (bool, error) {
	services.Logger.Println("Executing write statement")

	uid, parameter, ok := models.FromReferencedTarget(w.Arguments.Reference)
	if !ok {
		return false, fmt.Errorf("invalid reference target")
	}

	sdParameterInfo, err := referencedValueStore.ResolveDeviceInformationFunction(uid, parameter, "sdParameter", &deviceCommands)
	if err != nil {
		return false, err
	}
	services.Logger.Printf("Resolved parameter info: %v", sdParameterInfo)

	referencedValue, ok := models.NewReferencedValue(w.Arguments.Reference)
	if ok {
		referencedValueStore.AddReferencedValue(w.Arguments.Reference, referencedValue)
		referencedValueStore.SetReferencedValue(w.Arguments.Reference, sdParameterInfo.Snapshot)
	}

	if !ok {
		services.Logger.Printf("Failed to create referenced value for %s", w.Arguments.Reference)
	}

	return false, nil
}

func (w *Write) GetId() string {
	return w.Id
}

func (w *Write) Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error {
	// Implementation for validating the write statement
	return nil
}

func (w *Write) GetArguments() *models.TreeNode {
	return w.Arguments
}
