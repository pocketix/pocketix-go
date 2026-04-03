package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/stretchr/testify/assert"
)

func TestExecuteAlertWebpush(t *testing.T) {
	assert := assert.New(t)

	alert := &statements.Alert{
		Id:        "alert",
		Method:    "WEBPUSH",
		Addressee: "1",
		Message:   "{currentDate} {currentTime}: Alert!!!",
	}

	result, err := alert.Execute(nil, nil, nil, func(any) {})

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")
}

/*func TestExecuteAlertEmail(t *testing.T) {
	assert := assert.New(t)

	alert := &statements.Alert{
		Id:        "alert",
		Method:    "EMAIL",
		Addressee: "1",
		Message:   "Test message",
	}

	result, err := alert.Execute(nil, nil, nil, func(any) {})

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")
}*/

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
		Addressee: "addressee",
		Message:   "Test message",
	}

	result, err := alert.Execute(variableStore, nil, nil, func(any) {})

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	addressee, err := variableStore.GetVariable("addressee")
	assert.Nil(err, "Error should be nil")
	assert.Equal("1234567890", addressee.Value.Value, "Addressee value should be 1234567890")
}

func TestExecuteAlertVariableMessage(t *testing.T) {
	assert := assert.New(t)

	variable := &models.Variable{
		Name:  "message",
		Type:  "string",
		Value: &models.TreeNode{Type: "string", Value: "Test message", ResultValue: "Test message"},
	}
	variableStore := models.NewVariableStore()
	variableStore.AddVariable(*variable)

	alert := &statements.Alert{
		Id:        "alert",
		Method:    "WEBPUSH",
		Addressee: "1234567890",
		Message:   "message",
	}

	result, err := alert.Execute(variableStore, nil, nil, func(any) {})

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	message, err := variableStore.GetVariable("message")
	assert.Nil(err, "Error should be nil")
	assert.Equal("Test message", message.Value.Value, "Message value should be Test message")
}

func TestExecuteAlertInvalidMethod(t *testing.T) {
	assert := assert.New(t)

	alert := &statements.Alert{
		Id:        "alert",
		Method:    "invalid_method",
		Addressee: "1234567890",
		Message:   "Test message",
	}

	result, err := alert.Execute(nil, nil, nil, func(any) {})

	assert.False(result, "Result should be false")
	assert.NotNil(err, "Error should not be nil")
}
