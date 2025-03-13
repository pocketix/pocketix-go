package tests

import (
	"io"
	"testing"

	"github.com/pocketix/pocketix-go/src/services"
	// 	"github.com/pocketix/pocketix-go/src/models"
	// 	"github.com/pocketix/pocketix-go/src/parser"
	// "github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	services.SetOutput(io.Discard)
	m.Run()
}

// func LoadProgram(t *testing.T, filePath string) *models.Program {
// 	t.Helper()

// 	data := services.OpenFile(filePath)
// 	assert.NotNil(t, data, "File should not be empty")

// 	program, err := parser.Parse(data)
// 	assert.Nil(t, err, "Parsing should not return an error")
// 	assert.NotNil(t, program, "Parsed program should not be nil")

// 	return program
// }
