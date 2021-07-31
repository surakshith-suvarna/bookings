package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/surakshith-suvarna/bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	go func() {
		for {
			m := <-app.MailChan
			sendMsg(m)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.KeepAlive = false
	server.Port = 1025
	server.SendTimeout = 10 * time.Second
	server.ConnectTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}
	email := mail.NewMSG()
	email.AddTo(m.To).SetFrom(m.From).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			errorLog.Println(err)
		}
		template := string(data)
		message := strings.Replace(template, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, message)

	}

	if email.Error != nil {
		log.Println(email.Error)
	}

	err = email.Send(client)
	if err != nil {
		errorLog.Println(err)
	} else {
		log.Println("Mail Sent")
	}

}
