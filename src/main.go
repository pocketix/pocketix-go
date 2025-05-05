package main

import (
	"flag"
	"log"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/statements"
)

func main() {
	path := flag.String("path", "programs/basic/empty_block.json", "path to the program file")
	flag.Parse()

	data := services.OpenFile(*path)
	variableStore := models.NewVariableStore()
	procedureStore := models.NewProcedureStore()
	commandHandlingStore := models.NewCommandsHandlingStore()

	// ast := make([]commands.Command, 0)
	err := parser.Parse(data, variableStore, procedureStore, commandHandlingStore, &statements.NoOpCollector{})

	// err := parser.Parse(data, variableStore, procedureStore, commandHandlingStore, func(c commands.Command) {})

	if err != nil {
		log.Fatalln(err)
	}

	services.Logger.Println("Commands:", commandHandlingStore.CommandInvocationStore.Commands)

	// for _, block := range ast {
	// 	if _, err := block.Execute(variableStore, commandHandlingStore); err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }
	services.Logger.Println("Execution completed successfully")
}
