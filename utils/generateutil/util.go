package generateutil

import (
	"encoding/hex"
	"github.com/satori/go.uuid"
)

// GenerateID uses uuid v5 for generating
// the unique id. It accepts "name" which
// can be very useful for generating unique
// id for a database.
//
// Used hex in order to shorten the generated
// id.
func GenerateID(name string) string {
	gen := uuid.NewV5(uuid.NewV4(), name)
	return hex.EncodeToString(gen.Bytes())
}
