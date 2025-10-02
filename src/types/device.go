package types

type SDInformationFromBackend struct {
	DeviceID  uint32              `json:"deviceId"`            // Device ID
	DeviceUID string              `json:"deviceUID"`           // Device ID
	Snapshot  SDParameterSnapshot `json:"sdParameterSnapshot"` // SD parameter snapshot
	Command   SDCommand           `json:"sdCommand"`           // SD command
}

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

type SDCommand struct {
	CommandID         uint32 `json:"commandId"`         // Command ID
	CommandDenotation string `json:"commandDenotation"` // Command
	Payload           string `json:"payload"`           // Payload
}

type SDCommandInvocation struct {
	InstanceID        uint32 `json:"instanceId"`        // Instance ID
	InstanceUID       string `json:"instanceUID"`       // Instance ID
	CommandID         uint32 `json:"commandId"`         // Command ID
	CommandDenotation string `json:"commandDenotation"` // Command
	Payload           string `json:"payload,omitempty"` // Payload
	InvocationTime    string `json:"invocationTime"`    // Invocation time
}

type TypeValue struct {
	Type  string
	Value any
}

type CommandPayload struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Values any    `json:"possibleValues"`
}
