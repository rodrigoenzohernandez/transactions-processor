package main

import (
	"github.com/rodrigoenzohernandez/transactions-processor/internal/repository"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/services"
	ssm_services "github.com/rodrigoenzohernandez/transactions-processor/internal/services/ssm"
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

	db, _ := repository.Connect()

	transactionsRepo := repository.NewTransactionRepo(db)

	e := transactionsRepo.InsertMany(records)
	if e != nil {
		return
	}

	report := utils.GenerateReport(records)

	emailContent, err := utils.GenerateEmailContent(report)

	if err != nil {
		return
	}

	notificationEmail, _ := ssm_services.GetSSMParameter("/smtp/provider/sender")

	services.SendEmail(notificationEmail, "Summary of your transactions", emailContent)

}
