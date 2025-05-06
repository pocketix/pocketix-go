package parser

import (
	"encoding/json"

	"github.com/pocketix/pocketix-go/src/models"
)

func ParseProcedures(data json.RawMessage, procedureStore *models.ProcedureStore) error {
	var procedures map[string]json.RawMessage

	if err := json.Unmarshal(data, &procedures); err != nil {
		return err
	}

	for name, program := range procedures {
		procedure := models.Procedure{
			Name:    name,
			Program: program,
		}
		procedureStore.AddProcedure(procedure)
	}

	return nil
}
