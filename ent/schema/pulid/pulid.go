package pulid

import (
	"crypto/rand"
	"database/sql/driver"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/oklog/ulid/v2"
)

// The default entropy source.
var defaultEntropySource *ulid.MonotonicEntropy

func init() {
	// Seed the default entropy source.
	// TODO: To improve testability, this package should allow control of entropy sources and the time.Now implementation.
	defaultEntropySource = ulid.Monotonic(rand.Reader, 0)
}

const separator = '_'

// ErrIncorrectIDFormat is returned when the ID is not in the correct format.
var ErrIncorrectIDFormat = fmt.Errorf("pulid: incorrect id format")

// NewULID returns a new ULID for time.Now() using the default entropy source.
var NewULID = func() ulid.ULID {
	return ulid.MustNew(ulid.Timestamp(time.Now()), defaultEntropySource)
}

// ID implements a PULID - a prefixed ULID.
type ID string

// ParsePrefix return the prefix from a Prefixed-ULID.
func ParsePrefix(id ID) (string, error) {
	idx := strings.IndexRune(string(id), separator)
	if idx == -1 {
		return "", ErrIncorrectIDFormat
	}

	return string(id[:idx]), nil
}

// MustNew returns a new PULID for time.Now() given a prefix. This uses the default entropy source.
func MustNew(prefix string) ID {
	return ID(fmt.Sprintf("%s%c%s", prefix, separator, NewULID()))
}

// UnmarshalGQL implements the graphql.Unmarshaler interface
func (u *ID) UnmarshalGQL(v interface{}) error {
	return u.Scan(v)
}

// MarshalGQL implements the graphql.Marshaler interface
func (u ID) MarshalGQL(w io.Writer) {
	_, _ = io.WriteString(w, strconv.Quote(string(u)))
}

// Scan implements the Scanner interface.
func (u *ID) Scan(src interface{}) error {
	if src == nil {
		return fmt.Errorf("pulid: expected a value")
	}
	switch s := src.(type) {
	case string:
		*u = ID(s)
	case []byte:
		*u = ID(s)
	default:
		return fmt.Errorf("pulid: expected a string, got %t", src)
	}

	return nil
}

// Value implements the driver Valuer interface.
func (u ID) Value() (driver.Value, error) {
	return string(u), nil
}
