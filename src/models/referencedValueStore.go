package models

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/services"
)

type ReferencedValueStore struct {
	ReferencedValues                 map[string]ReferencedValue
	ResolveDeviceInformationFunction func(deviceUID string, paramDenotation string, infoType string, deviceCommands *[]SDInformationFromBackend) (SDInformationFromBackend, error)
}

func NewReferencedValueStore() *ReferencedValueStore {
	return &ReferencedValueStore{
		ReferencedValues: make(map[string]ReferencedValue),
	}
}

func (rvStore *ReferencedValueStore) AddReferencedValue(referenceTarget string, referencedValue *ReferencedValue) {
	rvStore.ReferencedValues[referenceTarget] = *referencedValue
}

func (rvStore *ReferencedValueStore) GetReferencedValues() map[string]ReferencedValue {
	return rvStore.ReferencedValues
}

func (rvStore *ReferencedValueStore) GetReferencedValueFromStore(referencedTarget string) (*ReferencedValue, bool) {
	referencedValue, ok := rvStore.ReferencedValues[referencedTarget]
	if !ok {
		return nil, false
	}
	return &referencedValue, true
}

func (rvStore *ReferencedValueStore) SetReferencedValue(referenceTarget string, snapshot SDParameterSnapshot) (any, error) {
	services.Logger.Printf("Setting referenced value for %s: %v", referenceTarget, snapshot)
	referencedValue, ok := rvStore.ReferencedValues[referenceTarget]
	if !ok {
		return nil, fmt.Errorf("referenced value %s not found", referenceTarget)
	}
	if snapshot.String.Set {
		referencedValue.Value = snapshot.String.Value
		referencedValue.Type = "string"
	} else if snapshot.Number.Set {
		referencedValue.Value = snapshot.Number.Value
		referencedValue.Type = "number"
	} else if snapshot.Boolean.Set {
		referencedValue.Value = snapshot.Boolean.Value
		referencedValue.Type = "boolean"
	} else {
		return nil, fmt.Errorf("no valid value found in the snapshot")
	}

	referencedValue.IsSet = true
	referencedValue.DeviceID = snapshot.DeviceID
	referencedValue.ParameterID = snapshot.SDParameter

	rvStore.ReferencedValues[referenceTarget] = referencedValue
	services.Logger.Printf("Set referenced value for %s: %v", referenceTarget, referencedValue)
	return referencedValue.Value, nil
}

func (rvStore *ReferencedValueStore) SetResolveParameterFunction(fn func(deviceUID string, paramDenotation string, infoType string, deviceCommands *[]SDInformationFromBackend) (SDInformationFromBackend, error)) {
	rvStore.ResolveDeviceInformationFunction = fn
}

func (rvStore *ReferencedValueStore) GetSetReferencedValues() map[string]ReferencedValue {
	setValues := make(map[string]ReferencedValue)
	for key, value := range rvStore.ReferencedValues {
		if value.IsSet {
			setValues[key] = value
		}
	}
	return setValues
}
