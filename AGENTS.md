# drone-email-webhook

Go webhook receiver that emails the commit author when a Drone build fails.

## Email template

- `email.html` is generated; edit `email-template/emails/email.tsx`, then run `bun run build` in `email-template` and commit both. CI fails when `email.html` is stale.
- `build.tsx` maps every template prop to a `{{.Field}}` of `emailData` in `email.go`; keep the names in sync. Go computes every value, including fallbacks, so the templates have no conditionals.
- `email.txt` is the hand-written plain-text part; update it with the HTML.
- Images live in `email-template/emails/static`. `email.go` embeds and attaches them inline, and the HTML references them as `cid:<file name>`, because Gmail and Outlook block `data:` URIs.
- Colors follow the Drone UI light and dark themes; dark mode is a `prefers-color-scheme` block in `email.tsx`.
- `bun run preview` writes `out/email.html` with sample data; `bun run lint` runs TypeScript and Ultracite.

## Checks

`go test -race ./...` starts Mailpit through testcontainers and needs Docker; `golangci-lint run` covers Go style.
