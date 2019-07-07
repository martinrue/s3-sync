package sync

import (
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
)

// CreateCASKey creates a content-addressable key for the given file.
func CreateCASKey(filename string) (string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
