package models

import (
	"fmt"

	"github.com/pocketix/pocketix-go/src/interfaces"
)

type If struct {
	Id        string               `json:"id"`
	Commands  []interfaces.Command `json:"block"`
	Arguments []Argument           `json:"arguments"`
}

func (i *If) Execute() error {
	return nil
}

func (i *If) HasBlock() bool {
	return true
}

func (i *If) SetBlock(commands []interfaces.Command) {
	i.Commands = append(i.Commands, commands...)
}

func (i *If) String() string {
	return fmt.Sprintf(i.Id+" statement, arguments: %v, body: %v", i.Arguments, i.Commands)
}
