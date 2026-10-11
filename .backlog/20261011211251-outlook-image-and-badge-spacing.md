# Classic Outlook drops image margins and badge padding

Classic (Word-engine) Outlook strips `margin` on images and supports `padding` only on table cells, so in `email-template/emails/email.tsx`:

- the Branch, Commit and Author icons rely on `margin-right: 6px`, and their text can touch the image;
- the Failed badge `span` uses `padding: 0 6px` and loses its inner spacing.

Evidence (guidelines, not rendered in Outlook): [Litmus: Outlook rendering differences](https://www.litmus.com/blog/a-guide-to-rendering-differences-in-microsoft-outlook-clients), [caniemail: padding, note 1](https://www.caniemail.com/features/css-padding/).

Possible fix: a presentation table per value with the image in its own cell (`display: block`, `padding-right: 6px` on the `td`) and the badge styles on an auto-width `td`.

Open decision: whether classic Outlook rendering is worth the extra markup.
