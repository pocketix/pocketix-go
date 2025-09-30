package models

import (
	"strings"
	"time"

	"github.com/pocketix/pocketix-go/src/utils"
)

type SDCommand struct {
	CommandID         uint32 `json:"commandId"`         // Command ID
	CommandDenotation string `json:"commandDenotation"` // Command
	Payload           string `json:"payload"`           // Payload
}

type SDCommandInvocation struct {
	InstanceID        uint32 `json:"instanceId"`        // Instance ID
	InstanceUID       string `json:"instanceUID"`       // Instance ID
	CommandID         uint32 `json:"commandId"`         // Command ID
	CommandDenotation string `json:"commandDenotation"` // Command
	Payload           string `json:"payload,omitempty"` // Payload
	InvocationTime    string `json:"invocationTime"`    // Invocation time
}

type TypeValue struct {
	Type  string
	Value any
}

type CommandPayload struct {
	Name  string `json:"name"`
	Value any    `json:"value"`
}

type DeviceCommand struct {
	DeviceUID         string
	CommandDenotation string
	Arguments         TypeValue
}

func (dc *DeviceCommand) PrepareCommandToSend(sdInstanceInformation SDInformationFromBackend) (*SDCommandInvocation, error) {
	command := sdInstanceInformation.Command

	if command.Payload == "" {
		return createSDCommandInvocationWithoutPayload(sdInstanceInformation, command)
	}
	cleanedPlayload := cleanPayloadString(command.Payload)

	payload, err := utils.UnmarshalData[[]map[string]any]([]byte(cleanedPlayload))
	if err != nil {
		return nil, err
	}

	// err := json.Unmarshal([]byte(cleanedPlayload), &payload)
	// if err != nil {
	// 	return nil, err
	// }
	// // TODO check for possible values
	// newPayload := map[string]any{
	// 	"name":  payload[0]["name"],
	// 	"value": dc.Arguments.Value,
	// }
	// serializedPayload, err := json.Marshal(newPayload)
	// if err != nil {
	// 	return nil, err
	// }
	return &SDCommandInvocation{
		InstanceID:        sdInstanceInformation.DeviceID,
		InstanceUID:       sdInstanceInformation.DeviceUID,
		CommandID:         command.CommandID,
		CommandDenotation: command.CommandDenotation,
		Payload:           command.Payload,
		InvocationTime:    time.Now().Format(time.RFC3339),
	}, nil
}

func createSDCommandInvocationWithoutPayload(sdInstanceInformation SDInformationFromBackend, command SDCommand) (*SDCommandInvocation, error) {
	return &SDCommandInvocation{
		InstanceID:        sdInstanceInformation.DeviceID,
		InstanceUID:       sdInstanceInformation.DeviceUID,
		CommandID:         command.CommandID,
		CommandDenotation: command.CommandDenotation,
		InvocationTime:    time.Now().Format(time.RFC3339),
	}, nil
}

func cleanPayloadString(payload string) string {
	return strings.ReplaceAll(payload, "\n", "")
}
