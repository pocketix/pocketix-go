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

func (rvStore *ReferencedValueStore) GetReferencedValues() map[string]ReferencedValue {
	return rvStore.ReferencedValues
}

func (rvStore *ReferencedValueStore) SetValuesToReferenced(referencedValuesToUpdate []ReferenceValueResponseFromBackend) error {
	for _, referencedValueToUpdate := range referencedValuesToUpdate {
		referencedTarget := referencedValueToUpdate.DeviceID + "." + referencedValueToUpdate.SDType.SDParameters[0].ParameterID
		referencedValue, ok := rvStore.ReferencedValues[referencedTarget]
		if !ok {
			return fmt.Errorf("referenced value %s not found", referencedTarget)
		}

		latestSnapshot := getLatestSnapshot(referencedValueToUpdate.SDParameterSnapshots)
		if latestSnapshot.String != nil {
			referencedValue.Value = *latestSnapshot.String
			referencedValue.Type = "string"
		} else if latestSnapshot.Number != nil {
			referencedValue.Value = *latestSnapshot.Number
			referencedValue.Type = "number"
		} else if latestSnapshot.Boolean != nil {
			referencedValue.Value = *latestSnapshot.Boolean
			referencedValue.Type = "boolean"
		} else {
			return fmt.Errorf("no valid value found in the latest snapshot")
		}
		referencedValue.DeviceID = referencedValueToUpdate.DeviceID
		referencedValue.ParameterName = referencedValueToUpdate.SDType.SDParameters[0].ParameterID
		rvStore.ReferencedValues[referencedTarget] = referencedValue
	}
	return nil
}

func getLatestSnapshot(snapshots []SDParameterSnapshot) SDParameterSnapshot {
	return snapshots[0]
}

// GetAndUpdateReferencedValue retrieves the value from the device and updates the referenced value in the store.
func (rvStore *ReferencedValueStore) GetAndUpdateReferencedValue(referencedTarget string) (*ReferencedValue, error) {
	referencedValue, ok := rvStore.ReferencedValues[referencedTarget]
	if !ok {
		return nil, fmt.Errorf("referenced value %s not found", referencedTarget)
	}
	return &referencedValue, nil
}
