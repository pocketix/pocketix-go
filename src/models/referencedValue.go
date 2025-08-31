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
	DeviceID    uint32          `json:"deviceId"`
	SDParameter uint32          `json:"sdParameter"`
	String      SnapshotString  `json:"string,omitempty"`
	Number      SnapshotNumber  `json:"number,omitempty"`
	Boolean     SnapshotBoolean `json:"boolean,omitempty"`
}

type SnapshotString struct {
	Value string
	Set   bool
}

type SnapshotNumber struct {
	Value float64
	Set   bool
}

type SnapshotBoolean struct {
	Value bool
	Set   bool
}

type ReferencedValue struct {
	DeviceID      uint32 // Device ID
	DeviceUID     string // Device UID
	ParameterID   uint32 // Parameter ID
	ParameterName string // Parameter name
	Type          string // Type of the value
	Value         any    // Value of the referenced value
	IsSet         bool   // Indicates if the referenced value is set
}

type SDInformationFromBackend struct {
	DeviceID  uint32              `json:"deviceId"`            // Device ID
	DeviceUID string              `json:"deviceUID"`           // Device ID
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
		DeviceUID:     deviceID,
		ParameterName: parameterName,
		Type:          "",
		Value:         nil,
	}, true
}

func (rv *ReferencedValue) ToReferenceTarget() string {
	return rv.DeviceUID + "." + rv.ParameterName
}
