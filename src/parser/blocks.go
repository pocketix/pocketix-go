package parser

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/pocketix/pocketix-go/src/types"
)

// ParseBlocks parses a block of code.
//
// Parameters:
//   - block: the block to parse.
//   - variableStore: store for variables to use for parsing.
//   - procedureStore: store for procedure definitions.
//   - commandHandlingStore: store for command handling.
//   - collector: collector for statements.
//
// Returns:
//   - an AST of statements.
//   - an error if any occurred during parsing.
func ParseBlocks(
	block types.Block,
	variableStore *models.VariableStore,
	procedureStore *models.ProcedureStore,
	commandHandlingStore *models.CommandsHandlingStore,
	collector statements.Collector,
) ([]statements.Statement, error) {
	argumentTree := make([]*models.TreeNode, len(block.Arguments))

	if len(block.Arguments) != 0 {
		err := ParseArguments(block.Arguments, argumentTree, variableStore, commandHandlingStore)
		if err != nil {
			return nil, err
		}
	}

	var previousSubStatement statements.Statement

	for _, subBlock := range block.Body {
		// Parse nested blocks
		subAst := make([]statements.Statement, 0)
		blockCollector := collector.NewCollectorBasedOnType(collector.Type(), &subAst)

		statementList, err := ParseBlocks(subBlock, variableStore, procedureStore, commandHandlingStore, blockCollector)
		if err != nil {
			return nil, err
		}
		if len(statementList) != 1 {
			for _, statement := range statementList {
				collector.Collect(statement)
			}
			continue
		}
		statement := statementList[0]

		if statement.GetId() == "if" {
			previousSubStatement = statement
		} else if statement.GetId() == "else" {
			if previousSubStatement != nil {
				previousSubStatement.(*statements.If).AddElseBlock(statement)
				collector.Collect(previousSubStatement)
				previousSubStatement = nil
			} else {
				services.Logger.Println("Error: Else without if")
				return nil, fmt.Errorf("else without if")
			}
		} else if statement.GetId() == "elseif" {
			if previousSubStatement != nil {
				previousSubStatement.(*statements.If).AddElseIfBlock(statement)
			} else {
				services.Logger.Println("Error: Elseif without if")
				return nil, fmt.Errorf("elseif without if")
			}
		} else {
			if previousSubStatement != nil {
				collector.Collect(previousSubStatement)
				previousSubStatement = nil
			}

			collector.Collect(statement)
		}
	}

	if previousSubStatement != nil {
		collector.Collect(previousSubStatement)
	}

	if procedureStore != nil && procedureStore.Has(block.Id) {
		procedure := procedureStore.Get(block.Id)
		statementList, err := ParseProcedureBody(procedure, variableStore, procedureStore, commandHandlingStore, collector)
		if err != nil {
			return nil, err
		}
		return statementList, nil
	}
	statement, err := statements.StatementFactory(block.Id, *collector.GetTarget(), argumentTree, procedureStore, commandHandlingStore.CommandInvocationStore)
	if err != nil {
		services.Logger.Println("Error creating statement", err)
		return nil, err
	}
	if statement == nil {
		services.Logger.Println("Statement is nil, therefore it is device statement")
		return nil, nil
	}
	err = statement.Validate(variableStore, commandHandlingStore.ReferencedValueStore)
	return []statements.Statement{statement}, err
}
