// Renders the email with sample data to out/email.html for a browser check.
import { render } from "react-email";

import Email from "./emails/email";
import type { EmailProps } from "./emails/email";

const props = {
  authorAvatar:
    "https://secure.gravatar.com/avatar/83c8d33e33a4999d1618d48ba0135e11?d=identicon",
  authorName: "Priya Raman",
  buildLink: "https://drone.example.com/octo-org/billing-api/4321/1/3",
  buildNumber: "4321",
  commitHash: "8f2e41f9",
  commitLink:
    "https://github.com/octo-org/billing-api/commit/8f2e41f9c1d04b7a9e3f6a2b5c8d7e1f0a9b4c3d",
  commitMessage: "Retry webhook delivery on 502 and 503 responses",
  duration: "3m12s",
  failedSteps: "test",
  imagesUrl: new URL("emails/static/", import.meta.url).href,
  refName: "feature/retry-webhooks",
  repository: "octo-org/billing-api",
  serverHost: "drone.example.com",
  serverLink: "https://drone.example.com",
} satisfies EmailProps;

await Bun.write(
  new URL("out/email.html", import.meta.url),
  await render(<Email {...props} />)
);
