package utils

import (
	"encoding/csv"
	"io"
	"log"
)

// Reads a buffer, removes the header and returns a slice with the records.
// e.g [[1 1/1 -150] [11 5/13 +10.3]] - Each records is a slice itself
func GetRecordsFromBuffer(buffer io.Reader) [][]string {

	r := csv.NewReader(buffer)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Remove header
	if len(records) > 0 {
		records = records[1:]
	}

	return records
}
