package helpers

import (
	"fmt"
	"net/smtp"
	"strings"
)

const (
	//SMTPSever const stror link email smtp
	SMTPSever = "smtp.gmail.com"
)

//Sender struct include of user password
type Sender struct {
	User     string
	Password string
}

//NewSender func to return a  new sender
func NewSender(Username, Password string) Sender {

	return Sender{Username, Password}
}

//SendMail func to send  a mail to many dest with subject and body message
func (sender Sender) SendMail(Dest []string, Subject, bodyMessage string) {

	msg := "From: " + sender.User + "\n" +
		"To: " + strings.Join(Dest, ",") + "\n" +
		"Subject: " + Subject + "\n" + bodyMessage

	err := smtp.SendMail(SMTPSever+":587",
		smtp.PlainAuth("", sender.User, sender.Password, SMTPSever),
		sender.User, Dest, []byte(msg))

	if err != nil {

		fmt.Printf("smtp error: %s\n", err)
		return
	}

	fmt.Println("Mail sent successfully!")
}
