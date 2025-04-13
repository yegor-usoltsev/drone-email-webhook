package main

import (
	"bytes"
	_ "embed"
	"fmt"
	htmlTemplate "html/template"
	"log"
	"net"
	"net/smtp"
	"net/textproto"
	"strconv"
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
	addr     string
	username string
	password string
	from     string
	html     *htmlTemplate.Template
	text     *textTemplate.Template
}

func NewEmailSender(settings Settings) *EmailSender {
	return &EmailSender{
		host:     settings.EmailSMTPHost,
		addr:     net.JoinHostPort(settings.EmailSMTPHost, strconv.Itoa(settings.EmailSMTPPort)),
		username: settings.EmailSMTPUsername,
		password: settings.EmailSMTPPassword,
		from:     settings.EmailFrom,
		html:     htmlTemplate.Must(htmlTemplate.New("html").Parse(htmlTemplateStr)),
		text:     textTemplate.Must(textTemplate.New("text").Parse(textTemplateStr)),
	}
}

func (s *EmailSender) Send(req *webhook.Request) error {
	author := req.Build.AuthorName
	if author == "" {
		author = req.Build.Author
	}

	commitHash := req.Build.After
	if len(commitHash) > 8 {
		commitHash = commitHash[:8]
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
		Subject:         fmt.Sprintf("[%s] Failed build #%d for %s (%s)", req.Repo.Slug, req.Build.Number, req.Build.Ref, commitHash),
		From:            fmt.Sprintf("%s <%s>", "Drone", s.from),
		To:              fmt.Sprintf("%s <%s>", author, req.Build.AuthorEmail),
		Header:          fmt.Sprintf("Build #%d has failed", req.Build.Number),
		Repository:      req.Repo.Slug,
		Reference:       req.Build.Ref,
		CommitHash:      commitHash,
		CommitMessage:   strings.TrimSpace(strings.Split(req.Build.Message, "\n")[0]),
		AuthorAvatar:    req.Build.AuthorAvatar,
		AuthorName:      author,
		DroneBuildLink:  fmt.Sprintf("%s/%s/%d", req.System.Link, req.Repo.Slug, req.Build.Number),
		DroneServerHost: req.System.Host,
		DroneServerLink: req.System.Link,
	}

	var html bytes.Buffer
	if err := s.html.Execute(&html, &data); err != nil {
		log.Printf("[ERROR] email: cannot execute HTML template for build #%d: %v", req.Build.Number, err)
		return fmt.Errorf("failed to execute HTML template: %w", err)
	}

	var text bytes.Buffer
	if err := s.text.Execute(&text, &data); err != nil {
		log.Printf("[ERROR] email: cannot execute text template for build #%d: %v", req.Build.Number, err)
		return fmt.Errorf("failed to execute text template: %w", err)
	}

	//nolint:exhaustruct
	emailMsg := &email.Email{
		From:    data.From,
		To:      []string{data.To},
		Subject: data.Subject,
		HTML:    html.Bytes(),
		Text:    text.Bytes(),
		Headers: textproto.MIMEHeader{},
	}

	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	if err := emailMsg.Send(s.addr, auth); err != nil {
		log.Printf("[ERROR] email: cannot send mail for build #%d to %s: %v", req.Build.Number, data.To, err)
		return fmt.Errorf("failed to send email: %w", err)
	}
	log.Printf("[INFO] email: successfully sent notification for build #%d to %s", req.Build.Number, data.To)
	return nil
}
