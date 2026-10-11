# drone-email-webhook

Go webhook receiver that emails the commit author when a Drone build fails.

## Email template

- `email.html` is generated; edit `email-template/emails/email.tsx` and commit both. The prek pre-commit hook in `.pre-commit-config.yaml` runs `bun run build` and `bun run lint` in `email-template` whenever a file there changes, and CI runs the same hooks. Renovate leaves `email-template` alone because the template changes rarely.
- `build.tsx` maps every template prop to a `{{.Field}}` of `emailData` in `email.go`; keep the names in sync. Go computes every value, including fallbacks, so the templates have no conditionals.
- `email.txt` is the hand-written plain-text part; update it with the HTML.
- Images live in `email-template/emails/static`. `email.go` embeds and attaches them inline, and the HTML references them as `cid:<file name>`, because Gmail and Outlook block `data:` URIs.
- Give every `<img>` both numeric `width`/`height` attributes and CSS `width`/`height` in `style`: macOS Mail ignores the attributes on `cid:` images and draws them at natural size, while classic Outlook ignores the CSS.
- Colors follow the Drone UI light and dark themes; dark mode is a `prefers-color-scheme` block in `email.tsx`.
- `bun run preview` writes `out/email.html` with sample data.
- `bun run lint` type-checks and runs Ultracite (oxlint and oxfmt); `bun run lint:fix` applies fixes.
- `tsc` is TypeScript 7 from `@typescript/native`. The `typescript` package stays on 6 because `eslint-plugin-github` loads typescript-eslint, which rejects 7.

## Checks

`go test -race ./...` starts Mailpit through testcontainers and needs Docker; `golangci-lint run` covers Go style.
