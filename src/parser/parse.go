package parser

import (
	"encoding/json"
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/pocketix/pocketix-go/src/types"
)

func CheckMissingBlock(data []byte) error {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if _, ok := raw["header"]; !ok {
		return fmt.Errorf("header is required")
	}

	if _, ok := raw["block"]; !ok {
		return fmt.Errorf("at least one block is required")
	}
	return nil
}

func ParseHeader(data []byte, variableStore *models.VariableStore, procedureStore *models.ProcedureStore, referencedValueStore *models.ReferencedValueStore) (*types.Program, error) {
	var program types.Program

	if err := CheckMissingBlock(data); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &program); err != nil {
		return nil, err
	}

	if err := ParseVariables(program.Header.Variables, variableStore, referencedValueStore); err != nil {
		return nil, err
	}

	if err := ParseProcedures(program.Header.Procedures, procedureStore); err != nil {
		return nil, err
	}

	return &program, nil
}

func ParseProcedureBody(
	procedure models.Procedure,
	variableStore *models.VariableStore,
	procedureStore *models.ProcedureStore,
	referencedValueStore *models.ReferencedValueStore,
	collector statements.Collector,
) ([]statements.Statement, error) {
	var blocks []types.Block
	if err := json.Unmarshal(procedure.Program, &blocks); err != nil {
		return nil, err
	}

	var commandList []statements.Statement
	for _, block := range blocks {
		statement, err := ParseBlocks(block, variableStore, procedureStore, referencedValueStore, collector)
		if err != nil {
			return nil, err
		}
		commandList = append(commandList, statement...)
	}
	return commandList, nil
}

// Parse parses the given program data and validates the blocks.
//
// Parameters:
//   - data: the program data in JSON format.
//   - variableStore: store for variables to use for parsing.
//   - procedureStore: store for procedure definitions.
//   - referencedValueStore: store for command-related services.
//   - collector: collector for statements.
//
// Returns:
//   - error: nil if parsing was successful, or an error if there was a problem.
func Parse(
	data []byte,
	variableStore *models.VariableStore,
	procedureStore *models.ProcedureStore,
	referencedValueStore *models.ReferencedValueStore,
	collector statements.Collector,
) error {
	program, err := ParseHeader(data, variableStore, procedureStore, referencedValueStore)
	if err != nil {
		return err
	}

	var previousStatement statements.Statement
	for _, block := range program.Blocks {
		subAst := make([]statements.Statement, 0)
		blockCollector := collector.NewCollectorBasedOnType(collector.Type(), &subAst)

		statementList, err := ParseBlocks(block, variableStore, procedureStore, referencedValueStore, blockCollector)
		if err != nil {
			return err
		}
		// If the ParseBlocks returns a list of statements, it means that the block is a procedure call,
		// it appends the statements to the statement list and continues.

		// Initialize device index counter for this block
		deviceIndex := 0

		if len(statementList) != 1 {
			// Process multiple statements (procedure call)
			for _, statement := range statementList {
				// Process the statement using the utility function
				processedStatement, _ := HandleDeviceTypeStatement(statement, block.Devices, &deviceIndex)

				// Collect the processed statement
				collector.Collect(processedStatement)
			}
			continue
		}

		if statementList == nil {
			continue
		}

		// Process single statement
		statement := statementList[0]

		// Check if the statement is a deviceType and replace it if needed
		processedStatement, _ := HandleDeviceTypeStatement(statement, block.Devices, &deviceIndex)

		// Handle if statement with the processed statement
		err = HandleIfStatement(processedStatement, &previousStatement, collector.Collect)
		if err != nil {
			services.Logger.Println("Error handling if statement", err)
			return err
		}
	}
	if previousStatement != nil {
		collector.Collect(previousStatement)
	}

	return nil
}
