package boil

import (
	"time"
)

var (
	// timestampLocation is the timezone used for the
	// automated setting of created_at/updated_at columns
	timestampLocation = time.UTC
)

// SetLocation sets the global timestamp Location.
// This is the timezone used by the generated package for the
// automated setting of created_at and updated_at columns.
// If the package was generated with the --no-auto-timestamps flag
// then this function has no effect.
func SetLocation(loc *time.Location) {
	timestampLocation = loc
}

// GetLocation retrieves the global timestamp Location.
// This is the timezone used by the generated package for the
// automated setting of created_at and updated_at columns
// if the package was not generated with the --no-auto-timestamps flag.
func GetLocation() *time.Location {
	return timestampLocation
}
