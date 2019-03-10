# Based on https://github.com/chemidy/smallest-secured-golang-docker-image

FROM golang:1.12.1-alpine3.9 as gobuilder

ARG git_tag
ARG git_hash
ARG build_time

# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk add --no-cache git build-base ca-certificates tzdata
RUN update-ca-certificates

RUN adduser -D -g '' -h /opt/app app
WORKDIR /opt/app
COPY --chown=app:app . .

USER app
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make hylc GIT_TAG=$git_tag GIT_HASH=$git_hash BUILD_TIME=$build_time

FROM scratch
COPY --from=gobuilder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=gobuilder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=gobuilder /etc/passwd /etc/passwd

COPY --from=gobuilder /opt/app/hylc /opt/app/hylc

USER app
WORKDIR /opt/app

ENTRYPOINT ["/opt/app/hylc"]
CMD ["serve"]
