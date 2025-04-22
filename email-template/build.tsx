import { render } from "@react-email/render";
import { mkdir, rm, writeFile } from "fs/promises";
import { join } from "path";
import { Email } from "./emails/email";

const ourDir = join(__dirname, "out");
await rm(ourDir, { recursive: true, force: true });
await mkdir(ourDir, { recursive: true });

const html = await render(<Email {...Email.BuildProps} />, { pretty: false });

const outFile = join(ourDir, "email.html");
await writeFile(outFile, html, "utf-8");
