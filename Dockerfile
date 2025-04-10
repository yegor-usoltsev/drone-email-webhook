FROM alpine:latest
ENTRYPOINT ["/drone-email-webhook"]
COPY drone-email-webhook /
