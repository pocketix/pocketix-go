package models

import (
	"github.com/pocketix/pocketix-go/src/interfaces"
)

type Block struct {
	Id   string               `json:"id"`
	Body []interfaces.Command `json:"block"`
}
