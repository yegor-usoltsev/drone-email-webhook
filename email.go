package main

import (
	"bytes"
	_ "embed"
	"fmt"
	htmlTemplate "html/template"
	"net/smtp"
	"strings"
	textTemplate "text/template"

	"github.com/drone/drone-go/plugin/webhook"
)

const headersTemplate = `Subject: {{.Subject}}
From: {{.From}}
To: {{.To}}
Mime-Version: 1.0;
Content-Type: text/html; charset=UTF-8;

`

//go:embed email.html
var bodyTemplate string

type EmailSender struct {
	host     string
	port     int
	username string
	password string
	from     string
	headers  *textTemplate.Template
	body     *htmlTemplate.Template
}

func NewEmailSender(settings Settings) *EmailSender {
	return &EmailSender{
		host:     settings.EmailSmtpHost,
		port:     settings.EmailSmtpPort,
		username: settings.EmailSmtpUsername,
		password: settings.EmailSmtpPassword,
		from:     settings.EmailFrom,
		headers:  textTemplate.Must(textTemplate.New("headers").Parse(headersTemplate)),
		body:     htmlTemplate.Must(htmlTemplate.New("body").Parse(bodyTemplate)),
	}
}

func (s *EmailSender) Send(req *webhook.Request) error {
	// TODO: logging
	addr := fmt.Sprintf("%s:%d", s.host, s.port)
	auth := smtp.PlainAuth("", s.username, s.password, s.host)

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
		From:            fmt.Sprintf("\"%s\" <%s>", "Drone", s.from),
		To:              fmt.Sprintf("\"%s\" <%s>", author, req.Build.AuthorEmail),
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
	err := s.headers.Execute(&msg, &data)
	if err != nil {
		return err
	}
	err = s.body.Execute(&msg, &data)
	if err != nil {
		return err
	}

	return smtp.SendMail(addr, auth, s.from, []string{req.Build.AuthorEmail}, msg.Bytes())
}
