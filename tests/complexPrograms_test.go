package tests

import (
	"testing"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/stretchr/testify/assert"
)

func TestExecuteWhileSetVar(t *testing.T) {
	assert := assert.New(t)

	data := services.OpenFile("../programs/complex/prog2.json")
	variableStore := models.NewVariableStore()
	procedureStore := models.NewProcedureStore()
	referencedValueStore := models.NewReferencedValueStore()

	statementAST := make([]statements.Statement, 0)
	err := parser.Parse(data, variableStore, procedureStore, referencedValueStore, &statements.ASTCollector{Target: &statementAST})
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.NotNil(statementAST, "Commands list should not be nil")

	for _, statement := range statementAST {
		_, err := statement.Execute(variableStore, referencedValueStore, nil, nil)
		assert.Nil(err, "Error should be nil, but got: %v", err)
	}

	variable, err := variableStore.GetVariable("count")
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.Equal(float64(5), variable.Value.Value, "Variable value should be 5, but got: %v", variable.Value.Value)
}

func MockResolveParameterFunctionComplexProgram(deviceUID string, paramDenotation string, infoType string, deviceCommands *[]models.SDInformationFromBackend) (models.SDInformationFromBackend, error) {
	return models.SDInformationFromBackend{
		DeviceUID: deviceUID,
		Snapshot: models.SDParameterSnapshot{
			SDParameter: paramDenotation,
			Number:      func(v float64) *float64 { return &v }(230.0),
		},
	}, nil
}

func TestExecuteProgramWithReferencedValue(t *testing.T) {
	assert := assert.New(t)

	data := services.OpenFile("../programs/complex/prog5.json")
	variableStore := models.NewVariableStore()
	procedureStore := models.NewProcedureStore()
	referencedValueStore := models.NewReferencedValueStore()
	referencedValueStore.SetResolveParameterFunction(MockResolveParameterFunctionComplexProgram)
	statementAST := make([]statements.Statement, 0)

	err := parser.Parse(data, variableStore, procedureStore, referencedValueStore, &statements.ASTCollector{Target: &statementAST})
	assert.Nil(err, "Error should be nil, but got: %v", err)
	assert.NotNil(statementAST, "Commands list should not be nil")

	for _, statement := range statementAST {
		_, err := statement.Execute(variableStore, referencedValueStore, nil, nil)
		assert.Nil(err, "Error should be nil, but got: %v", err)
	}

	referencedValue, ok := referencedValueStore.GetReferencedValueFromStore("DistanceSensor-1.waterLevel")
	assert.True(ok, "Referenced value should be found, but got: %v", ok)
	assert.Equal("DistanceSensor-1", referencedValue.DeviceID, "Device ID should be DistanceSensor-1, but got: %v", referencedValue.DeviceID)
	assert.Equal("waterLevel", referencedValue.ParameterName, "Parameter name should be waterLevel, but got: %v", referencedValue.ParameterName)
	assert.Equal("number", referencedValue.Type, "Type should be number, but got: %v", referencedValue.Type)
	assert.Equal(230.0, referencedValue.Value, "Value should be 230.0, but got: %v", referencedValue.Value)
}
