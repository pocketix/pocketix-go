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
	result, err := w.Arguments.Evaluate(variableStore, referencedValueStore)
	if err != nil {
		return false, err
	}

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

		snapshot, err := createSDParameterSnapshot(w.Arguments.Type, result, sdParameterInfo)
		if err != nil {
			return false, err
		}
		referencedValueStore.SetReferencedValue(referencedValue, snapshot, true)
	} else {
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

func createSDParameterSnapshot(snapshotType string, value any, sdParameterInfo models.SDInformationFromBackend) (models.SDParameterSnapshot, error) {
	sdParameterSnapshot := models.SDParameterSnapshot{
		DeviceID:    sdParameterInfo.Snapshot.DeviceID,
		SDParameter: sdParameterInfo.Snapshot.SDParameter,
	}

	switch snapshotType {
	case "string":
		strValue, ok := value.(string)
		if !ok {
			return models.SDParameterSnapshot{}, fmt.Errorf("expected string value, got %T", value)
		}
		sdParameterSnapshot.String = models.SnapshotString{
			Set:   true,
			Value: strValue,
		}
		return sdParameterSnapshot, nil
	case "number":
		numValue, ok := value.(float64)
		if !ok {
			return models.SDParameterSnapshot{}, fmt.Errorf("expected number value, got %T", value)
		}
		sdParameterSnapshot.Number = models.SnapshotNumber{
			Set:   true,
			Value: numValue,
		}
		return sdParameterSnapshot, nil
	case "boolean":
		boolValue, ok := value.(bool)
		if !ok {
			return models.SDParameterSnapshot{}, fmt.Errorf("expected boolean value, got %T", value)
		}
		sdParameterSnapshot.Boolean = models.SnapshotBoolean{
			Set:   true,
			Value: boolValue,
		}
		return sdParameterSnapshot, nil
	default:
		return models.SDParameterSnapshot{}, fmt.Errorf("unsupported snapshot type: %s", snapshotType)
	}
}
