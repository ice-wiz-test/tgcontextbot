package handling

import (
	"crypto/tls"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"log"
)

func HandleError(err error) {
	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", "zakimamka@gmail.com")

	// Set E-Mail receivers
	m.SetHeader("To", "lerest.go@gmail.com")

	// Set E-Mail subject
	m.SetHeader("Subject", err.Error())

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/plain", "This is Gomail test body")

	// Settings for SMTP server
	d := gomail.NewDialer("smtp.gmail.com", 587, "zakimamka@gmail.com", "Jmhy^m$SK#2dL!#N")

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		log.Println(err)
	}
}
