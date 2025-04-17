package main

import (
	"bytes"
	_ "embed"
	"fmt"
	htmlTemplate "html/template"
	"log/slog"
	"net"
	"net/smtp"
	"net/textproto"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	textTemplate "text/template"
	"time"

	"github.com/drone/drone-go/plugin/webhook"
	"github.com/jordan-wright/email"
)

const emailSenderShutdownTimeout = 60 * time.Second

var (
	//go:embed email.html
	htmlTemplStr string
	//go:embed email.txt
	textTemplStr string

	htmlTempl = htmlTemplate.Must(htmlTemplate.New("html").Parse(htmlTemplStr)) //nolint:gochecknoglobals
	textTempl = textTemplate.Must(textTemplate.New("text").Parse(textTemplStr)) //nolint:gochecknoglobals
)

type EmailSender struct {
	host     string
	addr     string
	username string
	password string
	from     string
	cc       []string
	bcc      []string

	closed atomic.Bool
	wg     sync.WaitGroup
}

func NewEmailSender(cfg Config) *EmailSender {
	return &EmailSender{
		host:     cfg.EmailSMTPHost,
		addr:     net.JoinHostPort(cfg.EmailSMTPHost, strconv.Itoa(int(cfg.EmailSMTPPort))),
		username: cfg.EmailSMTPUsername,
		password: cfg.EmailSMTPPassword,
		from:     cfg.EmailFrom,
		cc:       cfg.EmailCC,
		bcc:      cfg.EmailBCC,

		closed: atomic.Bool{},
		wg:     sync.WaitGroup{},
	}
}

func (s *EmailSender) SendAsync(req *webhook.Request) {
	if s.closed.Load() {
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		_ = s.Send(req)
	}()
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
		From:            s.from,
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
	if err := htmlTempl.Execute(&html, &data); err != nil {
		slog.Error("email sender cannot execute HTML template", "build_number", req.Build.Number, "error", err)
		return fmt.Errorf("email sender cannot execute HTML template: %w", err)
	}

	var text bytes.Buffer
	if err := textTempl.Execute(&text, &data); err != nil {
		slog.Error("email sender cannot execute text template", "build_number", req.Build.Number, "error", err)
		return fmt.Errorf("email sender cannot execute text template: %w", err)
	}

	//nolint:exhaustruct
	emailMsg := &email.Email{
		From:    data.From,
		To:      []string{data.To},
		Cc:      s.cc,
		Bcc:     s.bcc,
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
		slog.Error("email sender failed to send message", "build_number", req.Build.Number, "to", data.To, "error", err)
		return fmt.Errorf("email sender failed to send message: %w", err)
	}
	slog.Info("email sender successfully sent message", "build_number", req.Build.Number, "to", data.To)
	return nil
}

func (s *EmailSender) Shutdown() {
	if s.closed.Swap(true) {
		return
	}
	slog.Info("email sender initiating shutdown")

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		slog.Info("email sender completed shutdown")
	case <-time.After(emailSenderShutdownTimeout):
		slog.Error("email sender shutdown timed out")
	}
}
