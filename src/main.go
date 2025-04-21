package main

import (
	"flag"
	"log"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/services"
)

func main() {
	path := flag.String("path", "programs/basic/empty_block.json", "path to the program file")
	flag.Parse()

	data := services.OpenFile(*path)
	variableStore := models.NewVariableStore()
	procedureStore := models.NewProcedureStore()
	commandHandlingStore := models.NewCommandsHandlingStore()

	commandList, err := parser.Parse(data, variableStore, procedureStore, commandHandlingStore)
	if err != nil {
		log.Fatalln(err)
	}

	services.Logger.Println("Commands:", commandHandlingStore.CommandInvocationStore.Commands)

	for _, command := range commandList {
		if _, err := command.Execute(variableStore, commandHandlingStore); err != nil {
			log.Fatalln(err)
		}
	}
	services.Logger.Println("Execution completed successfully")
}
