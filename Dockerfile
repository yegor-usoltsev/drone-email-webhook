FROM alpine:latest
RUN apk add --no-cache --update ca-certificates tzdata
ARG TARGETPLATFORM
ENTRYPOINT ["/drone-email-webhook"]
COPY $TARGETPLATFORM/drone-email-webhook /
