package utils

import (
	"bytes"
	"log"
	"text/template"

	"github.com/rodrigoenzohernandez/transactions-processor/internal/types"
)

// Receives a report and creates a string email content using the html template called balance
func GenerateEmailContent(report types.Report) string {
	tmpl, err := template.ParseFiles("internal/templates/balance.html")
	if err != nil {
		log.Fatal("Error parsing template:", err)
	}

	var template bytes.Buffer
	err = tmpl.Execute(&template, report)
	if err != nil {
		log.Fatal("Error executing template:", err)
	}

	return template.String()
}
