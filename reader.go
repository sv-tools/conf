package conf

import (
	"context"
)

// Reader is an interface for the configuration readers
type Reader interface {
	// Read reads the data and returns a raw data
	Read(ctx context.Context) (interface{}, error)
	// Returns a prefix to be used for all keys of the values provided by the reader
	Prefix() string
}
