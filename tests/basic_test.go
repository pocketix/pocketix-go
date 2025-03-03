package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/services"
)

func TestBasicEmptyBlock(t *testing.T) {
	data := services.OpenFile("../programs/basic/empty_block.json")
	program, err := parser.Parse(data)

	if err != nil {
		t.Error(err)
	}

	if len(program.Blocks) != 0 {
		t.Errorf("Expected 0 blocks, got %d", len(program.Blocks))
	}
}

func TestBasicWithoutFirstBlock(t *testing.T) {
	data := services.OpenFile("../programs/basic/no_first_block.json")
	_, err := parser.Parse(data)

	if err == nil {
		t.Error("Expected error when parsing program without first block, got nil")
	}
}
