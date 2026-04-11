package tests

import (
	"regexp"
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/pocketix/pocketix-go/src/types"
	"github.com/stretchr/testify/assert"
)

func TestExecuteAlertWebpush(t *testing.T) {
	assert := assert.New(t)

	alert := &statements.Alert{
		Id:        "alert",
		Method:    "WEBPUSH",
		Addressee: "1",
		Content:   "message",
	}

	var invocation any
	result, err := alert.Execute(nil, nil, nil, func(inv any) {
		invocation = inv
	})

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	// Check that the invocation is a NotificationInvocation
	_, ok := invocation.(types.NotificationInvocation)
	assert.True(ok, "Invocation should be a NotificationInvocation")

}

func TestExecuteAlertVariableAddressee(t *testing.T) {
	assert := assert.New(t)

	variable := models.Variable{
		Name:  "addressee",
		Type:  "string",
		Value: &models.TreeNode{Type: "string", Value: "1234567890", ResultValue: "1234567890"},
	}
	variableStore := models.NewVariableStore()
	variableStore.AddVariable(variable)

	alert := &statements.Alert{
		Id:        "alert",
		Method:    "WEBPUSH",
		Addressee: "1",
		Content:   "Test content",
	}

	result, err := alert.Execute(variableStore, nil, nil, func(any) {})

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	addressee, err := variableStore.GetVariable("addressee")
	assert.Nil(err, "Error should be nil")
	assert.Equal("1234567890", addressee.Value.Value, "Addressee value should be 1234567890")
}

func TestExecuteAlertVariableContent(t *testing.T) {
	assert := assert.New(t)

	variable := &models.Variable{
		Name:  "content",
		Type:  "string",
		Value: &models.TreeNode{Type: "string", Value: "Test content", ResultValue: "Test content"},
	}
	variableStore := models.NewVariableStore()
	variableStore.AddVariable(*variable)

	alert := &statements.Alert{
		Id:        "alert",
		Method:    "WEBPUSH",
		Addressee: "1234567890",
		Content:   "{currentDate} {currentTime}: Alert!!!",
	}

	var invocation any
	result, err := alert.Execute(nil, nil, nil, func(inv any) {
		invocation = inv
	})

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	notificationInvocation, ok := invocation.(types.NotificationInvocation)
	assert.True(ok, "Invocation should be a NotificationInvocation")

	assert.NotContains(notificationInvocation.Content, "{", "Content should not contain unformatted placeholders")
	assert.NotContains(notificationInvocation.Content, "}", "Content should not contain unformatted placeholders")

	dateTimeRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}: Alert!!!$`)
	assert.True(dateTimeRegex.MatchString(notificationInvocation.Content), "Content should match the expected date and time format")

}

func TestExecuteAlertInvalidMethod(t *testing.T) {
	assert := assert.New(t)

	alert := &statements.Alert{
		Id:        "alert",
		Method:    "invalid_method",
		Addressee: "1234567890",
		Content:   "Test content",
	}

	result, err := alert.Execute(nil, nil, nil, func(any) {})

	assert.False(result, "Result should be false")
	assert.NotNil(err, "Error should not be nil")
}
