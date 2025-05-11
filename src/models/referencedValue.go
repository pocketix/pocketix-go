package models

import (
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

type SDParameterSnapshot struct {
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
	IsSet         bool   // Indicates if the referenced value is set
}

type SDInformationFromBackend struct {
	DeviceUID string              `json:"deviceId"`            // Device ID
	Snapshot  SDParameterSnapshot `json:"sdParameterSnapshot"` // SD parameter snapshot
	Command   SDCommand           `json:"sdCommand"`           // SD command
}

func FromReferencedTarget(referencedTarget string) (string, string, bool) {
	lastDot := strings.LastIndex(referencedTarget, ".")
	if lastDot == -1 {
		return "", "", false
	}

	prefix := referencedTarget[:lastDot]
	last := referencedTarget[lastDot+1:]

	if prefix == "" || last == "" {
		return "", "", false
	}

	return prefix, last, true
}

func NewReferencedValue(referencedTarget string) (*ReferencedValue, bool) {
	deviceID, parameterName, ok := FromReferencedTarget(referencedTarget)
	if !ok {
		return nil, ok
	}
	return &ReferencedValue{
		DeviceID:      deviceID,
		ParameterName: parameterName,
		Type:          "",
		Value:         nil,
	}, true
}

func (rv *ReferencedValue) ToReferenceTarget() string {
	return rv.DeviceID + "." + rv.ParameterName
}
