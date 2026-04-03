package types

type NotificationInvocation struct {
	EndpointType   string `json:"endpointType"`
	AddresseeID    string `json:"addresseeId"`
	Message        string `json:"payload"`
	InvocationTime string `json:"invocationTime"` // Invocation time
}
