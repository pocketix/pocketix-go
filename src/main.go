package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/statements"
)

func main() {
	path := flag.String("path", "programs/basic/empty_block.json", "path to the program file")
	flag.Parse()

	// Load the original program
	data := services.OpenFile(*path)

	// Simulate procedures coming from backend/RabbitMQ
	// This would be the procedures[] entry from the payload
	testProcedures := []byte(`{
		"testProcedure": [
			{
				"id": "alert",
				"arguments": [
					{
						"type": "str_opt",
						"value": "phone_number"
					},
					{
						"type": "string",
						"value": "123456789"
					},
					{
						"type": "string",
						"value": "This is a test procedure"
					}
				]
			},
			{
				"id": "alert",
				"arguments": [
					{
						"type": "str_opt",
						"value": "email"
					},
					{
						"type": "string",
						"value": "test@example.com"
					},
					{
						"type": "string",
						"value": "Email from procedure"
					}
				]
			}
		]
	}`)

	// Print original program for debugging
	fmt.Println("Original Program:")
	var originalProgram map[string]interface{}
	json.Unmarshal(data, &originalProgram)
	originalJSON, _ := json.MarshalIndent(originalProgram, "", "  ")
	fmt.Println(string(originalJSON))

	// Add procedures to the program
	modifiedData, err := services.AddProceduresToProgram(data, testProcedures)
	if err != nil {
		log.Fatalf("Failed to add procedures to program: %v", err)
	}

	// Print modified program for debugging
	fmt.Println("\nModified Program with Added Procedures:")
	var modifiedProgram map[string]interface{}
	json.Unmarshal(modifiedData, &modifiedProgram)
	modifiedJSON, _ := json.MarshalIndent(modifiedProgram, "", "  ")
	fmt.Println(string(modifiedJSON))

	// Parse the modified program
	variableStore := models.NewVariableStore()
	procedureStore := models.NewProcedureStore()
	commandHandlingStore := models.NewCommandsHandlingStore()

	// ast := make([]commands.Command, 0)
	err = parser.Parse(modifiedData, variableStore, procedureStore, commandHandlingStore, &statements.NoOpCollector{})

	// err := parser.Parse(data, variableStore, procedureStore, commandHandlingStore, func(c commands.Command) {})

	if err != nil {
		log.Fatalln(err)
	}

	// Check if the procedure was added to the store
	fmt.Println("\nProcedure Store Contents:")
	for name := range procedureStore.Procedures {
		fmt.Printf("- Procedure: %s\n", name)
	}

	services.Logger.Println("Commands:", commandHandlingStore.CommandInvocationStore.Commands)

	// for _, block := range ast {
	// 	if _, err := block.Execute(variableStore, commandHandlingStore); err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }
	services.Logger.Println("Execution completed successfully")
}
