FROM golang:1.23.5-alpine AS build
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/drone-email-webhook .
CMD ["/app/drone-email-webhook"]
