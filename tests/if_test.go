package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfWithoutBodyAndArguments(t *testing.T) {
	assert := assert.New(t)
	program := LoadProgram(t, "../programs/if/simple_if.json")

	assert.Equal(1, len(program.Blocks), "Expected 1 block, got %d", len(program.Blocks))

	ifStatement := program.Blocks[0]

	assert.NotNil(ifStatement)
	assert.Equal("if", ifStatement.Id, "Expected 'if' command, got %s", ifStatement.Id)

	assert.Equal(0, len(ifStatement.Body), "Expected 0 commands, got %d", len(ifStatement.Body))
	assert.Equal(1, len(ifStatement.Arguments), "Expected 1 arguments, got %d", len(ifStatement.Arguments))
}

func TestIfTwoInRow(t *testing.T) {
	assert := assert.New(t)
	program := LoadProgram(t, "../programs/if/two_ifs.json")

	assert.Equal(2, len(program.Blocks), "Expected 2 blocks, got %d", len(program.Blocks))
}

func TestNestedIf(t *testing.T) {
	assert := assert.New(t)
	program := LoadProgram(t, "../programs/if/nested_if.json")

	assert.Equal(1, len(program.Blocks), "Expected 1 block, got %d", len(program.Blocks))

	ifStatement := program.Blocks[0]
	nestedIfStatement := ifStatement.Body[0]

	assert.NotNil(ifStatement)
	assert.NotNil(nestedIfStatement)

	assert.Equal("if", ifStatement.Id, "Expected 'if' command, got %s", ifStatement.Id)
	assert.Equal("if", nestedIfStatement.Id, "Expected 'if' command, got %s", nestedIfStatement.Id)
}
