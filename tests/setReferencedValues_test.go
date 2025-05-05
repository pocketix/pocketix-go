package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/stretchr/testify/assert"
)

// func(deviceUID string, paramDenotation string) (any, error)
func MockResolveParameterFunction(deviceUID string, paramDenotation string) (any, string, error) {
	return 10.0, "number", nil
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

	parameterValue, valueType, err := referencedValueStore.ResolveParameterFunction("Device-1", "test")
	assert.Nil(err)
	assert.Equal(10.0, parameterValue)
	assert.Equal("number", valueType)
	referencedValueStore.SetReferencedValue("Device-1.test", parameterValue, valueType)

	assert.Equal(10.0, referencedValueStore.GetReferencedValues()["Device-1.test"].Value)
	assert.Equal("number", referencedValueStore.GetReferencedValues()["Device-1.test"].Type)
	assert.Equal("Device-1", referencedValueStore.GetReferencedValues()["Device-1.test"].DeviceID)
	assert.Equal("test", referencedValueStore.GetReferencedValues()["Device-1.test"].ParameterName)
}
