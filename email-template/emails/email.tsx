import {
  Body,
  Button,
  Column,
  Container,
  Head,
  Heading,
  Html,
  Img,
  Link,
  Row,
  Section,
  Tailwind,
  Text,
} from "@react-email/components";
import { readFileSync } from "fs";
import { join } from "path";

const droneLogoPng = `data:image/png;base64,${readFileSync(join(__dirname, "../images/drone-logo.png")).toString("base64")}`;
const referencePng = `data:image/png;base64,${readFileSync(join(__dirname, "../images/reference.png")).toString("base64")}`;
const commitPng = `data:image/png;base64,${readFileSync(join(__dirname, "../images/commit.png")).toString("base64")}`;

export interface EmailProps {
  subject: string;
  from: string;
  to: string;
  header: string;
  repository: string;
  reference: string;
  commitHash: string;
  commitMessage: string;
  authorAvatar: string;
  authorName: string;
  droneBuildLink: string;
  droneServerHost: string;
  droneServerLink: string;
}

export const Email = ({
  subject,
  from,
  to,
  header,
  repository,
  reference,
  commitHash,
  commitMessage,
  authorAvatar,
  authorName,
  droneBuildLink,
  droneServerHost,
  droneServerLink,
}: EmailProps) => {
  return (
    <Tailwind
      config={{
        presets: [require("tailwindcss-preset-email")],
        important: false,
      }}
    >
      <Html>
        <Head />
        <Body className="bg-slate-100 font-sans text-[16px] text-slate-800 dark:bg-slate-900 dark:text-slate-200">
          <Container>
            <Img
              className="mx-auto my-6"
              height="64"
              src={droneLogoPng}
              width="64"
            />
            <Section className="rounded-lg bg-slate-50 p-4 shadow dark:bg-slate-950">
              <Heading className="m-0 rounded bg-red-500 px-4 py-2 text-center text-lg text-slate-100 dark:bg-red-700">
                {header}
              </Heading>
              <Section className="my-6 min-w-80 text-sm">
                <Row className="pb-2">
                  <Column className="w-1/4 pr-1">Repository</Column>
                  <Column className="line-clamp-3 text-ellipsis break-all">
                    {repository}
                  </Column>
                </Row>
                <Row className="pb-2">
                  <Column className="w-1/4 pr-1">Reference</Column>
                  <Column className="line-clamp-3 text-ellipsis break-all">
                    <Img
                      className="inline align-middle"
                      height="24"
                      src={referencePng}
                      width="24"
                    />{" "}
                    {reference}
                  </Column>
                </Row>
                <Row className="pb-2">
                  <Column className="w-1/4 pr-1">Commit</Column>
                  <Column className="line-clamp-3 text-ellipsis break-all">
                    <Img
                      className="inline align-middle"
                      height="24"
                      src={commitPng}
                      width="24"
                    />{" "}
                    {commitHash}
                    <br />
                    {commitMessage}
                  </Column>
                </Row>
                <Row>
                  <Column className="w-1/4 pr-1">Author</Column>
                  <Column className="line-clamp-3 text-ellipsis break-all">
                    <Img
                      className="inline rounded-full bg-white align-middle"
                      height="24"
                      src={authorAvatar}
                      width="24"
                    />{" "}
                    {authorName}
                  </Column>
                </Row>
              </Section>
              <Section className="text-center">
                <Button
                  className="rounded bg-sky-500 px-6 py-3 text-center font-semibold text-slate-100 no-underline dark:bg-sky-700"
                  href={droneBuildLink}
                >
                  View build
                </Button>
              </Section>
            </Section>
            <Text className="text-center text-xs text-slate-500">
              You&apos;re receiving this email because of your account on{" "}
              <Link
                className="text-sky-500 no-underline dark:text-sky-700"
                href={droneServerLink}
              >
                {droneServerHost}
              </Link>
            </Text>
          </Container>
        </Body>
      </Html>
    </Tailwind>
  );
};

Email.PreviewProps = {
  subject:
    "[harness/drone] Failed build #4321 for refs/heads/feature/add-notifications (8f2e41f9)",
  from: "ci@example.com",
  to: "sarah.johnson@harness.io",
  header: "Build #4321 has failed",
  repository: "harness/drone",
  reference: "refs/heads/feature/add-notifications",
  commitHash: "8f2e41f9",
  commitMessage: "fix: Handle edge case in notification delivery",
  authorAvatar:
    "https://secure.gravatar.com/avatar/83c8d33e33a4999d1618d48ba0135e11?d=identicon",
  authorName: "Sarah Johnson",
  droneBuildLink: "https://ci.harness.io/harness/drone/4321",
  droneServerHost: "ci.harness.io",
  droneServerLink: "https://ci.harness.io",
} as EmailProps;

Email.BuildProps = {
  subject: "{{.Subject}}",
  from: "{{.From}}",
  to: "{{.To}}",
  header: "{{.Header}}",
  repository: "{{.Repository}}",
  reference: "{{.Reference}}",
  commitHash: "{{.CommitHash}}",
  commitMessage: "{{.CommitMessage}}",
  authorAvatar: "{{.AuthorAvatar}}",
  authorName: "{{.AuthorName}}",
  droneBuildLink: "{{.DroneBuildLink}}",
  droneServerHost: "{{.DroneServerHost}}",
  droneServerLink: "{{.DroneServerLink}}",
} as EmailProps;

export default Email;
