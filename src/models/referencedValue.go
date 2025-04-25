package models

import (
	"fmt"
	"strings"
)

type ReferenceValueResponseFromBackend struct {
	DeviceID             string                `json:"deviceId"`             // Device ID
	SDType               SDType                `json:"sdType"`               // SD type
	SDParameterSnapshots []SDParameterSnapshot `json:"sdParameterSnapshots"` // List of SD parameter snapshots
}

type SDType struct {
	SDParameters []SDParameter `json:"sdParameters"` // List of SD parameters
	SDCommands   []SDCommand   `json:"sdCommands"`   // List of SD commands
}

type SDParameter struct {
	ParameterID         uint32 `json:"parameterName"`       // Parameter name
	ParameterDenotation string `json:"parameterDenotation"` // Parameter denotation
}

type SDCommand struct {
	CommandID         uint32 `json:"deviceId"`          // Device CommandID
	CommandDenotation string `json:"commandDenotation"` // Command denotation
}

type SDParameterSnapshot struct {
	DeviceID    string   `json:"deviceId"`
	SDParameter string   `json:"sdParameter"`
	String      *string  `json:"string,omitempty"`
	Number      *float64 `json:"number,omitempty"`
	Boolean     *bool    `json:"boolean,omitempty"`
}

type ReferencedValue struct {
	DeviceID      string // Device ID
	ParameterName string // Parameter name
	Type          string // Type of the value
	Value         any    // Value of the referenced value
}

func FromReferencedTarget(referencedTarget string) (string, string, error) {
	parts := strings.Split(referencedTarget, ".")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid referenced target: %s", referencedTarget)
	}
	return parts[0], parts[1], nil
}

func NewReferencedValue(referencedTarget string) (*ReferencedValue, error) {
	deviceID, parameterName, err := FromReferencedTarget(referencedTarget)
	if err != nil {
		return nil, err
	}
	return &ReferencedValue{
		DeviceID:      deviceID,
		ParameterName: parameterName,
		Type:          "",
		Value:         nil,
	}, nil
}

func (rv *ReferencedValue) ToReferenceTarget() string {
	return rv.DeviceID + "." + rv.ParameterName
}
