package models

import (
	"fmt"
)

type TypeValue struct {
	Type  string
	Value any
}

type DeviceCommand struct {
	DeviceID  string
	Command   string
	Arguments []TypeValue // list of arguments, for now only one argument is assumed
}

func (dc *DeviceCommand) SendCommandToDevice() (*DeviceCommand, *SDParameterSnapshot, error) {
	if len(dc.Arguments) == 0 {
		return dc, nil, nil
	}
	var sdParameterSnapshot SDParameterSnapshot
	argument := dc.Arguments[0]

	switch argument.Type {
	case "string":
		value := argument.Value.(string)
		sdParameterSnapshot.String = &value
	case "number":
		value := argument.Value.(float64)
		sdParameterSnapshot.Number = &value
	case "boolean":
		value := argument.Value.(bool)
		argument.Value = value
	default:
		return nil, nil, fmt.Errorf("unsupported type: %s", argument.Type)
	}

	sdParameterSnapshot.DeviceID = dc.DeviceID
	sdParameterSnapshot.SDParameter = dc.Command
	return dc, &sdParameterSnapshot, nil
}
