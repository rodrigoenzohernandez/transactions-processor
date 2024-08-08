package utils

import (
	"strconv"
	s "strings"
	"time"

	"github.com/rodrigoenzohernandez/transactions-processor/internal/types"
)

/*
Receives records and calculates:
- Total balance
- Transactions count by month
- Average credit by month
- Average debit by month
Then returns it as Report
*/
func GenerateReport(records [][]string) types.Report {
	var creditBalance, debitBalance float64

	report := types.Report{
		TotalBalance:        0,
		TransactionsByMonth: make(types.TransactionsByMonth),
	}

	for _, record := range records {
		date := record[1]
		amount, _ := strconv.ParseFloat(record[2][1:], 64)
		sign := string(record[2][0])
		monthNumber, _ := strconv.Atoi(s.Split(date, "/")[0])
		monthName := time.Month(monthNumber).String()
		monthBalance := report.TransactionsByMonth[monthNumber]
		monthBalance.Count++
		monthBalance.Name = monthName

		if sign == "+" {
			creditBalance += amount
			monthBalance.AvgCredit += amount
		} else {
			debitBalance += amount
			monthBalance.AvgDebit += amount
		}

		report.TransactionsByMonth[monthNumber] = monthBalance

	}

	report.TotalBalance = ToFixed(creditBalance - debitBalance)

	for month, balance := range report.TransactionsByMonth {
		if balance.Count > 0 {

			balance.AvgCredit = ToFixed(balance.AvgCredit / float64(balance.Count))
			balance.AvgDebit = ToFixed(balance.AvgDebit / float64(balance.Count))
			report.TransactionsByMonth[month] = balance
		}
	}

	return report
}
