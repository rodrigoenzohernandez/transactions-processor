package utils

import (
	"encoding/csv"
	"fmt"
	"io"
)

// Reads a buffer, removes the header and returns a slice with the records.
// e.g [[1 1/1 -150] [11 5/13 +10.3]] - Each records is a slice itself
func GetRecordsFromBuffer(buffer io.Reader) ([][]string, error) {

	r := csv.NewReader(buffer)

	records, err := r.ReadAll()
	if err != nil {
		log.Error(fmt.Sprintf("Error getting records from buffer %v", err))
		return nil, err

	}

	if len(records) == 0 {
		err := fmt.Errorf("the file uploaded is empty")
		log.Error(err.Error())
		return nil, err
	}

	// Remove header
	if len(records) > 0 {
		records = records[1:]
	}

	log.Info(fmt.Sprintf("Records: %s", records))

	return records, nil
}
