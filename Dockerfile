FROM golang:1.20-alpine AS build
WORKDIR /app
COPY . .
RUN go mod download && \
    go build -ldflags="-s -w"

FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache --update ca-certificates
COPY --from=build /app/drone-email-webhook .
CMD ["/app/drone-email-webhook"]
