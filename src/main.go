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
	referencedValueStore := models.NewReferencedValueStore()

	commandList, err := parser.Parse(data, variableStore, referencedValueStore)
	if err != nil {
		log.Fatalln(err)
	}

	for _, command := range commandList {
		if _, err := command.Execute(variableStore, referencedValueStore); err != nil {
			log.Fatalln(err)
		}
	}
	services.Logger.Println("Execution completed successfully")
}
