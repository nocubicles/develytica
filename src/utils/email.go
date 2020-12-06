package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func sendEmail(to []string, msg []byte) error {

	emailHost := os.Getenv("emailHost")
	smtpPort := os.Getenv("smtpPort")
	fromEmail := os.Getenv("fromEmail")
	smtpUser := os.Getenv("smtpUser")
	smtpPsw := os.Getenv("smtpPsw")

	auth := smtp.PlainAuth(
		"",
		smtpUser,
		smtpPsw,
		emailHost,
	)

	err := smtp.SendMail(emailHost+":"+smtpPort, auth, fromEmail, to, msg)

	if err != nil {
		return err
	}
	return nil
}

func SendSignInEmail(to string, signInToken string) error {
	server := os.Getenv("SERVER")
	fromEmail := os.Getenv("fromEmail")
	fromMime := fmt.Sprintf("From: %s\n", fromEmail)
	loginLink := fmt.Sprintf("<a href='%s/logisisse?signintoken=%s'>lingile</a>", server, signInToken)
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	msg := fmt.Sprintf(
		"To: %s\r\n"+
			fromMime+
			"Subject: Teie sisse logimise link skillbase.io\n"+
			mime+
			"<html><body>"+
			"<p>Palun vajutage %s, et logida sisse.</p>"+
			"<p>skillbase.io</p>"+
			"</body></html>",
		to, loginLink,
	)

	err := sendEmail([]string{to}, []byte(msg))

	if err != nil {
		return err
	}
	return nil
}
