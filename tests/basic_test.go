package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/services"
)

func TestBasicEmptyBlock(t *testing.T) {
	assert := assert.New(t)
	program := LoadProgram(t, "../programs/basic/empty_block.json")

	assert.Equal(0, len(program.Blocks), "Expected 0 blocks, got %d", len(program.Blocks))
}

func TestBasicWithoutFirstBlock(t *testing.T) {
	data := services.OpenFile("../programs/basic/no_first_block.json")
	_, err := parser.Parse(data)

	assert.NotNil(t, err, "Expected error when parsing program without first block, got nil")
}
