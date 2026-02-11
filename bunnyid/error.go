package bunnyid

import "fmt"

// InvalidIDError is returned when trying to decode an invalid ID.
type InvalidIDError struct {
	Value        []byte
	Type         string
	DetectedType string
}

// Error implements the error interface
func (e *InvalidIDError) Error() string {
	if e.DetectedType != "" {
		return fmt.Sprintf("Invalid %s ID '%s': You're passing a %s ID, but we need a %s ID here.", e.Type, e.Value, e.DetectedType, e.Type)
	}
	return fmt.Sprintf("Invalid %s ID '%s'", e.Type, e.Value)
}
