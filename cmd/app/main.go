package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	s "strings"
	"time"
)

type MonthBalance struct {
	Count     int
	AvgCredit float64
	AvgDebit  float64
}

type Result struct {
	TotalBalance        float64
	TransactionsByMonth map[string]MonthBalance
}

func main() {

	// Read from CSV file and remove the header
	file, err := os.Open("files/txns.csv")

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

	// Calculate summary information
	var creditBalance, debitBalance float64

	result := Result{
		TotalBalance:        0,
		TransactionsByMonth: make(map[string]MonthBalance),
	}

	for _, record := range records {
		date := record[1]
		amount, _ := strconv.ParseFloat(record[2][1:], 64)
		sign := string(record[2][0])
		monthNumber, _ := strconv.Atoi(s.Split(date, "/")[0])
		monthName := time.Month(monthNumber).String()
		monthBalance := result.TransactionsByMonth[monthName]
		monthBalance.Count++

		if sign == "+" {
			creditBalance += amount
			monthBalance.AvgCredit += amount // accumulate the credit amount, then it will be divided by count outside this loop
		} else {
			debitBalance += amount
			monthBalance.AvgDebit += amount // accumulate the debit amount, then it will be divided by count outside this loop
		}

		result.TransactionsByMonth[monthName] = monthBalance

	}
	result.TotalBalance = creditBalance - debitBalance

	for month, balance := range result.TransactionsByMonth {
		if balance.Count > 0 {

			balance.AvgCredit = balance.AvgCredit / float64(balance.Count)
			balance.AvgDebit = balance.AvgDebit / float64(balance.Count)
			result.TransactionsByMonth[month] = balance
		}
	}

	// Prepare the result in the desired email format
	tmpl, err := template.ParseFiles("internal/templates/balance.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	var template bytes.Buffer
	err = tmpl.Execute(&template, result)
	if err != nil {
		log.Fatal("Error executing template:", err)
	}

	emailContent := template.String()

	fmt.Println(emailContent)

}
