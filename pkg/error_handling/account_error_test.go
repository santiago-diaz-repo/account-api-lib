package error_handling

import (
	"reflect"
	"testing"
)

func getAccountError() AccountError {
	return AccountError{
		operation: "test",
		code:      1,
		message:   "test error",
	}
}

func TestAccountError_ShouldReturnANewAccountError(t *testing.T) {

	want := "*error_handling.AccountError"
	subject := NewAccountError("test", 1, "test")
	got := reflect.TypeOf(subject).String()

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}
}

func TestAccountError_ShouldReturnErrorMessage(t *testing.T) {

	want := "test: 1 - test error"
	subject := getAccountError()
	got := subject.Error()

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}
}

func TestAccountError_ShouldReturnCode(t *testing.T) {

	want := 1
	subject := getAccountError()

	got := subject.GetCode()

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}
}

func TestAccountError_ShouldReturnMessage(t *testing.T) {

	want := "test error"
	subject := getAccountError()

	got := subject.GetMessage()

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}
}

func TestAccountError_ShouldReturnOperation(t *testing.T) {

	want := "test"
	subject := getAccountError()

	got := subject.GetOperation()

	if got != want {
		t.Errorf("wanted: %v\n got: %v", want, got)
	}
}
