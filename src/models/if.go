package models

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/interfaces"
)

type If struct {
	Id        string
	Block     []interfaces.Command
	Arguments []Argument
}

func (i *If) Execute() error {
	return nil
}

func (i *If) String() string {
	return fmt.Sprintf("If, Body: %v, Arguments: %v", i.Block, i.Arguments)
}
