package main

import (
	"fmt"

	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils"
)

func main() {

	buffer, err := utils.GetBufferFromFile("files/txns.csv")
	if err != nil {
		return
	}
	defer buffer.Close()

	records, err := utils.GetRecordsFromBuffer(buffer)
	if err != nil {
		return
	}

	report := utils.GenerateReport(records)

	emailContent, err := utils.GenerateEmailContent(report)

	if err != nil {
		return
	}

	fmt.Println(emailContent)

}
