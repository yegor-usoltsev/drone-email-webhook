# Classic Outlook drops image margins and badge padding

Classic (Word-engine) Outlook strips `margin` on images and supports `padding` only on table cells, so in `email-template/emails/email.tsx`:

- the Author avatar relies on `margin-right: 6px`, and its text can touch the image;
- the Failed badge `span` uses `padding: 0 6px` and loses its inner spacing.

Evidence (guidelines, not rendered in Outlook): [Litmus: Outlook rendering differences](https://www.litmus.com/blog/a-guide-to-rendering-differences-in-microsoft-outlook-clients), [caniemail: padding, note 1](https://www.caniemail.com/features/css-padding/).

Branch and Commit now use presentation tables with the image in its own cell and spacing on the `td`.

Possible fix for the remaining issues: the same table layout for the Author avatar and the badge styles on an auto-width `td`.

Open decision: whether classic Outlook rendering is worth the extra markup.
