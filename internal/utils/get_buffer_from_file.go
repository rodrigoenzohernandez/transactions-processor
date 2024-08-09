package utils

import (
	"fmt"
	"os"
)

// Reads a file from the given path and returns the buffer
// e.g buffer(0xc000570940)
func GetBufferFromFile(path string) (*os.File, error) {
	buffer, err := os.Open(path)

	if err != nil {
		log.Error(fmt.Sprintf("Error while reading the file %v", err))
		return nil, err
	}

	log.Info(fmt.Sprintf("Buffer %v", buffer))

	return buffer, nil

}
