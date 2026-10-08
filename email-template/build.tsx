// Renders the email as a Go html/template; email.go embeds the result and
// attaches emails/static as inline images.
import { pretty, render } from "react-email";

import Email from "./emails/email";
import type { EmailProps } from "./emails/email";

const props = {
  authorAvatar: "{{.AuthorAvatar}}",
  authorName: "{{.AuthorName}}",
  buildLink: "{{.BuildLink}}",
  buildNumber: "{{.BuildNumber}}",
  commitHash: "{{.CommitHash}}",
  commitLink: "{{.CommitLink}}",
  commitMessage: "{{.CommitMessage}}",
  duration: "{{.Duration}}",
  failedSteps: "{{.FailedSteps}}",
  imagesUrl: "cid:",
  refName: "{{.RefName}}",
  repository: "{{.Repository}}",
  serverHost: "{{.ServerHost}}",
  serverLink: "{{.ServerLink}}",
} satisfies EmailProps;

const html = await pretty(await render(<Email {...props} />));
await Bun.write(new URL("../email.html", import.meta.url), html);
