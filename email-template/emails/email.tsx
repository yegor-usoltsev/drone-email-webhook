import type { CSSProperties, ReactNode } from "react";
import {
  Body,
  Button,
  Column,
  Container,
  Head,
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
        <style>{css}</style>
      </Head>
      <Body className="page" style={page}>
        <Preview>{`${props.failedSteps} failed: ${props.commitMessage}`}</Preview>
        <Container className="container" style={container}>
          <Section className="card" style={card}>
            <Section className="header" style={header}>
              <Row>
                <Column style={logoColumn}>
                  <Img
                    alt="Drone"
                    height="36"
                    src={`${props.imagesUrl}logo.png`}
                    width="36"
                  />
                </Column>
                <Column>
                  <Text style={heading}>Build #{props.buildNumber} failed</Text>
                  <Text style={repository}>{props.repository}</Text>
                </Column>
              </Row>
            </Section>
            <Section className="body" style={body}>
              <Field label="Failed">
                <span className="failed" style={failed}>
                  {props.failedSteps}
                </span>
              </Field>
              <Field label="Time">{props.duration}</Field>
              <Field label="Branch">
                <Icon src={`${props.imagesUrl}branch.png`} />
                <span style={refName}>{props.refName}</span>
              </Field>
              <Field label="Commit">
                <Icon src={`${props.imagesUrl}commit.png`} />
                <Link className="link" href={props.commitLink} style={link}>
                  {props.commitHash}
                </Link>
              </Field>
              <Field label="Message">{props.commitMessage}</Field>
              <Field label="Author">
                <Img
                  alt=""
                  height="18"
                  src={props.authorAvatar}
                  style={avatar}
                  width="18"
                />
                {props.authorName}
              </Field>
              <Button href={props.buildLink} style={button}>
                View build
              </Button>
            </Section>
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

function Field({ children, label }: { children: ReactNode; label: string }) {
  return (
    <Row className="row" style={row}>
      <Column className="muted" style={key}>
        {label}
      </Column>
      <Column className="content" style={value}>
        {children}
      </Column>
    </Row>
  );
}

function Icon({ src }: { src: string }) {
  return <Img alt="" height="16" src={src} style={icon} width="16" />;
}

// Colors follow the Drone UI themes. Apple Mail and some other clients apply
// the dark one; the rest keep the light theme.
// Body repeats its background on an inner cell, so the dark one covers both.
// Sections put their padding on an inner cell too, so phones override it there.
const css = `
@media (prefers-color-scheme: dark) {
  .page, .page > table > tbody > tr > td { background-color: #151a1e !important; }
  .card { background-color: #0b0d0f !important; border-color: #35354b !important; }
  .content { color: #b8b9c7 !important; }
  .muted { color: #9fa0b2 !important; }
  .link { color: #3886fa !important; }
  .row { border-color: #35354b !important; }
  .failed { background-color: #2b1214 !important; border-color: #ef554d !important; color: #ef554d !important; }
}
@media only screen and (max-width: 600px) {
  .container > tbody > tr > td { padding: 0 !important; }
  .header > tbody > tr > td { padding: 16px !important; }
  .body > tbody > tr > td { padding: 4px 16px 20px !important; }
}
`
  .replaceAll(/\s+/gu, " ")
  .trim();

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

const card = {
  backgroundColor: "#ffffff",
  border: "1px solid #e4e4eb",
  borderRadius: "6px",
} satisfies CSSProperties;

const header = {
  backgroundColor: "#0a3364",
  borderRadius: "5px 5px 0 0",
  padding: "20px 24px",
} satisfies CSSProperties;

const logoColumn = {
  verticalAlign: "top",
  width: "48px",
} satisfies CSSProperties;

const heading = {
  color: "#ffffff",
  fontSize: "20px",
  fontWeight: 600,
  lineHeight: "26px",
  margin: "0",
} satisfies CSSProperties;

const repository = {
  color: "#a9bcd8",
  fontSize: "14px",
  lineHeight: "20px",
  margin: "0",
} satisfies CSSProperties;

const body = {
  padding: "8px 24px 24px",
} satisfies CSSProperties;

const row = {
  borderBottom: "1px solid #e4e4eb",
} satisfies CSSProperties;

const key = {
  color: "#6b6d85",
  fontSize: "14px",
  lineHeight: "20px",
  padding: "12px 12px 12px 0",
  verticalAlign: "top",
  width: "72px",
} satisfies CSSProperties;

const value = {
  color: "#383946",
  fontSize: "14px",
  lineHeight: "20px",
  padding: "12px 0",
  verticalAlign: "top",
  wordBreak: "break-word",
} satisfies CSSProperties;

// A single box, so a long list of steps does not wrap into broken borders.
const failed = {
  backgroundColor: "#fff5f5",
  border: "1px solid #e43326",
  borderRadius: "4px",
  color: "#e43326",
  display: "inline-block",
  fontWeight: 600,
  padding: "0 6px",
} satisfies CSSProperties;

// Refs are identifiers, so they break anywhere instead of leaving the icon
// alone on a line.
const refName = {
  wordBreak: "break-all",
} satisfies CSSProperties;

const icon = {
  display: "inline",
  marginRight: "6px",
  verticalAlign: "-3px",
} satisfies CSSProperties;

const avatar = {
  ...icon,
  borderRadius: "50%",
  verticalAlign: "-4px",
} satisfies CSSProperties;

const link = {
  color: "#0063f7",
  textDecoration: "none",
} satisfies CSSProperties;

const button = {
  backgroundColor: "#0278d5",
  borderRadius: "4px",
  boxSizing: "border-box",
  color: "#ffffff",
  fontSize: "14px",
  fontWeight: 600,
  letterSpacing: "0.5px",
  lineHeight: "20px",
  marginTop: "24px",
  padding: "12px 20px",
  textAlign: "center",
  textTransform: "uppercase",
  width: "100%",
} satisfies CSSProperties;

const footer = {
  color: "#6b6d85",
  fontSize: "12px",
  lineHeight: "18px",
  margin: "20px 0",
  padding: "0 16px",
  textAlign: "center",
} satisfies CSSProperties;

const footerLink = {
  color: "#6b6d85",
  textDecoration: "underline",
} satisfies CSSProperties;
