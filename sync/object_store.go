package sync

import "io"

// ObjectStore describes an object capable of storing and retrieving data for a given key.
type ObjectStore interface {
	Get(key string) ([]byte, error)
	Put(key string, contentType string, data io.ReadSeeker) error
}
