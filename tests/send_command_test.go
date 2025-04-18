package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/stretchr/testify/assert"
)

func TestSendSimpleCommandWithoutArguments(t *testing.T) {
	assert := assert.New(t)

	deviceCommand := &models.DeviceCommand{
		DeviceID:  "Device-1",
		Command:   "Command",
		Arguments: []models.TypeValue{},
	}

	dc, parameterSnapshotsToUpdate, err := deviceCommand.SendCommandToDevice()
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.Equal("Device-1", dc.DeviceID, "Device ID should be 'Device-1', but got: %s", dc.DeviceID)
	assert.Equal("Command", dc.Command, "Command should be 'Command', but got: %s", dc.Command)
	assert.Nil(parameterSnapshotsToUpdate, "Parameter snapshots should be nil, but got: %v", parameterSnapshotsToUpdate)
}

func TestSendCommandWithArguments(t *testing.T) {
	assert := assert.New(t)

	deviceCommand := &models.DeviceCommand{
		DeviceID: "Device-1",
		Command:  "Command",
		Arguments: []models.TypeValue{
			{
				Type:  "string",
				Value: "Test",
			},
		},
	}

	dc, parameterSnapshotsToUpdate, err := deviceCommand.SendCommandToDevice()
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.Equal("Device-1", dc.DeviceID, "Device ID should be 'Device-1', but got: %s", dc.DeviceID)
	assert.Equal("Command", dc.Command, "Command should be 'Command', but got: %s", dc.Command)
	assert.NotNil(parameterSnapshotsToUpdate, "Parameter snapshots should not be nil, but got: %v", parameterSnapshotsToUpdate)
	assert.Equal(1, len(*parameterSnapshotsToUpdate), "Parameter snapshots length should be 1, but got: %d", len(*parameterSnapshotsToUpdate))
	assert.Equal("Device-1", (*parameterSnapshotsToUpdate)[0].DeviceID, "Device ID in parameter snapshots should be 'Device-1', but got: %s", (*parameterSnapshotsToUpdate)[0].DeviceID)
	assert.Equal("Command", (*parameterSnapshotsToUpdate)[0].SDParameter, "Command in parameter snapshots should be 'Command', but got: %s", (*parameterSnapshotsToUpdate)[0].SDParameter)
	assert.Equal("Test", *(*parameterSnapshotsToUpdate)[0].String, "String value in parameter snapshots should be 'Test', but got: %s", *(*parameterSnapshotsToUpdate)[0].String)
	assert.Nil((*parameterSnapshotsToUpdate)[0].Number, "Number value in parameter snapshots should be nil, but got: %v", (*parameterSnapshotsToUpdate)[0].Number)
	assert.Nil((*parameterSnapshotsToUpdate)[0].Boolean, "Boolean value in parameter snapshots should be nil, but got: %v", (*parameterSnapshotsToUpdate)[0].Boolean)
}

func TestSendTwoCommands(t *testing.T) {
	assert := assert.New(t)

	deviceCommandList := []*models.DeviceCommand{
		{
			DeviceID:  "Device-1",
			Command:   "Command-1",
			Arguments: []models.TypeValue{},
		},
		{
			DeviceID: "Device-2",
			Command:  "Command-2",
			Arguments: []models.TypeValue{
				{
					Type:  "number",
					Value: 10.0,
				},
			},
		},
	}

	deviceCmd, parameterSnapshotsToUpdate, err := deviceCommandList[0].SendCommandToDevice()
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.Equal("Device-1", deviceCmd.DeviceID, "Device ID should be 'Device-1', but got: %s", deviceCmd.DeviceID)
	assert.Equal("Command-1", deviceCmd.Command, "Command should be 'Command-1', but got: %s", deviceCmd.Command)
	assert.Nil(parameterSnapshotsToUpdate, "Parameter snapshots should be nil, but got: %v", parameterSnapshotsToUpdate)

	deviceCmd, parameterSnapshotsToUpdate, err = deviceCommandList[1].SendCommandToDevice()
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.Equal("Device-2", deviceCmd.DeviceID, "Device ID should be 'Device-2', but got: %s", deviceCmd.DeviceID)
	assert.Equal("Command-2", deviceCmd.Command, "Command should be 'Command-2', but got: %s", deviceCmd.Command)
	assert.NotNil(parameterSnapshotsToUpdate, "Parameter snapshots should not be nil, but got: %v", parameterSnapshotsToUpdate)
	assert.Equal(1, len(*parameterSnapshotsToUpdate), "Parameter snapshots length should be 1, but got: %d", len(*parameterSnapshotsToUpdate))
	assert.Equal("Device-2", (*parameterSnapshotsToUpdate)[0].DeviceID, "Device ID in parameter snapshots should be 'Device-2', but got: %s", (*parameterSnapshotsToUpdate)[0].DeviceID)
	assert.Equal("Command-2", (*parameterSnapshotsToUpdate)[0].SDParameter, "Command in parameter snapshots should be 'Command-2', but got: %s", (*parameterSnapshotsToUpdate)[0].SDParameter)
	assert.Nil((*parameterSnapshotsToUpdate)[0].String, "String value in parameter snapshots should be nil, but got: %v", (*parameterSnapshotsToUpdate)[0].String)
	assert.Equal(10.0, *(*parameterSnapshotsToUpdate)[0].Number, "Number value in parameter snapshots should be 10.0, but got: %f", *(*parameterSnapshotsToUpdate)[0].Number)
	assert.Nil((*parameterSnapshotsToUpdate)[0].Boolean, "Boolean value in parameter snapshots should be nil, but got: %v", (*parameterSnapshotsToUpdate)[0].Boolean)
}
