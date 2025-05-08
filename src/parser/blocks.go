package parser

import (
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
	referencedValueStore *models.ReferencedValueStore,
	collector statements.Collector,
) ([]statements.Statement, error) {
	argumentTree := make([]*models.TreeNode, len(block.Arguments))

	if len(block.Arguments) != 0 {
		err := ParseArguments(block.Arguments, argumentTree, variableStore, referencedValueStore)
		if err != nil {
			return nil, err
		}
	}

	var previousSubStatement statements.Statement

	for _, subBlock := range block.Body {
		// Parse nested blocks
		subAst := make([]statements.Statement, 0)
		blockCollector := collector.NewCollectorBasedOnType(collector.Type(), &subAst)

		statementList, err := ParseBlocks(subBlock, variableStore, procedureStore, referencedValueStore, blockCollector)
		if err != nil {
			return nil, err
		}
		// Initialize device index counter for this block
		deviceIndex := 0

		if len(statementList) != 1 {
			for _, statement := range statementList {
				// Process the statement using the utility function
				processedStatement, _ := HandleDeviceTypeStatement(statement, block.Devices, &deviceIndex)

				// Collect the processed statement
				collector.Collect(processedStatement)
			}
			continue
		}

		// Process single statement
		statement := statementList[0]

		// Check if the statement is a deviceType and replace it if needed
		processedStatement, _ := HandleDeviceTypeStatement(statement, block.Devices, &deviceIndex)

		// Handle if statement with the processed statement
		err = HandleIfStatement(processedStatement, &previousSubStatement, collector.Collect)
		if err != nil {
			services.Logger.Println("Error handling if statement", err)
			return nil, err
		}
	}

	if previousSubStatement != nil {
		collector.Collect(previousSubStatement)
	}

	if procedureStore != nil && procedureStore.Has(block.Id) {
		procedure := procedureStore.Get(block.Id)
		statementList, err := ParseProcedureBody(procedure, variableStore, procedureStore, referencedValueStore, collector)
		if err != nil {
			return nil, err
		}
		return statementList, nil
	}
	statement, err := statements.StatementFactory(block.Id, *collector.GetTarget(), argumentTree, procedureStore)
	if err != nil {
		services.Logger.Println("Error creating statement", err)
		return nil, err
	}
	if statement == nil {
		services.Logger.Println("Statement is nil, therefore it is device statement")
		return nil, nil
	}

	// Initialize device index counter for this block
	deviceIndex := 0

	// Check if the statement is a deviceType and replace it if needed
	processedStatement, _ := HandleDeviceTypeStatement(statement, block.Devices, &deviceIndex)

	err = processedStatement.Validate(variableStore, referencedValueStore)
	return []statements.Statement{processedStatement}, err
}
