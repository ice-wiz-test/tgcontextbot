package handling

import (
	"crypto/tls"
	gmail "gopkg.in/mail.v2"
	"log"
)

func HandleError(err error) {
	m := gmail.NewMessage()
	m.SetHeader("From", "zakimamka@gmail.com")
	m.SetHeader("To", "lerest.go@gmail.com")
	m.SetHeader("Subject", "ERROR")
	m.SetBody("text/plain", err.Error())
	d := gmail.NewDialer("smtp.gmail.com", 587, "zakimamka@gmail.com", "Jmhy^m$SK#2dL!#N")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
	m.SetHeader("To", "prokhoridze16@gmail.com")
	d = gmail.NewDialer("smtp.gmail.com", 587, "zakimamka@gmail.com", "Jmhy^m$SK#2dL!#N")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
}
