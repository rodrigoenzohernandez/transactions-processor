package utils

import (
	"encoding/csv"
	"log"
	"os"
)

// Reads a csv file from the given path, removes the header and returns a slice with the records.
// e.g [[1 1/1 -150] [11 5/13 +10.3]] - Each records is a slice itself
func GetRecordsFromCSV(path string) [][]string {
	// Read from CSV file and remove the header
	file, err := os.Open(path)

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}

	defer file.Close()

	r := csv.NewReader(file)

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
