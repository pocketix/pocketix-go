package statements

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/types"
)

type Alert struct {
	Id            string
	Method        string
	Addressee     string
	AddresseeType string
	Content       string
	ContentType   string
}

func (a *Alert) Execute(
	variableStore *models.VariableStore,
	_ *models.ReferencedValueStore,
	_ []types.SDInformationFromBackend,
	callback func(invocation any),
) (bool, error) {
	services.Logger.Println("Executing alert")

	if a.Method != "WEBPUSH" {
		return false, fmt.Errorf("invalid alert method")
	}

	addressee := a.Addressee
	content := a.Content

	if a.AddresseeType == "variable" {
		variable, err := variableStore.GetVariable(a.Addressee)
		if err != nil {
			return false, err
		}
		addressee = variable.Value.Value.(string)
	}

	if a.ContentType == "variable" {
		variable, err := variableStore.GetVariable(a.Content)
		if err != nil {
			return false, err
		}
		content = variable.Value.Value.(string)
	}

	dynamicValues := map[string]string{
		"currentTime": time.Now().Format("15:04:05"),
		"currentDate": time.Now().Format("2006-01-02"),
	}
	re := regexp.MustCompile(`\{([^{}]+)\}`)
	content = re.ReplaceAllStringFunc(content, func(match string) string {
		matchTrimmed := strings.TrimSpace(match[1 : len(match)-1])
		if len(matchTrimmed) > 0 {
			if matchTrimmed[0] == '$' {
				variableName := matchTrimmed[1:]
				variable, err := variableStore.GetVariable(variableName)
				if err == nil && variable != nil {
					return fmt.Sprint(variable.Value.Value)
				}
			} else {
				if value, ok := dynamicValues[matchTrimmed]; ok {
					return value
				}
			}
		}
		return match
	})

	notificationToSend := types.NotificationInvocation{
		AddresseeID:    addressee,
		Content:        content,
		EndpointType:   a.Method,
		InvocationTime: time.Now().Format(time.RFC3339),
	}
	callback(notificationToSend)
	return true, nil
}

func (a *Alert) GetId() string {
	return a.Id
}

func (a *Alert) GetMethod() string {
	return a.Method
}

func (a *Alert) GetAddressee() (string, string) {
	return a.Addressee, a.AddresseeType
}

func (a *Alert) GetContent() (string, string) {
	return a.Content, a.ContentType
}

func (a *Alert) Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
