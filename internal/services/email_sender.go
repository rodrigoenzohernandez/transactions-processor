package services

import (
	"fmt"

	smpt_provider "github.com/mailjet/mailjet-apiv3-go/v4"
	ssm_services "github.com/rodrigoenzohernandez/transactions-processor/internal/services/ssm"
	"github.com/rodrigoenzohernandez/transactions-processor/internal/utils/logger"
)

var log = logger.GetLogger("email_sender")

func SendEmail(to string, subject string, emailContent string) {

	publicKey, _ := ssm_services.GetSSMParameter("/smtp/provider/public_key")
	privateKey, _ := ssm_services.GetSSMParameter("/smtp/provider/private_key")
	emailSender, _ := ssm_services.GetSSMParameter("/smtp/provider/sender")

	client := smpt_provider.NewMailjetClient(publicKey, privateKey)
	messagesInfo := []smpt_provider.InfoMessagesV31{
		{
			From: &smpt_provider.RecipientV31{
				Email: emailSender,
				Name:  "transactions-processor",
			},
			To: &smpt_provider.RecipientsV31{
				smpt_provider.RecipientV31{
					Email: to,
				},
			},
			Subject:  subject,
			HTMLPart: emailContent,
		},
	}
	messages := smpt_provider.MessagesV31{Info: messagesInfo}
	res, err := client.SendMailV31(&messages)

	if err != nil {
		log.Error(fmt.Sprintf("Error sending the email %v", err))
		return
	}
	log.Debug(fmt.Sprintf("Email sent %v", res))

}
