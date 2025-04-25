package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/utils"
	"github.com/stretchr/testify/assert"
)

func TestToBool(t *testing.T) {
	assert := assert.New(t)

	examples := []any{0, 1, -1, 0.0, 1.0, -1.0, "", "test", true, false, nil, []any{}, map[string]any{}, "true", "false"}
	expected := []bool{false, true, true, false, true, true, false, false, true, false, false, false, false, true, false}
	expectedErr := []bool{false, false, false, false, false, false, true, true, false, false, true, true, true, false, false}

	for i, example := range examples {
		result, err := utils.ToBool(example)
		if expectedErr[i] {
			assert.Error(err, "Example %v should return an error", example)
			continue
		}
		assert.Equal(expected[i], result, "Example %v should be %v", examples[i], expected[i])
	}
}
