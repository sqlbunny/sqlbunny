package bunnyid

import "fmt"

type InvalidIDError struct {
	Value        []byte
	Type         string
	DetectedType string
}

func (e *InvalidIDError) Error() string {
	if e.DetectedType != "" {
		return fmt.Sprintf("Invalid %s ID '%s': You're passing a %s ID, but we need a %s ID here.", e.Type, e.Value, e.DetectedType, e.Type)
	}
	return fmt.Sprintf("Invalid %s ID '%s'", e.Type, e.Value)
}
