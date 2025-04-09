package models

import "fmt"

type ReferencedValueStore struct {
	ReferencedValues map[string]ReferencedValue
}

func NewReferencedValueStore() *ReferencedValueStore {
	return &ReferencedValueStore{
		ReferencedValues: make(map[string]ReferencedValue),
	}
}

func (rvStore *ReferencedValueStore) AddReferencedValue(referenceTarget string, referencedValue *ReferencedValue) {
	rvStore.ReferencedValues[referenceTarget] = *referencedValue
}

func (rvStore *ReferencedValueStore) GetReferencedValues() map[string]string {
	referencedValues := make(map[string]string)

	for _, value := range rvStore.ReferencedValues {
		referencedValues[value.DeviceID] = value.ParameterName
	}
	return referencedValues
}

// SetValuesToReferenced sets the values of the referenced values to the referenced value store.
// Values are then used in the command execution.
func (rvStore *ReferencedValueStore) SetValuesToReferenced() {

}

// GetAndUpdateReferencedValue retrieves the value from the device and updates the referenced value in the store.
func (rvStore *ReferencedValueStore) GetAndUpdateReferencedValue(referencedTarget string) (*ReferencedValue, error) {
	referencedValue, ok := rvStore.ReferencedValues[referencedTarget]
	if !ok {
		return nil, fmt.Errorf("referenced value %s not found", referencedTarget)
	}
	return &referencedValue, nil
}
