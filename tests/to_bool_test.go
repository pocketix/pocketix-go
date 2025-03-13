package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestToBool(t *testing.T) {
	assert := assert.New(t)

	examples := []any{0, 1, -1, 0.0, 1.0, -1.0, "", "test", true, false, nil, []any{}, map[string]any{}}
	expected := []bool{false, true, true, false, true, true, false, true, true, false, false, false, false}

	services.Logger.Println(len(examples), len(expected))
	for i, example := range examples {
		assert.Equal(expected[i], utils.ToBool(example), "Example %v should be %v", examples[i], expected[i])
	}
}
