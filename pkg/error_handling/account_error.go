package error_handling

import "fmt"

type AccountError struct {
	operation string
	code      int
	message   string
}

func NewAccountError(operation string, code int, message string) error {
	return &AccountError{
		operation: operation,
		code:      code,
		message:   message,
	}
}

func (ce *AccountError) Error() string {
	return fmt.Sprintf("%s: %d - %s", ce.operation, ce.code, ce.message)
}

func (ce *AccountError) GetOperation() string {
	return ce.operation
}

func (ce *AccountError) GetCode() int {
	return ce.code
}

func (ce *AccountError) GetMessage() string {
	return ce.message
}
