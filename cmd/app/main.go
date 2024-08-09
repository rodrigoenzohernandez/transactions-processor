package main

import (
	"fmt"

	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils"
)

func main() {

	buffer := utils.GetBufferFromFile("files/txns.csv")

	defer buffer.Close()

	records := utils.GetRecordsFromBuffer(buffer)

	report := utils.GenerateReport(records)

	emailContent := utils.GenerateEmailContent(report)

	fmt.Println(emailContent)

}
