package main

import (
	"bytes"
	_ "embed"
	"fmt"
	htmlTemplate "html/template"
	"log"
	"net/smtp"
	"net/textproto"
	"strings"
	textTemplate "text/template"

	"github.com/drone/drone-go/plugin/webhook"
	"github.com/jordan-wright/email"
)

//go:embed email.html
var htmlTemplateStr string

//go:embed email.txt
var textTemplateStr string

type EmailSender struct {
	host     string
	port     int
	username string
	password string
	from     string
	html     *htmlTemplate.Template
	text     *textTemplate.Template
}

func NewEmailSender(settings Settings) *EmailSender {
	return &EmailSender{
		host:     settings.EmailSmtpHost,
		port:     settings.EmailSmtpPort,
		username: settings.EmailSmtpUsername,
		password: settings.EmailSmtpPassword,
		from:     settings.EmailFrom,
		html:     htmlTemplate.Must(htmlTemplate.New("html").Parse(htmlTemplateStr)),
		text:     textTemplate.Must(textTemplate.New("text").Parse(textTemplateStr)),
	}
}

func (s *EmailSender) Send(req *webhook.Request) error {
	var author string
	if req.Build.AuthorName != "" {
		author = req.Build.AuthorName
	} else {
		author = req.Build.Author
	}

	data := struct {
		Subject         string
		From            string
		To              string
		Header          string
		Repository      string
		Reference       string
		CommitHash      string
		CommitMessage   string
		AuthorAvatar    string
		AuthorName      string
		DroneBuildLink  string
		DroneServerHost string
		DroneServerLink string
	}{
		Subject:         fmt.Sprintf("[%s] Failed build #%d for %s (%s)", req.Repo.Slug, req.Build.Number, req.Build.Ref, req.Build.After[:8]),
		From:            fmt.Sprintf("%s <%s>", "Drone", s.from),
		To:              fmt.Sprintf("%s <%s>", author, req.Build.AuthorEmail),
		Header:          fmt.Sprintf("Build #%d has failed", req.Build.Number),
		Repository:      req.Repo.Slug,
		Reference:       req.Build.Ref,
		CommitHash:      req.Build.After[:8],
		CommitMessage:   strings.TrimSpace(req.Build.Message),
		AuthorAvatar:    req.Build.AuthorAvatar,
		AuthorName:      author,
		DroneBuildLink:  fmt.Sprintf("%s/%s/%d", req.System.Link, req.Repo.Slug, req.Build.Number),
		DroneServerHost: req.System.Host,
		DroneServerLink: req.System.Link,
	}

	var html bytes.Buffer
	err := s.html.Execute(&html, &data)
	if err != nil {
		log.Println("email: cannot execute html template:", err)
		return err
	}

	var text bytes.Buffer
	err = s.text.Execute(&text, &data)
	if err != nil {
		log.Println("email: cannot execute text template:", err)
		return err
	}

	err = (&(email.Email{
		From:    data.From,
		To:      []string{data.To},
		Subject: data.Subject,
		HTML:    html.Bytes(),
		Text:    text.Bytes(),
		Headers: textproto.MIMEHeader{},
	})).Send(fmt.Sprintf("%s:%d", s.host, s.port), smtp.PlainAuth("", s.username, s.password, s.host))
	if err != nil {
		log.Println("email: cannot send mail:", err)
		return err
	}
	return nil
}
