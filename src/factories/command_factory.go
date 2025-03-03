package factories

import (
	"encoding/json"
	"fmt"

	"github.com/pocketix/pocketix-go/src/interfaces"
	"github.com/pocketix/pocketix-go/src/models"
)

func CommandFactory(id string, data json.RawMessage) (interfaces.Command, error) {
	switch id {
	case "if":
		fmt.Println("Creating if command...")
		// var IfStatement models.If
		// if err := json.Unmarshal(data, &IfStatement); err != nil {
		// 	return nil, err
		// }
		return &models.If{Id: id}, nil
	}
	return nil, nil
}
