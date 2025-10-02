package utils

type PayloadTypeMismatchError struct {
	CommandDenotation string
	ExpectedType      string
	ActualType        string
}

func (e *PayloadTypeMismatchError) Error() string {
	return "Payload type mismatch for command " + e.CommandDenotation + ": expected " + e.ExpectedType + ", got " + e.ActualType
}
