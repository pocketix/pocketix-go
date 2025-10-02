package main

import (
	"flag"
	"log"

	"github.com/pocketix/pocketix-go/src/models"
	"github.com/pocketix/pocketix-go/src/parser"
	"github.com/pocketix/pocketix-go/src/services"
	"github.com/pocketix/pocketix-go/src/statements"
	"github.com/pocketix/pocketix-go/src/types"
)

func MockResolveDeviceInformationFunction(deviceUID string, paramDenotation string, infoType string, deviceCommands *[]types.SDInformationFromBackend) (types.SDInformationFromBackend, error) {
	// Mock implementation: return a dummy SDInformationFromBackend
	return types.SDInformationFromBackend{
		DeviceID:  1,
		DeviceUID: deviceUID,
		Snapshot: types.SDParameterSnapshot{
			SDParameter: 1,
			Number:      types.SnapshotNumber{Value: 42, Set: true},
			String:      types.SnapshotString{Value: "mocked", Set: true},
			Boolean:     types.SnapshotBoolean{Value: true, Set: true},
		},
		Command: types.SDCommand{
			CommandID:         1,
			CommandDenotation: paramDenotation,
			Payload:           "",
		},
	}, nil
}

func main() {
	path := flag.String("path", "programs/basic/empty_block.json", "path to the program file")
	flag.Parse()

	// Load the original program
	data := services.OpenFile(*path)

	// Parse the modified program
	variableStore := models.NewVariableStore()
	procedureStore := models.NewProcedureStore()
	referencedValueStore := models.NewReferencedValueStore()
	referencedValueStore.SetResolveParameterFunction(MockResolveDeviceInformationFunction)

	ast := make([]statements.Statement, 0)
	collector := &statements.ASTCollector{Target: &ast}
	// err = parser.Parse(modifiedData, variableStore, procedureStore, referencedValueStore, &statements.NoOpCollector{})

	err := parser.Parse(data, variableStore, procedureStore, referencedValueStore, collector)

	if err != nil {
		log.Fatalln(err)
	}

	var interpretInvocationsToSend []types.SDCommandInvocation
	for _, block := range ast {
		if _, err := block.Execute(variableStore, referencedValueStore, collector.DeviceCommands, func(deviceCommand types.SDCommandInvocation) {
			interpretInvocationsToSend = append(interpretInvocationsToSend, deviceCommand)
		}); err != nil {
			log.Fatalln(err)
		}
	}
	services.Logger.Println("Execution completed successfully")
}
