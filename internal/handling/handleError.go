package handling

import (
	"crypto/tls"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"log"
)

func HandleError(err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", "zakimamka@gmail.com")
	m.SetHeader("To", "lerest.go@gmail.com")
	m.SetHeader("Subject", err.Error())
	m.SetBody("text/plain", "This is Gomail test body")
	d := gomail.NewDialer("smtp.gmail.com", 587, "zakimamka@gmail.com", "Jmhy^m$SK#2dL!#N")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		log.Println(err)
	}
}
