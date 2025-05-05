package models

type CommandsHandlingStore struct {
	ReferencedValueStore   *ReferencedValueStore
	CommandInvocationStore *CommandInvocationStore
}

func NewCommandsHandlingStore() *CommandsHandlingStore {
	return &CommandsHandlingStore{
		ReferencedValueStore:   NewReferencedValueStore(),
		CommandInvocationStore: NewCommandInvocationStore(),
	}
}
