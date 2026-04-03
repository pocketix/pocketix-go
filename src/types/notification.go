package types

type NotificationInvocation struct {
	EndpointType   string `json:"endpointType"`
	AddresseeID    uint32 `json:"addresseeId"`
	Content        string `json:"payload"`
	InvocationTime string `json:"invocationTime"` // Invocation time
}
