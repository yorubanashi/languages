package db

import (
	"os"

	"gopkg.in/yaml.v3"
)

// FetchYAML retrieves the file specified and serializes it into the passed object.
//
// This is just a fun attempt at abstracting a DB that uses the local filesystem.
func FetchYAML(filepath string, obj interface{}) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(data, obj)
}
