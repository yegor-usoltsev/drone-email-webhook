package main

import (
	"bytes"
	"cmp"
	"embed"
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

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/webhook"
	"github.com/jordan-wright/email"
)

const emailSenderShutdownTimeout = 60 * time.Second

var (
	//go:embed email.html
	htmlTemplStr string
	//go:embed email.txt
	textTemplStr string
	//go:embed email-template/emails/static
	images embed.FS

	htmlTempl = htmlTemplate.Must(htmlTemplate.New("html").Parse(htmlTemplStr))
	textTempl = textTemplate.Must(textTemplate.New("text").Parse(textTemplStr))
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

	s.wg.Go(func() {
		_ = s.Send(req)
	})
}

func (s *EmailSender) Send(req *webhook.Request) error {
	data := newEmailData(req)
	to := fmt.Sprintf("%s <%s>", data.AuthorName, req.Build.AuthorEmail)

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

	emailMsg := &email.Email{
		From:    s.from,
		To:      []string{to},
		Cc:      s.cc,
		Bcc:     s.bcc,
		Subject: data.Subject,
		HTML:    html.Bytes(),
		Text:    text.Bytes(),
		Headers: textproto.MIMEHeader{},
	}
	if err := attachImages(emailMsg); err != nil {
		slog.Error("email sender cannot attach images", "build_number", req.Build.Number, "error", err)
		return fmt.Errorf("email sender cannot attach images: %w", err)
	}

	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	if err := emailMsg.Send(s.addr, auth); err != nil {
		slog.Error("email sender failed to send message", "build_number", req.Build.Number, "to", to, "error", err)
		return fmt.Errorf("email sender failed to send message: %w", err)
	}
	slog.Info("email sender successfully sent message", "build_number", req.Build.Number, "to", to)
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

// attachImages adds the template images as inline parts that email.html references by "cid:<file name>".
func attachImages(msg *email.Email) error {
	entries, err := images.ReadDir("email-template/emails/static")
	if err != nil {
		return fmt.Errorf("cannot read images: %w", err)
	}
	for _, entry := range entries {
		data, err := images.ReadFile("email-template/emails/static/" + entry.Name())
		if err != nil {
			return fmt.Errorf("cannot read image %s: %w", entry.Name(), err)
		}
		attachment, err := msg.Attach(bytes.NewReader(data), entry.Name(), "image/png")
		if err != nil {
			return fmt.Errorf("cannot attach image %s: %w", entry.Name(), err)
		}
		attachment.HTMLRelated = true
	}
	return nil
}

// emailData holds the values that email.html and email.txt render.
type emailData struct {
	Subject       string
	AuthorAvatar  string
	AuthorName    string
	BuildLink     string
	BuildNumber   int64
	CommitHash    string
	CommitLink    string
	CommitMessage string
	Duration      string
	FailedSteps   string
	RefName       string
	Repository    string
	ServerHost    string
	ServerLink    string
}

func newEmailData(req *webhook.Request) emailData {
	build := req.Build
	ref := refName(build)
	buildLink := fmt.Sprintf("%s/%s/%d", req.System.Link, req.Repo.Slug, build.Number)
	steps, stepLink := failedSteps(build, buildLink)
	return emailData{
		Subject:       fmt.Sprintf("[%s] Build #%d failed on %s", req.Repo.Slug, build.Number, ref),
		AuthorAvatar:  build.AuthorAvatar,
		AuthorName:    cmp.Or(build.AuthorName, build.Author),
		BuildLink:     stepLink,
		BuildNumber:   build.Number,
		CommitHash:    build.After[:min(len(build.After), 8)],
		CommitLink:    cmp.Or(build.Link, buildLink),
		CommitMessage: strings.TrimSpace(strings.Split(build.Message, "\n")[0]),
		Duration:      (time.Duration(max(build.Finished-build.Started, 0)) * time.Second).String(),
		FailedSteps:   steps,
		RefName:       ref,
		Repository:    req.Repo.Slug,
		ServerHost:    req.System.Host,
		ServerLink:    req.System.Link,
	}
}

// refName shortens branch and tag refs and shows pull requests as source → target.
func refName(build *drone.Build) string {
	if build.Event == drone.EventPullRequest && build.Source != "" && build.Target != "" {
		return build.Source + " → " + build.Target
	}
	if name, ok := strings.CutPrefix(build.Ref, "refs/heads/"); ok {
		return name
	}
	if name, ok := strings.CutPrefix(build.Ref, "refs/tags/"); ok {
		return name
	}
	return build.Ref
}

const maxFailedSteps = 3

// failedSteps names the steps that failed the build and links to the log of the first one.
func failedSteps(build *drone.Build, buildLink string) (names, link string) {
	var failed []string
	link = buildLink
	for _, stage := range build.Stages {
		for _, step := range stage.Steps {
			if step.Status != drone.StatusFailing || step.ErrIgnore {
				continue
			}
			if len(failed) == 0 {
				link = fmt.Sprintf("%s/%d/%d", buildLink, stage.Number, step.Number)
			}
			name := step.Name
			if len(build.Stages) > 1 {
				name = stage.Name + " › " + step.Name
			}
			failed = append(failed, name)
		}
	}
	switch {
	case len(failed) == 0:
		return "unknown", link
	case len(failed) > maxFailedSteps:
		return fmt.Sprintf("%s and %d more", strings.Join(failed[:maxFailedSteps], ", "), len(failed)-maxFailedSteps), link
	default:
		return strings.Join(failed, ", "), link
	}
}
