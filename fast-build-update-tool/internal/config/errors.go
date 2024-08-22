package config

import "fmt"

// MissingArgumentError is used when the application is called with a missing required argument
type MissingArgumentError struct {
	ArgumentName string
}

func (m *MissingArgumentError) Error() string {
	return fmt.Sprintf("missing required argument %s", m.ArgumentName)
}

// InvalidArgumentError is used when the application, or a function is called with an invalid argument
type InvalidArgumentError struct {
	ArgumentName      string
	ValidationMessage string
}

func (m *InvalidArgumentError) Error() string {
	return fmt.Sprintf("argument %s was invalid: %s", m.ArgumentName, m.ValidationMessage)
}

func missingFileError(arg string) error {
	return &InvalidArgumentError{ArgumentName: arg, ValidationMessage: "could not find file"}
}

func missingArgumentError(arg string) error {
	return &MissingArgumentError{ArgumentName: arg}
}

func invalidArgumentError(arg, validation string) error {
	return &InvalidArgumentError{ArgumentName: arg, ValidationMessage: validation}
}

// UnknownOperatingSystemError is an error used when an operating system is not known, and cannot be handled by this application
func UnknownOperatingSystemError(os string) error {
	return &InvalidArgumentError{ArgumentName: "operatingSystem", ValidationMessage: fmt.Sprintf("unknown operating system %s", os)}
}
