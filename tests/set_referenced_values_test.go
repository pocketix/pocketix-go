package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/stretchr/testify/assert"
)

func TestSetReferencedValues(t *testing.T) {
	assert := assert.New(t)

	referencedValueStore := models.NewReferencedValueStore()
	referencedValueStore.AddReferencedValue("Device-1.test", &models.ReferencedValue{
		DeviceID:      "Device-1",
		ParameterName: "test",
	})
	assert.Equal(1, len(referencedValueStore.GetReferencedValues()))
	assert.Equal("Device-1", referencedValueStore.GetReferencedValues()["Device-1.test"].DeviceID)
	assert.Equal("test", referencedValueStore.GetReferencedValues()["Device-1.test"].ParameterName)
	assert.Equal("", referencedValueStore.GetReferencedValues()["Device-1.test"].Type)
	assert.Equal(nil, referencedValueStore.GetReferencedValues()["Device-1.test"].Value)

	referencedValueStore.SetValuesToReferenced([]models.ReferenceValueResponseFromBackend{
		{
			DeviceID: "Device-1",
			SDType: models.SDType{
				SDParameters: []models.SDParameter{
					{
						ParameterID:         "test",
						ParameterDenotation: "test",
					},
				},
				SDCommands: []models.SDCommand{},
			},
			SDParameterSnapshots: []models.SDParameterSnapshot{
				{
					SDParameter: "test",
					String:      nil,
					Number:      func(v float64) *float64 { return &v }(10),
					Boolean:     nil,
				},
			},
		},
	})

	assert.Equal(1, len(referencedValueStore.GetReferencedValues()))
	assert.Equal("Device-1", referencedValueStore.GetReferencedValues()["Device-1.test"].DeviceID)
	assert.Equal("test", referencedValueStore.GetReferencedValues()["Device-1.test"].ParameterName)
	assert.Equal("number", referencedValueStore.GetReferencedValues()["Device-1.test"].Type)
	assert.Equal(10.0, referencedValueStore.GetReferencedValues()["Device-1.test"].Value)
}
