package models

import (
	"fmt"
	"strings"
)

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
