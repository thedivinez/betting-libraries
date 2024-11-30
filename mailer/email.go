package mailer

import (
	"bytes"
	"crypto/tls"
	"embed"
	"html/template"
	"log"
	"strconv"

	"github.com/pkg/errors"
	"gopkg.in/gomail.v2"
)

//go:embed templates
var emailTemplate embed.FS

type EmailMailer struct {
	port   int
	server string
	email  string
	pwd    string
	sender string
}

type MailerConfigs struct {
	MailerPort   string `json:"MAILER_PORT"`
	MailerPwd    string `json:"MAILER_PWD"`
	MailerEmail  string `json:"MAILER_EMAIL"`
	MailerSender string `json:"MAILER_SENDER"`
	MailerServer string `json:"MAILER_SERVER"`
}

func NewEmailMailer(mailer MailerConfigs) *EmailMailer {
	port, err := strconv.Atoi(mailer.MailerPort)
	if err != nil {
		log.Fatal(err)
	}
	return &EmailMailer{email: mailer.MailerEmail, server: mailer.MailerServer, sender: mailer.MailerSender, pwd: mailer.MailerPwd, port: port}
}

/****************************************
*	READ EMAIL BODY FRO HTML FILE METHOD	*
****************************************/
func (mailer *EmailMailer) ReadTemplate(file string, data interface{}) (string, error) {
	t, err := template.ParseFS(emailTemplate, "templates/"+file+".html")
	if err != nil {
		log.Fatal(err)
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		return "", errors.WithStack(err)
	}
	return buf.String(), nil
}

/****************************
*	SEND EMAIL TO USER METHOD	*
****************************/
func (mailer *EmailMailer) Send(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetBody("text/html", body)
	msg.SetHeader("Subject", subject)
	msg.SetHeader("From", mailer.sender+" <"+mailer.email+">")
	n := gomail.NewDialer(mailer.server, mailer.port, mailer.email, mailer.pwd)
	n.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := n.DialAndSend(msg); err != nil {
		return errors.WithStack(err)
	}
	return nil
}