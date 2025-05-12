package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/stretchr/testify/assert"
)

// func(deviceUID string, paramDenotation string) (any, error)
func MockResolveParameterFunction(deviceUID string, paramDenotation string, infoType string, deviceCommands *[]models.SDInformationFromBackend) (models.SDInformationFromBackend, error) {
	return models.SDInformationFromBackend{
		DeviceUID: deviceUID,
		Snapshot: models.SDParameterSnapshot{
			SDParameter: paramDenotation,
			Number:      func(v float64) *float64 { return &v }(10.0),
		},
	}, nil
}

func MockResolveCommandFunction(deviceUID string, commandDenotation string, infoType string, deviceCommands *[]models.SDInformationFromBackend) (models.SDInformationFromBackend, error) {
	info := models.SDInformationFromBackend{
		DeviceUID: deviceUID,
		Command: models.SDCommand{
			CommandDenotation: commandDenotation,
			Payload:           "",
		},
	}

	*deviceCommands = append(*deviceCommands, info)

	return info, nil
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

	sdInformation, err := referencedValueStore.ResolveDeviceInformationFunction("Device-1", "test", "sdParameter", nil)
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

func TestRepeatedCommandInvocation(t *testing.T) {
	assert := assert.New(t)

	program := `
		{
			"header": {
				"userVariables": {},
				"userProcedures": {}
			},
			"block": [
				{
					"id": "TemperatureDevice-1.setTemperature",
					"arguments": [
						{
							"type": "str_opt",
							"value": "ideal"
						}
					]
				},
				{
					"id": "TemperatureDevice-1.setTemperature",
					"arguments": [
						{
							"type": "str_opt",
							"value": "ideal"
						}
					]
				}
			]
		}
	`

	variableStore := models.NewVariableStore()
	procedureStore := models.NewProcedureStore()
	referencedValueStore := models.NewReferencedValueStore()
	referencedValueStore.SetResolveParameterFunction(MockResolveCommandFunction)

	var interpretInvocationsToSend []models.SDCommandInvocation
	callback := func(deviceCommand models.SDCommandInvocation) {
		interpretInvocationsToSend = append(interpretInvocationsToSend, deviceCommand)
	}
	statementList := make([]statements.Statement, 0)
	collector := &statements.ASTCollector{Target: &statementList, DeviceCommands: make([]models.SDInformationFromBackend, 0)}

	err := parser.Parse([]byte(program), variableStore, procedureStore, referencedValueStore, collector)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	_, err = statementList[0].Execute(variableStore, referencedValueStore, collector.DeviceCommands, callback)
	assert.Nil(err, "Error should be nil, but got: %v", err)

	deviceCommand, ok := statementList[0].(*statements.DeviceCommand)
	assert.True(ok, "Expected DeviceCommand type, but got: %T", statementList[0])

	deviceCommand2, ok := deviceCommand.DeviceCommand2ModelsDeviceCommand()
	assert.True(ok, "Expected DeviceCommand2ModelsDeviceCommand type, but got: %T", deviceCommand)

	_, isIn := statements.Contains(collector.DeviceCommands, deviceCommand2)
	assert.False(isIn)

	_, err = statementList[1].Execute(variableStore, referencedValueStore, collector.DeviceCommands, callback)
	assert.Nil(err, "Error should be nil, but got: %v", err)
}
