FROM golang:1.17.1 as go-builder
COPY go /go/nwn-do
RUN apt update \
    && cd /go/nwn-do \
    && go mod download \
    && go build -o ./bin/b \
    && mv bin/b /usr/local/bin/
FROM ubuntu:latest 
LABEL maintainer "urothis@gmail.com"
COPY --from=go-builder /usr/local/bin/b /usr/local/bin/upload
COPY nwsync /nwsync
RUN apt update \
    && apt-get install ca-certificates -y
CMD upload