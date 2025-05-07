package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/stretchr/testify/assert"
)

func TestExecuteAlertPhoneNumber(t *testing.T) {
	assert := assert.New(t)

	alert := &statements.Alert{
		Id:           "alert",
		Method:       "phone_number",
		Receiver:     "1234567890",
		ReceiverType: "phone_number",
		Message:      "Test message",
		MessageType:  "string",
	}

	result, err := alert.Execute(nil, nil)

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")
}

func TestExecuteAlertEmail(t *testing.T) {
	assert := assert.New(t)

	alert := &statements.Alert{
		Id:           "alert",
		Method:       "email",
		Receiver:     "mail@mail.com",
		ReceiverType: "email",
		Message:      "Test message",
		MessageType:  "string",
	}

	result, err := alert.Execute(nil, nil)

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")
}

func TestExecuteAlertVariableReceiver(t *testing.T) {
	assert := assert.New(t)

	variable := models.Variable{
		Name:  "receiver",
		Type:  "string",
		Value: &models.TreeNode{Type: "string", Value: "1234567890", ResultValue: "1234567890"},
	}
	variableStore := models.NewVariableStore()
	variableStore.AddVariable(variable)

	alert := &statements.Alert{
		Id:           "alert",
		Method:       "phone_number",
		Receiver:     "receiver",
		ReceiverType: "variable",
		Message:      "Test message",
		MessageType:  "string",
	}

	result, err := alert.Execute(variableStore, nil)

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	receiver, err := variableStore.GetVariable("receiver")
	assert.Nil(err, "Error should be nil")
	assert.Equal("1234567890", receiver.Value.Value, "Receiver value should be 1234567890")
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
		Id:           "alert",
		Method:       "phone_number",
		Receiver:     "1234567890",
		ReceiverType: "phone_number",
		Message:      "message",
		MessageType:  "variable",
	}

	result, err := alert.Execute(variableStore, nil)

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	message, err := variableStore.GetVariable("message")
	assert.Nil(err, "Error should be nil")
	assert.Equal("Test message", message.Value.Value, "Message value should be Test message")
}

func TestExecuteAlertInvalidMethod(t *testing.T) {
	assert := assert.New(t)

	alert := &statements.Alert{
		Id:           "alert",
		Method:       "invalid_method",
		Receiver:     "1234567890",
		ReceiverType: "phone_number",
		Message:      "Test message",
		MessageType:  "string",
	}

	result, err := alert.Execute(nil, nil)

	assert.False(result, "Result should be false")
	assert.NotNil(err, "Error should not be nil")
}
