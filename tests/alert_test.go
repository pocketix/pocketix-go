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
		Content:   "{currentDate} {currentTime}: Alert!!!",
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
		Content:   "Test content",
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
		Content:   "content",
	}

	result, err := alert.Execute(variableStore, nil, nil, func(any) {})

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	content, err := variableStore.GetVariable("content")
	assert.Nil(err, "Error should be nil")
	assert.Equal("Test content", content.Value.Value, "Content value should be Test content")
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
