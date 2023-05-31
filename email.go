package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"net/textproto"
	"strings"

	"github.com/drone/drone-go/plugin/webhook"
	"github.com/jordan-wright/email"
)

//go:embed email.html
var bodyTemplate string

type EmailSender struct {
	host     string
	port     int
	username string
	password string
	from     string
	body     *template.Template
}

func NewEmailSender(settings Settings) *EmailSender {
	return &EmailSender{
		host:     settings.EmailSmtpHost,
		port:     settings.EmailSmtpPort,
		username: settings.EmailSmtpUsername,
		password: settings.EmailSmtpPassword,
		from:     settings.EmailFrom,
		body:     template.Must(template.New("body").Parse(bodyTemplate)),
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
		Subject:         fmt.Sprintf("[%s] Failed build for %s (%s)", req.Repo.Slug, req.Build.Ref, req.Build.After[:8]),
		From:            fmt.Sprintf("%s <%s>", "Drone", s.from),
		To:              fmt.Sprintf("%s <%s>", author, req.Build.AuthorEmail),
		Header:          fmt.Sprintf("Build #%d has failed", req.Build.ID),
		Repository:      req.Repo.Slug,
		Reference:       req.Build.Ref,
		CommitHash:      req.Build.After[:8],
		CommitMessage:   strings.TrimSpace(req.Build.Message),
		AuthorAvatar:    req.Build.AuthorAvatar,
		AuthorName:      author,
		DroneBuildLink:  fmt.Sprintf("%s/%s/%d", req.System.Link, req.Repo.Slug, req.Build.ID),
		DroneServerHost: req.System.Host,
		DroneServerLink: req.System.Link,
	}
	var msg bytes.Buffer
	err := s.body.Execute(&msg, &data)
	if err != nil {
		log.Println("email: cannot execute body template:", err)
		return err
	}

	err = (&(email.Email{
		From:    data.From,
		To:      []string{data.To},
		Subject: data.Subject,
		Text:    []byte(data.Header),
		HTML:    msg.Bytes(),
		Headers: textproto.MIMEHeader{},
	})).Send(fmt.Sprintf("%s:%d", s.host, s.port), smtp.PlainAuth("", s.username, s.password, s.host))
	if err != nil {
		log.Println("email: cannot send mail:", err)
		return err
	}
	return nil
}
