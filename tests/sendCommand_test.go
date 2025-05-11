package tests

import (
	"testing"
	"time"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/stretchr/testify/assert"
)

func TestSimpleCommandWithoutArguments(t *testing.T) {
	assert := assert.New(t)

	deviceCommand := &models.DeviceCommand{
		DeviceUID:         "Device-1",
		CommandDenotation: "Command",
		Arguments:         models.TypeValue{},
	}

	response := &models.SDInformationFromBackend{
		DeviceUID: "Device-1",
		Command: models.SDCommand{
			CommandDenotation: "Command",
			Payload:           "",
		},
	}

	dc, err := deviceCommand.PrepareCommandToSend(*response)
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.Equal("Device-1", dc.DeviceID, "Device ID should be 'Device-1', but got: %s", dc.DeviceID)
	assert.Equal("Command", dc.CommandDenotation, "Command should be 'Command', but got: %s", dc.CommandDenotation)
	assert.Equal("", dc.Payload, "Payload should be empty, but got: %s", dc.Payload)
	assert.Equal(time.Now().Format(time.RFC3339), dc.InvocationTime, "Invocation time should be current time, but got: %s", dc.InvocationTime)
}

func TestCommandWithArguments(t *testing.T) {
	assert := assert.New(t)

	deviceCommand := &models.DeviceCommand{
		DeviceUID:         "Device-1",
		CommandDenotation: "Command",
		Arguments: models.TypeValue{
			Type:  "string",
			Value: "Test",
		},
	}

	response := &models.SDInformationFromBackend{
		DeviceUID: "Device-1",
		Command: models.SDCommand{
			CommandDenotation: "Command",
			Payload:           `{"name":"testing", "possibleValues":["Test", "Test2"]}`,
		},
	}

	dc, err := deviceCommand.PrepareCommandToSend(*response)
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.Equal("Device-1", dc.DeviceID, "Device ID should be 'Device-1', but got: %s", dc.DeviceID)
	assert.Equal("Command", dc.CommandDenotation, "Command should be 'Command', but got: %s", dc.CommandDenotation)
	assert.Equal(`{"name":"testing","value":"Test"}`, dc.Payload, "Payload should be `{\"name\":\"testing\",\"value\":\"Test\"}`, but got: %s", dc.Payload)
	assert.Equal(time.Now().Format(time.RFC3339), dc.InvocationTime, "Invocation time should be current time, but got: %s", dc.InvocationTime)
}

func TestDeviceCommand2ModelsDeviceCommand(t *testing.T) {
	assert := assert.New(t)

	deviceCommand := &statements.DeviceCommand{
		Id:        "Device-1.Command",
		Arguments: &models.TreeNode{},
	}

	modelsDeviceCommand, ok := deviceCommand.DeviceCommand2ModelsDeviceCommand()
	assert.True(ok, "Conversion should be successful, but got: %v", ok)
	assert.Equal("Device-1", modelsDeviceCommand.DeviceUID, "Device UID should be 'Device-1', but got: %s", modelsDeviceCommand.DeviceUID)
	assert.Equal("Command", modelsDeviceCommand.CommandDenotation, "Command denotation should be 'Command', but got: %s", modelsDeviceCommand.CommandDenotation)

	deviceCommand = &statements.DeviceCommand{
		Id:        "Device-1",
		Arguments: &models.TreeNode{},
	}

	modelsDeviceCommand, ok = deviceCommand.DeviceCommand2ModelsDeviceCommand()
	assert.False(ok, "Conversion should fail, but got: %v", ok)

	deviceCommand = &statements.DeviceCommand{
		Id:        "Device-1.Command.",
		Arguments: &models.TreeNode{},
	}

	modelsDeviceCommand, ok = deviceCommand.DeviceCommand2ModelsDeviceCommand()
	assert.False(ok, "Conversion should fail, but got: %v", ok)

	deviceCommand = &statements.DeviceCommand{
		Id:        ".Command",
		Arguments: &models.TreeNode{},
	}

	modelsDeviceCommand, ok = deviceCommand.DeviceCommand2ModelsDeviceCommand()
	assert.False(ok, "Conversion should fail, but got: %v", ok)

	deviceCommand = &statements.DeviceCommand{
		Id:        "Device-1.Test.Command",
		Arguments: &models.TreeNode{},
	}

	modelsDeviceCommand, ok = deviceCommand.DeviceCommand2ModelsDeviceCommand()
	assert.True(ok, "Conversion should be successful, but got: %v", ok)
	assert.Equal("Device-1.Test", modelsDeviceCommand.DeviceUID, "Device UID should be 'Device-1.Test', but got: %s", modelsDeviceCommand.DeviceUID)
	assert.Equal("Command", modelsDeviceCommand.CommandDenotation, "Command denotation should be 'Command', but got: %s", modelsDeviceCommand.CommandDenotation)

	deviceCommand = &statements.DeviceCommand{
		Id:        "Device-1.Test.Test2.Command",
		Arguments: &models.TreeNode{},
	}

	modelsDeviceCommand, ok = deviceCommand.DeviceCommand2ModelsDeviceCommand()
	assert.True(ok, "Conversion should be successful, but got: %v", ok)
	assert.Equal("Device-1.Test.Test2", modelsDeviceCommand.DeviceUID, "Device UID should be 'Device-1.Test.Test2', but got: %s", modelsDeviceCommand.DeviceUID)
	assert.Equal("Command", modelsDeviceCommand.CommandDenotation, "Command denotation should be 'Command', but got: %s", modelsDeviceCommand.CommandDenotation)
}
