package models

import (
	"strings"
)

type ReferencedValue struct {
	DeviceID      uint32 // Device ID
	DeviceUID     string // Device UID
	ParameterID   uint32 // Parameter ID
	ParameterName string // Parameter name
	Type          string // Type of the value
	Value         any    // Value of the referenced value
	IsSet         bool   // Indicates if the referenced value is set
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
