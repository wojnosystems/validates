package ifaces

import (
	"golang.org/x/text/message"
)

type ValidateError interface {
	// Creates an internationalized string representing the error. Uses message.Sprintf so you can customize the messages for which region your user is currently in
	// To use the default, specify
	ErrorI18n(*message.Printer) string

	// IsEqual returns true if the two ValidateErrors are the same type and contain the same data, false if not
	IsEqual(ValidateError) bool
}
