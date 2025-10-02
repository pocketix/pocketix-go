package models

import (
	"log"
	"strings"
	"time"

	"github.com/pocketix/pocketix-go/src/types"
	"github.com/pocketix/pocketix-go/src/utils"
)

type DeviceCommand struct {
	DeviceUID         string
	CommandDenotation string
	Arguments         types.TypeValue
}

func (dc *DeviceCommand) PrepareCommandToSend(sdInstanceInformation types.SDInformationFromBackend) (*types.SDCommandInvocation, error) {
	command := sdInstanceInformation.Command

	if command.Payload == "" {
		return createSDCommandInvocationWithoutPayload(sdInstanceInformation, command)
	}
	cleanedPlayload := cleanPayloadString(command.Payload)
	log.Printf("Cleaned Payload: %s", cleanedPlayload)

	payload, err := utils.UnmarshalData[[]types.CommandPayload]([]byte(cleanedPlayload))
	if err != nil {
		return nil, err
	}
	dc.checkPayloadValues(*payload)
	log.Printf("----------------- Command: %+v", command)

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
	return &types.SDCommandInvocation{
		InstanceID:        sdInstanceInformation.DeviceID,
		InstanceUID:       sdInstanceInformation.DeviceUID,
		CommandID:         command.CommandID,
		CommandDenotation: command.CommandDenotation,
		Payload:           command.Payload,
		InvocationTime:    time.Now().Format(time.RFC3339),
	}, nil
}

func createSDCommandInvocationWithoutPayload(sdInstanceInformation types.SDInformationFromBackend, command types.SDCommand) (*types.SDCommandInvocation, error) {
	return &types.SDCommandInvocation{
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

func (dc *DeviceCommand) checkPayloadValues(payload []types.CommandPayload) error {
	for _, p := range payload {
		if p.Type != dc.Arguments.Type {
			return &utils.PayloadTypeMismatchError{
				CommandDenotation: dc.CommandDenotation,
				ExpectedType:      dc.Arguments.Type,
				ActualType:        p.Type,
			}
		}

	}
	return nil
}
