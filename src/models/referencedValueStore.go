package models

import "fmt"

type ReferencedValueStore struct {
	ReferencedValues                 map[string]ReferencedValue
	ResolveDeviceInformationFunction func(deviceUID string, paramDenotation string, infoType string) (SDInformationFromBackend, error)
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
	referencedValue, ok := rvStore.ReferencedValues[referenceTarget]
	if !ok {
		return nil, fmt.Errorf("referenced value %s not found", referenceTarget)
	}
	if snapshot.String != nil {
		referencedValue.Value = *snapshot.String
		referencedValue.Type = "string"
	} else if snapshot.Number != nil {
		referencedValue.Value = *snapshot.Number
		referencedValue.Type = "number"
	} else if snapshot.Boolean != nil {
		referencedValue.Value = *snapshot.Boolean
		referencedValue.Type = "boolean"
	} else {
		return nil, fmt.Errorf("no valid value found in the snapshot")
	}
	referencedValue.IsSet = true
	rvStore.ReferencedValues[referenceTarget] = referencedValue
	return referencedValue.Value, nil
}

func (rvStore *ReferencedValueStore) SetResolveParameterFunction(fn func(deviceUID string, paramDenotation string, infoType string) (SDInformationFromBackend, error)) {
	rvStore.ResolveDeviceInformationFunction = fn
}
