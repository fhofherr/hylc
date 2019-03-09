FROM golang:1.12.0-alpine3.9 as go-build

ARG git_tag
ARG git_hash
ARG build_time

RUN adduser -D go-build
WORKDIR /home/go-build
COPY --chown=go-build:go-build . .

USER go-build
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build\
    -ldflags "-s -X github.com/fhofherr/hylc/cmd.GitHash=$git_hash -X github.com/fhofherr/hylc/cmd.Version=$git_tag -X github.com/fhofherr/hylc/cmd.BuildTime=$build_time"\
    -a\
    -o hylc\
    .

FROM alpine:3.9 as run
RUN adduser -D hylc
WORKDIR /home/hylc
COPY --from=go-build --chown=hylc:hylc /home/go-build/hylc .

USER hylc
ENTRYPOINT ["./hylc"]
