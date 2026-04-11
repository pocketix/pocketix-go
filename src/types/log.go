package types

type LogInvocation struct {
	Content        string `json:"payload"`
	InvocationTime string `json:"invocationTime"` // Invocation time
}
