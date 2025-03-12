package models

import (
	"github.com/pocketix/pocketix-go/src/interfaces"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/tree"
)

type If struct {
	Id        string
	Block     []interfaces.Command
	Arguments *tree.TreeNode
}

func (i *If) Execute() error {
	services.Logger.Println("Executing if")
	if result, err, _ := i.Arguments.Evaluate(); err != nil {
		services.Logger.Println("Error executing if arguments", err)
	} else {
		if result.(bool) {
			services.Logger.Println("If is true, can execute body")
		} else {
			services.Logger.Println("If is false, skip body")
		}
	}

	return nil
}

// func (i *If) String() string {
// 	services.Logger.Println("Printing If")
// 	return fmt.Sprintf("If, Body: %v, Arguments: %v", i.Block, i.Arguments)
// }
