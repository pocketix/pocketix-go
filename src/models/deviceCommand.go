package models

import (
	"encoding/json"
	"time"
)

type SDCommand struct {
	CommandDenotation string `json:"commandDenotation"` // Command
	Payload           string `json:"payload"`           // Payload
}

type SDCommandInvocation struct {
	DeviceID          string `json:"deviceId"`          // Device CommandID
	CommandDenotation string `json:"commandDenotation"` // Command
	Payload           string `json:"payload,omitempty"` // Payload
	InvocationTime    string `json:"invocationTime"`    // Invocation time
}

type TypeValue struct {
	Type  string
	Value any
}

type DeviceCommand struct {
	DeviceUID         string
	CommandDenotation string
	Arguments         TypeValue
}

func (dc *DeviceCommand) PrepareCommandToSend(sdInstanceInformation SDInformationFromBackend) (*SDCommandInvocation, error) {
	command := sdInstanceInformation.Command
	var payload map[string]any

	if command.Payload == "" {
		return &SDCommandInvocation{
			DeviceID:          sdInstanceInformation.DeviceUID,
			CommandDenotation: command.CommandDenotation,
			InvocationTime:    time.Now().Format(time.RFC3339),
		}, nil
	}

	err := json.Unmarshal([]byte(command.Payload), &payload)
	if err != nil {
		return nil, err
	}
	// TODO check for possible values
	newPayload := map[string]any{
		"name":  payload["name"],
		"value": dc.Arguments.Value,
	}
	serializedPayload, err := json.Marshal(newPayload)
	if err != nil {
		return nil, err
	}
	return &SDCommandInvocation{
		DeviceID:          sdInstanceInformation.DeviceUID,
		CommandDenotation: command.CommandDenotation,
		Payload:           string(serializedPayload),
		InvocationTime:    time.Now().Format(time.RFC3339),
	}, nil
}
