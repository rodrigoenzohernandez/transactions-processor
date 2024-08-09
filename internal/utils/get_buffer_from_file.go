package utils

import (
	"log"
	"os"
)

// Reads a file from the given path and returns the buffer
// e.g buffer(0xc000570940)
func GetBufferFromFile(path string) *os.File {
	buffer, err := os.Open(path)

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	log.Println("Downloaded object", buffer)

	return buffer

}
