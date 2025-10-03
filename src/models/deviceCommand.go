package models

import (
	"log"
	"regexp"
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

	cleanedPayload := normalizePossibleValues(command.Payload)
	cleanedPayload = cleanPayloadString(cleanedPayload)
	log.Printf("Cleaned Payload: %s", cleanedPayload)

	payload, err := utils.UnmarshalData[[]types.CommandPayload]([]byte(cleanedPayload))
	if err != nil {
		return nil, err
	}

	payloadErr := dc.checkPayloadValues(*payload)
	if payloadErr != nil {
		return nil, payloadErr
	}

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

func normalizePossibleValues(payload string) string {
	re := regexp.MustCompile(`"possibleValues"\s*:\s*""`)
	return re.ReplaceAllString(payload, `"possibleValues":[]`)
}

func (dc *DeviceCommand) checkPayloadValues(payload []types.CommandPayload) error {
	for _, p := range payload {
		if p.Type != dc.Arguments.Type {
			return utils.NewErrorOf[utils.PayloadTypeMismatchError](dc.CommandDenotation, p.Type, dc.Arguments.Type)
		}
		if len(p.Values) > 0 {
			valueFound := false
			for _, v := range p.Values {
				if v == dc.Arguments.Value {
					valueFound = true
					break
				}
			}
			if !valueFound {
				return utils.NewErrorOf[utils.PayloadValueMissingError](dc.CommandDenotation, p.Values, dc.Arguments.Value)
			}
		}
	}
	return nil
}
