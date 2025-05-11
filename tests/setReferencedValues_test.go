package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/statements"
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

func MockResolveCommandFunction(deviceUID string, commandDenotation string, infoType string) (models.SDInformationFromBackend, error) {
	return models.SDInformationFromBackend{
		DeviceUID: deviceUID,
		Command: models.SDCommand{
			CommandDenotation: commandDenotation,
			Payload:           "",
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

	statementList := make([]statements.Statement, 0)
	collector := &statements.ASTCollector{Target: &statementList}

	err := parser.Parse([]byte(program), variableStore, procedureStore, referencedValueStore, collector)
	assert.Nil(err, "Error should be nil, but got: %v", err)
	var wasCommandSet bool

	if deviceCommand, ok := statementList[0].(*statements.DeviceCommand); ok {
		dvcCommand, ok := deviceCommand.DeviceCommand2ModelsDeviceCommand()
		assert.True(ok, "Expected DeviceCommand, but got: %T", deviceCommand)

		assert.Equal("TemperatureDevice-1", dvcCommand.DeviceUID, "Device ID should be 'TemperatureDevice-1', but got: %s", dvcCommand.DeviceUID)
		assert.Equal("setTemperature", dvcCommand.CommandDenotation, "Command should be 'setTemperature', but got: %s", dvcCommand.CommandDenotation)

		var sdCommandInformation models.SDInformationFromBackend
		if cmd, ok := collector.IsDeviceCommandCollected(dvcCommand); ok {
			sdCommandInformation = cmd
		} else {
			sdCommandInfo, err2 := referencedValueStore.ResolveDeviceInformationFunction(dvcCommand.DeviceUID, dvcCommand.CommandDenotation, "sdCommand")
			collector.CollectDeviceCommand(sdCommandInfo)
			assert.Nil(err2, "Error should be nil, but got: %v", err2)
			sdCommandInformation = sdCommandInfo
			wasCommandSet = false
		}
		assert.False(wasCommandSet, "Expected command to be collected, but it was not")
		assert.Equal("TemperatureDevice-1", sdCommandInformation.DeviceUID, "Device ID should be 'TemperatureDevice-1', but got: %s", sdCommandInformation.DeviceUID)
		assert.Equal("setTemperature", sdCommandInformation.Command.CommandDenotation, "Command should be 'setTemperature', but got: %s", sdCommandInformation.Command.CommandDenotation)

	} else {
		t.Errorf("Expected DeviceCommand, got %T", statementList[0])
	}

	if deviceCommand, ok := statementList[1].(*statements.DeviceCommand); ok {
		dvcCommand, ok := deviceCommand.DeviceCommand2ModelsDeviceCommand()
		assert.True(ok, "Expected DeviceCommand, but got: %T", deviceCommand)

		assert.Equal("TemperatureDevice-1", dvcCommand.DeviceUID, "Device ID should be 'TemperatureDevice-1', but got: %s", dvcCommand.DeviceUID)
		assert.Equal("setTemperature", dvcCommand.CommandDenotation, "Command should be 'setTemperature', but got: %s", dvcCommand.CommandDenotation)

		var sdCommandInformation models.SDInformationFromBackend
		if cmd, ok := collector.IsDeviceCommandCollected(dvcCommand); ok {
			sdCommandInformation = cmd
			wasCommandSet = true
		} else {
			sdCommandInfo, err2 := referencedValueStore.ResolveDeviceInformationFunction(dvcCommand.DeviceUID, dvcCommand.CommandDenotation, "sdCommand")
			assert.Nil(err2, "Error should be nil, but got: %v", err2)
			sdCommandInformation = sdCommandInfo
		}
		assert.True(wasCommandSet, "Expected command to be collected, but it was not")
		assert.Equal("TemperatureDevice-1", sdCommandInformation.DeviceUID, "Device ID should be 'TemperatureDevice-1', but got: %s", sdCommandInformation.DeviceUID)
		assert.Equal("setTemperature", sdCommandInformation.Command.CommandDenotation, "Command should be 'setTemperature', but got: %s", sdCommandInformation.Command.CommandDenotation)
	} else {
		t.Errorf("Expected DeviceCommand, got %T", statementList[1])
	}
}
