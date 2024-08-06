package utils

import (
	"crypto/tls"
	"strconv"

	"github.com/gofiber/fiber/v2/log"
	gomail "gopkg.in/mail.v2"
)

func SendEmail(subject string, body string, email string) error {
	log.Info("Sending email to: ", email)
	m := gomail.NewMessage()
	m.SetHeader("From", Get("MAIL_SMTP"))
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)

	m.SetBody("text/html", body)

	port, err := strconv.Atoi(Get("SMTP_PORT"))

	if err != nil {
		return err
	}

	d := gomail.NewDialer(Get("SMTP_HOST"), port, Get("MAIL_SMTP"), Get("MAIL_PASSWORD"))

	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
