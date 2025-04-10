# drone-email-webhook

[![Build Status](https://github.com/yegor-usoltsev/drone-email-webhook/actions/workflows/ci.yml/badge.svg)](https://github.com/yegor-usoltsev/drone-email-webhook/actions)
[![GitHub Release](https://img.shields.io/github/v/release/yegor-usoltsev/drone-email-webhook?sort=semver)](https://github.com/yegor-usoltsev/drone-email-webhook/releases)
[![Docker Image (docker.io)](https://img.shields.io/docker/v/yusoltsev/drone-email-webhook?label=docker.io&sort=semver)](https://hub.docker.com/r/yusoltsev/drone-email-webhook)
[![Docker Image (ghcr.io)](https://img.shields.io/docker/v/yusoltsev/drone-email-webhook?label=ghcr.io&sort=semver)](https://github.com/yegor-usoltsev/drone-email-webhook/pkgs/container/drone-email-webhook)
[![Docker Image Size](https://img.shields.io/docker/image-size/yusoltsev/drone-email-webhook?sort=semver&arch=amd64)](https://hub.docker.com/r/yusoltsev/drone-email-webhook/tags)

Webhook listener for Drone CI / CD notifying commit authors of failed builds via email.

![Screenshot](https://raw.githubusercontent.com/yegor-usoltsev/drone-email-webhook/main/.github/screenshot.png)

## Usage

### Starting the webhook listener

Pull and run the container:

```bash
docker run -d \
  -p 3000:3000 \
  -e DRONE_SECRET=your_webhook_secret \
  yusoltsev/drone-email-webhook:latest
```

See the [Environment Variables](#environment-variables) table below for all available configuration options.

### Configuring Drone

Configure your Drone server to send webhooks by setting the following environment variables:

```yaml
DRONE_WEBHOOK_ENDPOINT: http://your-webhook-host:3000
DRONE_WEBHOOK_SECRET: your_webhook_secret # Must match DRONE_SECRET in webhook container
```

For more information about Drone webhooks configuration, please refer to
the [official Drone documentation](https://docs.drone.io/webhooks/overview/).

Note: Make sure to replace the placeholder values with your actual configuration values.

### Environment Variables

| KEY                         | TYPE      | DEFAULT           |
|-----------------------------|-----------|-------------------|
| `DRONE_SECRET`              | `String`  |                   |
| `DRONE_SERVER_HOST`         | `String`  | `0.0.0.0`         |
| `DRONE_SERVER_PORT`         | `Integer` | `3000`            |
| `DRONE_EMAIL_SMTP_HOST`     | `String`  | `localhost`       |
| `DRONE_EMAIL_SMTP_PORT`     | `Integer` | `25`              |
| `DRONE_EMAIL_SMTP_USERNAME` | `String`  | `drone`           |
| `DRONE_EMAIL_SMTP_PASSWORD` | `String`  | `drone`           |
| `DRONE_EMAIL_FROM`          | `String`  | `drone@localhost` |

## Docker Images

This application is delivered as a multi-platform Docker image and is available for download from two image registries
of choice: [yusoltsev/drone-email-webhook](https://hub.docker.com/r/yusoltsev/drone-email-webhook)
and [ghcr.io/yegor-usoltsev/drone-email-webhook](https://github.com/yegor-usoltsev/drone-email-webhook/pkgs/container/drone-email-webhook).
Images are tagged as follows:

- `latest` - Tracks the latest released version, which is typically tagged with a version number. This tag is
  recommended for most users as it provides the most stable version.
- `edge` - Tracks the latest commits to the `main` branch.
- `vX.Y.Z` (e.g., `v1.2.3`) - Represents a specific released version.

## Versioning

This project uses [Semantic Versioning](https://semver.org)

## Contributing

Pull requests are welcome. For major changes,
please [open an issue](https://github.com/yegor-usoltsev/drone-email-webhook/issues/new) first to discuss what you would
like to change. Please make sure to update tests as appropriate.

## License

[MIT](https://github.com/yegor-usoltsev/drone-email-webhook/blob/main/LICENSE)
