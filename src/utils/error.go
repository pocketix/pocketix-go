package utils

import (
	"reflect"
	"strings"
)

type InterpretError interface {
	error
	GetLine() int32
	GetContext() any
}

type InternalInterpretError struct {
	Msg string
}

func (e *InternalInterpretError) Error() string {
	return "InternalInterpretError: " + e.Msg
}

func (e *InternalInterpretError) GetLine() int32 {
	return -1
}

func (e *InternalInterpretError) GetContext() any {
	return nil
}

type PayloadTypeMismatchError struct {
	CommandDenotation string
	ExpectedType      string
	ActualType        string
}

func (e *PayloadTypeMismatchError) Error() string {
	return "PayloadTypeMismatchError: Command '" + e.CommandDenotation + "' expects payload type '" + e.ExpectedType + "', but got '" + e.ActualType + "'"
}

func (e *PayloadTypeMismatchError) GetLine() int32 {
	return -1
}

func (e *PayloadTypeMismatchError) GetContext() any {
	return map[string]string{
		"CommandDenotation": e.CommandDenotation,
		"ExpectedType":      e.ExpectedType,
		"ActualType":        e.ActualType,
	}
}

type PayloadValueMissingError struct {
	CommandDenotation string
	MissingValueName  string
	PossibleValues    []string
}

func (e *PayloadValueMissingError) Error() string {
	return "PayloadValueMissingError: Command '" + e.CommandDenotation + "' is missing value '" + e.MissingValueName + "'. Possible values are: " + joinStrings(e.PossibleValues, ", ")
}

func (e *PayloadValueMissingError) GetLine() int32 {
	return -1
}

func (e *PayloadValueMissingError) GetContext() any {
	return map[string]any{
		"CommandDenotation": e.CommandDenotation,
		"MissingValueName":  e.MissingValueName,
		"PossibleValues":    e.PossibleValues,
	}
}

func joinStrings(s1 []string, s2 string) string {
	return strings.Join(s1, s2)
}

func NewErrorOf[T any](args ...any) InterpretError {
	errorType := reflect.TypeOf((*T)(nil)).Elem()
	if errorType.Kind() != reflect.Struct {
		return &InternalInterpretError{Msg: "Error during error creation, T must be a struct type. Got: " + errorType.Kind().String()}
	}

	if len(args) != errorType.NumField() {
		return &InternalInterpretError{Msg: "Error during error creation, number of arguments does not match number of fields, expected " + string(rune(errorType.NumField())) + ", got " + string(rune(len(args)))}
	}

	errorPtr := reflect.New(errorType)
	errorValue := errorPtr.Elem()
	for i := 0; i < errorType.NumField(); i++ {
		arg := reflect.ValueOf(args[i])
		field := errorValue.Field(i)
		if !arg.Type().AssignableTo(field.Type()) {
			return &InternalInterpretError{Msg: "Error during error creation: argument type does not match field type, expected " + field.Type().String() + ", got " + arg.Type().String()}
		}
		field.Set(arg)
	}

	return errorPtr.Interface().(InterpretError)
}
