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

func ParseHeader(data []byte, variableStore *models.VariableStore, procedureStore *models.ProcedureStore, commandHandlingStore *models.CommandsHandlingStore) (*types.Program, error) {
	var program types.Program

	if err := CheckMissingBlock(data); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &program); err != nil {
		return nil, err
	}

	if err := ParseVariables(program.Header.Variables, variableStore, commandHandlingStore); err != nil {
		return nil, err
	}

	if err := ParseProcedures(program.Header.Procedures, procedureStore, commandHandlingStore); err != nil {
		return nil, err
	}

	return &program, nil
}

func ParseProcedureBody(
	procedure models.Procedure,
	variableStore *models.VariableStore,
	procedureStore *models.ProcedureStore,
	commandHandlingStore *models.CommandsHandlingStore,
	collector statements.Collector,
) ([]statements.Statement, error) {
	var blocks []types.Block
	if err := json.Unmarshal(procedure.Program, &blocks); err != nil {
		return nil, err
	}

	var commandList []statements.Statement
	for _, block := range blocks {
		statement, err := ParseBlocks(block, variableStore, procedureStore, commandHandlingStore, collector)
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
//   - commandHandlingStore: store for command-related services.
//   - collector: collector for statements.
//
// Returns:
//   - error: nil if parsing was successful, or an error if there was a problem.
func Parse(
	data []byte,
	variableStore *models.VariableStore,
	procedureStore *models.ProcedureStore,
	commandHandlingStore *models.CommandsHandlingStore,
	collector statements.Collector,
) error {
	program, err := ParseHeader(data, variableStore, procedureStore, commandHandlingStore)
	if err != nil {
		return err
	}

	var previousStatement statements.Statement
	for _, block := range program.Blocks {
		subAst := make([]statements.Statement, 0)
		blockCollector := collector.NewCollectorBasedOnType(collector.Type(), &subAst)

		statementList, err := ParseBlocks(block, variableStore, procedureStore, commandHandlingStore, blockCollector)
		if err != nil {
			return err
		}
		// If the ParseBlocks returns a list of statements, it means that the block is a procedure call,
		// it appends the statements to the statement list and continues.
		if len(statementList) != 1 {
			for _, statement := range statementList {
				collector.Collect(statement)
			}
			continue
		}
		if statementList == nil {
			continue
		}
		statement := statementList[0]

		if statement.GetId() == "if" {
			previousStatement = statement
		} else if statement.GetId() == "else" {
			if previousStatement != nil {
				previousStatement.(*statements.If).AddElseBlock(statement)
				collector.Collect(previousStatement)
				previousStatement = nil
			} else {
				services.Logger.Println("Error: Else without if")
				return fmt.Errorf("else without if")
			}
		} else if statement.GetId() == "elseif" {
			if previousStatement != nil {
				previousStatement.(*statements.If).AddElseIfBlock(statement)
			} else {
				services.Logger.Println("Error: Elseif without if")
				return fmt.Errorf("elseif without if")
			}
		} else {
			if previousStatement != nil {
				collector.Collect(previousStatement)
				previousStatement = nil
			}

			collector.Collect(statement)
		}
	}
	if previousStatement != nil {
		collector.Collect(previousStatement)
	}
	return nil
}
