# drone-email-webhook

[![Build Status](https://github.com/yegor-usoltsev/drone-email-webhook/actions/workflows/ci.yml/badge.svg)](https://github.com/yegor-usoltsev/drone-email-webhook/actions)
[![GitHub Release](https://img.shields.io/github/v/release/yegor-usoltsev/drone-email-webhook?sort=semver)](https://github.com/yegor-usoltsev/drone-email-webhook/releases)
[![Docker Image (docker.io)](https://img.shields.io/docker/v/yusoltsev/drone-email-webhook?label=docker.io&sort=semver)](https://hub.docker.com/r/yusoltsev/drone-email-webhook)
[![Docker Image (ghcr.io)](https://img.shields.io/docker/v/yusoltsev/drone-email-webhook?label=ghcr.io&sort=semver)](https://github.com/yegor-usoltsev/drone-email-webhook/pkgs/container/drone-email-webhook)

Webhook listener for Drone CI / CD notifying commit authors of failed builds via email.

## Usage

TODO

### Environment Variables

| KEY                           | TYPE       | DEFAULT                            |
|-------------------------------|------------|------------------------------------|
| `APP_SERVER_HOST`             | `String`   | `0.0.0.0`                          |
| `APP_SERVER_PORT`             | `Integer`  | `8080`                             |
| `APP_SERVER_MAX_HEADER_BYTES` | `Integer`  | `4096` (4 * 1024 = 4 KB)           |
| `APP_SERVER_MAX_BODY_BYTES`   | `Integer`  | `1048576` (1 * 1024 * 1024 = 1 MB) |
| `APP_SERVER_READ_TIMEOUT`     | `Duration` | `15s`                              |
| `APP_SERVER_HANDLER_TIMEOUT`  | `Duration` | `10s`                              |
| `APP_SERVER_WRITE_TIMEOUT`    | `Duration` | `15s`                              |
| `APP_SERVER_IDLE_TIMEOUT`     | `Duration` | `120s`                             |
| `APP_SERVER_SHUTDOWN_TIMEOUT` | `Duration` | `15s`                              |
| `APP_EMAIL_SMTP_HOST`         | `String`   | `localhost`                        |
| `APP_EMAIL_SMTP_PORT`         | `Integer`  | `1025`                             |
| `APP_EMAIL_SMTP_USERNAME`     | `String`   | `maildev`                          |
| `APP_EMAIL_SMTP_PASSWORD`     | `String`   | `maildev`                          |
| `APP_EMAIL_FROM`              | `String`   | `Drone <drone@example.com>`        |

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
