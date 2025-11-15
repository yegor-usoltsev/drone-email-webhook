FROM alpine:latest
RUN apk add --no-cache --update ca-certificates tzdata
ENTRYPOINT ["/drone-email-webhook"]
COPY drone-email-webhook /
