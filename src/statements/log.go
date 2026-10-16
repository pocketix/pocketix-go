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

type Log struct {
	Id          string
	Content     string
	ContentType string
}

func (l *Log) Execute(
	variableStore *models.VariableStore,
	_ *models.ReferencedValueStore,
	_ []types.SDInformationFromBackend,
	callback func(invocation any),
) (bool, error) {
	services.Logger.Println("Executing Log")

	content := l.Content
	if l.ContentType == "variable" {
		variable, err := variableStore.GetVariable(l.Content)
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

	logToSave := types.LogInvocation{
		Content:        content,
		InvocationTime: time.Now().Format(time.RFC3339),
	}
	callback(logToSave)
	return true, nil
}

func (l *Log) GetId() string {
	return l.Id
}

func (l *Log) GetContent() (string, string) {
	return l.Content, l.ContentType
}

func (l *Log) Validate(variableStore *models.VariableStore, referencedValueStore *models.ReferencedValueStore, args ...any) error {
	return nil
}
