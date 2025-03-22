package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/commands"
	"github.com/stretchr/testify/assert"
)

// MockRepeatExecute is a mock implementation of the Repeat.Execute method.
// This funcion uses the same logic as the original implementation,
// but it has been modified with iterations to make it testable.
func MockRepeatExecute(r commands.Repeat) (bool, int, error) {
	iterations := 0

	for range r.Count {
		iterations++
		result, err := commands.ExecuteCommands(r.Block, nil)
		if err != nil {
			return result, -1, err
		}
	}

	return true, iterations, nil
}

func TestRepeatZeroTimes(t *testing.T) {
	assert := assert.New(t)

	repeatStatement := commands.Repeat{
		Id:    "repeat",
		Count: 0,
		Block: []commands.Command{},
	}

	result, err := repeatStatement.Execute(nil)
	assert.True(result)
	assert.Nil(err)
}

func TestRepeatTenTimes(t *testing.T) {
	assert := assert.New(t)

	repeatStatement := commands.Repeat{
		Id:    "repeat",
		Count: 10,
		Block: []commands.Command{},
	}

	result, iterations, err := MockRepeatExecute(repeatStatement)
	assert.True(result)
	assert.Nil(err)
	assert.Equal(10, iterations)
}

func TestRepeatNegativeCount(t *testing.T) {
	assert := assert.New(t)

	repeatStatement := commands.Repeat{
		Id:    "repeat",
		Count: -1,
		Block: []commands.Command{},
	}

	result, err := repeatStatement.Execute(nil)
	assert.False(result)
	assert.NotNil(err)
	assert.Equal("count cannot be negative", err.Error())
}
