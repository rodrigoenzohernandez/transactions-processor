package utils

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/rodrigoenzohernandez/transactions-processor/internal/types"
)

// Receives a report and creates a string email content using the html template called balance
func GenerateEmailContent(report types.Report) (string, error) {
	tmpl, err := template.ParseFiles("internal/templates/balance.html")
	if err != nil {
		log.Error(fmt.Sprintf("Error parsing template %v", err))
		return "", err
	}

	var template bytes.Buffer
	err = tmpl.Execute(&template, report)
	if err != nil {
		log.Error(fmt.Sprintf("Error executing template %v", err))
		return "", err
	}

	return template.String(), nil
}
