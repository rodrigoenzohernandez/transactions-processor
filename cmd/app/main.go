package main

import (
	"fmt"

	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils"
)

func main() {

	records := utils.GetRecordsFromCSV("files/txns.csv")

	report := utils.GenerateReport(records)

	emailContent := utils.GenerateEmailContent(report)

	fmt.Println(emailContent)

}
