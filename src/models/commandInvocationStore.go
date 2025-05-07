package models

type CommandInvocationStore struct {
	Commands []DeviceCommand
}

func NewCommandInvocationStore() *CommandInvocationStore {
	return &CommandInvocationStore{
		Commands: make([]DeviceCommand, 0),
	}
}

func (s *CommandInvocationStore) AddCommand(command DeviceCommand) {
	s.Commands = append(s.Commands, command)
}
