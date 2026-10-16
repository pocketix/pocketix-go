package tests

import (
	"regexp"
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/pocketix/pocketix-go/src/types"
	"github.com/stretchr/testify/assert"
)

func TestExecuteLog(t *testing.T) {
	assert := assert.New(t)

	log := &statements.Log{
		Id:      "log",
		Content: "message",
	}

	var invocation any
	result, err := log.Execute(nil, nil, nil, func(inv any) {
		invocation = inv
	})

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	// Check that the invocation is a LogInvocation
	_, ok := invocation.(types.LogInvocation)
	assert.True(ok, "Invocation should be a LogInvocation")

}

func TestExecuteLogVariableContent(t *testing.T) {
	assert := assert.New(t)

	variable := &models.Variable{
		Name:  "content",
		Type:  "string",
		Value: &models.TreeNode{Type: "string", Value: "Test content", ResultValue: "Test content"},
	}
	variableStore := models.NewVariableStore()
	variableStore.AddVariable(*variable)

	log := &statements.Log{
		Id:      "log",
		Content: "{currentDate} {currentTime}: Log!!!",
	}

	var invocation any
	result, err := log.Execute(nil, nil, nil, func(inv any) {
		invocation = inv
	})

	assert.True(result, "Result should be true")
	assert.Nil(err, "Error should be nil")

	logInvocation, ok := invocation.(types.LogInvocation)
	assert.True(ok, "Invocation should be a LogInvocation")

	assert.NotContains(logInvocation.Content, "{", "Content should not contain unformatted placeholders")
	assert.NotContains(logInvocation.Content, "}", "Content should not contain unformatted placeholders")

	dateTimeRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}: Log!!!$`)
	assert.True(dateTimeRegex.MatchString(logInvocation.Content), "Content should match the expected date and time format")

}
