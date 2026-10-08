import type { CSSProperties, ReactNode } from "react";
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
  Preview,
  Row,
  Section,
  Text,
} from "react-email";

export interface EmailProps {
  authorAvatar: string;
  authorName: string;
  buildLink: string;
  buildNumber: string;
  commitHash: string;
  commitLink: string;
  commitMessage: string;
  duration: string;
  failedSteps: string;
  // Prefix for the files in emails/static, "cid:" in the sent email.
  imagesUrl: string;
  refName: string;
  repository: string;
  serverHost: string;
  serverLink: string;
}

export default function Email(props: EmailProps) {
  return (
    <Html>
      <Head>
        <meta content="width=device-width, initial-scale=1" name="viewport" />
        <meta content="light dark" name="color-scheme" />
        <meta content="light dark" name="supported-color-schemes" />
        <style>{darkMode}</style>
      </Head>
      <Body className="page" style={page}>
        <Preview>{`${props.failedSteps} failed: ${props.commitMessage}`}</Preview>
        <Container style={container}>
          <Img
            alt="Drone"
            height="40"
            src={`${props.imagesUrl}logo.png`}
            style={logo}
            width="40"
          />
          <Section className="card" style={card}>
            <Text className="muted" style={repository}>
              {props.repository}
            </Text>
            <Row>
              <Column style={statusColumn}>
                <Img
                  alt="Failed"
                  height="24"
                  src={`${props.imagesUrl}failed.png`}
                  width="24"
                />
              </Column>
              <Column>
                <Heading className="heading" style={heading}>
                  Build #{props.buildNumber} failed
                </Heading>
              </Column>
            </Row>
            <Text className="content" style={message}>
              {props.commitMessage}
            </Text>
            <Text style={pills}>
              <Pill>
                <Img
                  alt=""
                  height="16"
                  src={props.authorAvatar}
                  style={avatar}
                  width="16"
                />
                {props.authorName}
              </Pill>
              <Pill>
                <Icon alt="Commit" src={`${props.imagesUrl}commit.png`} />
                <Link className="link" href={props.commitLink} style={link}>
                  {props.commitHash}
                </Link>
              </Pill>
              <Pill>
                <Icon alt="Ref" src={`${props.imagesUrl}branch.png`} />
                {props.refName}
              </Pill>
            </Text>
            <Text className="content" style={summary}>
              <span className="failed" style={failed}>
                {props.failedSteps}
              </span>{" "}
              failed after {props.duration}
            </Text>
            <Button className="button" href={props.buildLink} style={button}>
              View build
            </Button>
          </Section>
          <Text className="muted" style={footer}>
            Sent by Drone at{" "}
            <Link className="muted" href={props.serverLink} style={footerLink}>
              {props.serverHost}
            </Link>
          </Text>
        </Container>
      </Body>
    </Html>
  );
}

function Pill({ children }: { children: ReactNode }) {
  return (
    <span className="pill" style={pill}>
      {children}
    </span>
  );
}

function Icon({ alt, src }: { alt: string; src: string }) {
  return <Img alt={alt} height="16" src={src} style={icon} width="16" />;
}

// Colors follow the Drone UI themes. Apple Mail and some other clients apply
// the dark one; the rest keep the light theme.
// Body repeats its background on an inner cell, so the dark one covers both.
const darkMode = `
@media (prefers-color-scheme: dark) {
  .page, .page > table > tbody > tr > td { background-color: #151a1e !important; }
  .card { background-color: #0b0d0f !important; border-color: #35354b !important; }
  .heading { color: #d5d5dd !important; }
  .content { color: #b8b9c7 !important; }
  .muted { color: #9fa0b2 !important; }
  .pill { background-color: #262636 !important; color: #b8b9c7 !important; }
  .link { color: #3886fa !important; }
  .failed { background-color: #2b1214 !important; border-color: #ef554d !important; color: #ef554d !important; }
}
`;

const fontFamily =
  '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif';

const page = {
  backgroundColor: "#f8f9fa",
  color: "#383946",
  fontFamily,
} satisfies CSSProperties;

const container = {
  maxWidth: "600px",
  padding: "32px 16px",
} satisfies CSSProperties;

const logo = {
  margin: "0 auto 24px",
} satisfies CSSProperties;

const card = {
  backgroundColor: "#ffffff",
  border: "1px solid #e4e4eb",
  borderRadius: "6px",
  padding: "24px",
} satisfies CSSProperties;

const repository = {
  color: "#6b6d85",
  fontSize: "14px",
  lineHeight: "20px",
  margin: "0 0 8px",
} satisfies CSSProperties;

const statusColumn = {
  verticalAlign: "middle",
  width: "34px",
} satisfies CSSProperties;

const heading = {
  color: "#22222a",
  fontSize: "22px",
  fontWeight: 600,
  lineHeight: "28px",
  margin: "0",
} satisfies CSSProperties;

const message = {
  color: "#383946",
  fontSize: "16px",
  lineHeight: "24px",
  margin: "16px 0 12px",
} satisfies CSSProperties;

const pills = {
  fontSize: "14px",
  lineHeight: "20px",
  margin: "0 0 16px",
} satisfies CSSProperties;

const pill = {
  backgroundColor: "#f3f3fa",
  borderRadius: "4px",
  color: "#383946",
  display: "inline-block",
  margin: "0 6px 6px 0",
  padding: "4px 8px",
  wordBreak: "break-word",
} satisfies CSSProperties;

const avatar = {
  borderRadius: "50%",
  display: "inline",
  marginRight: "6px",
  verticalAlign: "-3px",
} satisfies CSSProperties;

const icon = {
  display: "inline",
  marginRight: "4px",
  verticalAlign: "-3px",
} satisfies CSSProperties;

const link = {
  color: "#0063f7",
  textDecoration: "none",
} satisfies CSSProperties;

const summary = {
  color: "#383946",
  fontSize: "14px",
  lineHeight: "24px",
  margin: "0 0 24px",
} satisfies CSSProperties;

const failed = {
  backgroundColor: "#fff5f5",
  border: "1px solid #e43326",
  borderRadius: "4px",
  color: "#e43326",
  fontWeight: 600,
  padding: "2px 8px",
} satisfies CSSProperties;

const button = {
  backgroundColor: "#0278d5",
  borderRadius: "4px",
  color: "#ffffff",
  fontSize: "14px",
  fontWeight: 600,
  letterSpacing: "0.5px",
  lineHeight: "20px",
  padding: "10px 20px",
  textTransform: "uppercase",
} satisfies CSSProperties;

const footer = {
  color: "#6b6d85",
  fontSize: "12px",
  lineHeight: "16px",
  margin: "24px 0 0",
  textAlign: "center",
} satisfies CSSProperties;

const footerLink = {
  color: "#6b6d85",
  textDecoration: "underline",
} satisfies CSSProperties;
