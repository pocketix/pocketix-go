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

func (rvStore *ReferencedValueStore) SetReferencedValue(referencedValue *ReferencedValue, snapshot SDParameterSnapshot, isSet bool) (any, error) {
	referenceTarget := referencedValue.ToReferenceTarget()
	services.Logger.Printf("Setting referenced value for %s: %v", referenceTarget, snapshot)

	valueFromStore, ok := rvStore.ReferencedValues[referenceTarget]
	if !ok {
		return nil, fmt.Errorf("referenced value %s not found", referenceTarget)
	}

	if snapshot.String.Set {
		valueFromStore.Value = snapshot.String.Value
		valueFromStore.Type = "string"
	} else if snapshot.Number.Set {
		valueFromStore.Value = snapshot.Number.Value
		valueFromStore.Type = "number"
	} else if snapshot.Boolean.Set {
		valueFromStore.Value = snapshot.Boolean.Value
		valueFromStore.Type = "boolean"
	}

	valueFromStore.IsSet = isSet
	valueFromStore.DeviceID = snapshot.DeviceID
	valueFromStore.ParameterID = snapshot.SDParameter

	rvStore.ReferencedValues[referenceTarget] = valueFromStore
	services.Logger.Printf("Set referenced value for %s: %v", referenceTarget, valueFromStore)
	return valueFromStore.Value, nil
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
