package statements

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
)

type Alert struct {
	Id           string
	Method       string
	Receiver     string
	ReceiverType string
	Message      string
	MessageType  string
}

func (a *Alert) Execute(variableStore *models.VariableStore, _ *models.ReferencedValueStore, _ []models.SDInformationFromBackend) (any, bool, error) {
	services.Logger.Println("Executing alert")

	if a.Method != "phone_number" && a.Method != "email" {
		return nil, false, fmt.Errorf("invalid alert method")
	}

	receiver := a.Receiver
	message := a.Message

	if a.ReceiverType == "variable" {
		variable, err := variableStore.GetVariable(a.Receiver)
		if err != nil {
			return nil, false, err
		}
		receiver = variable.Value.Value.(string)
	}

	if a.MessageType == "variable" {
		variable, err := variableStore.GetVariable(a.Message)
		if err != nil {
			return nil, false, err
		}
		message = variable.Value.Value.(string)
	}

	services.Logger.Println("Sending alert with method", a.Method, "to receiver", receiver, "with message", message)
	return a, true, nil
}

func (a *Alert) GetId() string {
	return a.Id
}

func (a *Alert) GetMethod() string {
	return a.Method
}

func (a *Alert) GetReceiver() (string, string) {
	return a.Receiver, a.ReceiverType
}

func (a *Alert) GetMessage() (string, string) {
	return a.Message, a.MessageType
}

func (a *Alert) Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
