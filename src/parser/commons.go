package parser

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/statements"
)

func HandleIfStatement(
	statement statements.Statement,
	previousStatement *statements.Statement,
	collect func(statement statements.Statement),
) error {
	if statement.GetId() == "if" {
		*previousStatement = statement
	} else if statement.GetId() == "else" {
		if *previousStatement != nil {
			(*previousStatement).(*statements.If).AddElseBlock(statement)
			collect(*previousStatement)
			*previousStatement = nil
		} else {
			services.Logger.Println("Error: Else without if")
			return fmt.Errorf("else without if")
		}
	} else if statement.GetId() == "elseif" {
		if *previousStatement != nil {
			(*previousStatement).(*statements.If).AddElseIfBlock(statement)
		} else {
			services.Logger.Println("Error: Elseif without if")
			return fmt.Errorf("elseif without if")
		}
	} else {
		if *previousStatement != nil {
			collect(*previousStatement)
			*previousStatement = nil
		}

		collect(statement)
	}
	return nil
}
