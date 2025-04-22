import { render } from "@react-email/render";
import { exec } from "child_process";
import { mkdir, rm, writeFile } from "fs/promises";
import { join } from "path";
import { Email } from "./emails/email";

const ourDir = join(__dirname, "out");
await rm(ourDir, { recursive: true, force: true });
await mkdir(ourDir, { recursive: true });

const html = await render(<Email {...Email.PreviewProps} />, { pretty: true });

const outFile = join(ourDir, "email.html");
await writeFile(outFile, html, "utf-8");

exec(`open ${outFile}`);
