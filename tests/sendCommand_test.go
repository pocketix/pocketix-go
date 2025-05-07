package tests

import (
	"testing"
	"time"

	"github.com/pocketix/pocketix-go/src/models"
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
