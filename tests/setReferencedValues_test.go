package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/stretchr/testify/assert"
)

// func(deviceUID string, paramDenotation string) (any, error)
func MockResolveParameterFunction(deviceUID string, paramDenotation string, infoType string) (models.SDInformationFromBackend, error) {
	return models.SDInformationFromBackend{
		DeviceUID: deviceUID,
		Snapshot: models.SDParameterSnapshot{
			SDParameter: paramDenotation,
			Number:      func(v float64) *float64 { return &v }(10.0),
		},
	}, nil
}

func TestSetReferencedValues(t *testing.T) {
	assert := assert.New(t)

	referencedValueStore := models.NewReferencedValueStore()
	referencedValueStore.SetResolveParameterFunction(MockResolveParameterFunction)
	referencedValueStore.AddReferencedValue("Device-1.test", &models.ReferencedValue{
		DeviceID:      "Device-1",
		ParameterName: "test",
	})
	assert.Equal(1, len(referencedValueStore.GetReferencedValues()))
	assert.Equal("Device-1", referencedValueStore.GetReferencedValues()["Device-1.test"].DeviceID)
	assert.Equal("test", referencedValueStore.GetReferencedValues()["Device-1.test"].ParameterName)
	assert.Equal("", referencedValueStore.GetReferencedValues()["Device-1.test"].Type)
	assert.Equal(nil, referencedValueStore.GetReferencedValues()["Device-1.test"].Value)

	sdInformation, err := referencedValueStore.ResolveDeviceInformationFunction("Device-1", "test", "sdParameter")
	assert.Nil(err)
	assert.Equal("Device-1", sdInformation.DeviceUID)
	assert.Equal("test", sdInformation.Snapshot.SDParameter)
	assert.Equal(10.0, *sdInformation.Snapshot.Number)

	value, err := referencedValueStore.SetReferencedValue("Device-1.test", sdInformation.Snapshot)
	assert.Nil(err)
	assert.Equal(10.0, value)

	assert.Equal(10.0, referencedValueStore.GetReferencedValues()["Device-1.test"].Value)
	assert.Equal("number", referencedValueStore.GetReferencedValues()["Device-1.test"].Type)
	assert.Equal("Device-1", referencedValueStore.GetReferencedValues()["Device-1.test"].DeviceID)
	assert.Equal("test", referencedValueStore.GetReferencedValues()["Device-1.test"].ParameterName)
}
