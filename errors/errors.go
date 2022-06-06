package errors

import "fmt"

var _ error = &Error{}

type Error struct {
	// Type is an exit code if required
	Type int
	// Message is the message to display to the user
	Message string
}

// Error implements error
func (e *Error) Error() string {
	return fmt.Sprintf("dbadger:%d:%s", e.Type, e.Message)
}

var (
	NoSuchKey *Error = &Error{
		Type:    100,
		Message: "Specified Key Does Not Exist",
	}
	KeyNotBlob *Error = &Error{
		Type:    101,
		Message: "Specified Key Is Not A Blob",
	}
	KeyError *Error = &Error{
		Type:    199,
		Message: "Error With Specified Key",
	}
	DestIsDirectory *Error = &Error{
		Type:    200,
		Message: "Specified Destination Is A Directory",
	}
	DestError *Error = &Error{
		Type:    299,
		Message: "Error With Destination",
	}
	CopyError *Error = &Error{
		Type:    900,
		Message: "Error With Copy Function",
	}
)
