FROM golang:1.22.2-alpine AS build
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w"

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/drone-email-webhook .
CMD ["/app/drone-email-webhook"]
