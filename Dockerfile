# golang:alpine is the alpine image with the go tools added.. manually add git
FROM golang:alpine as builder

ENV GO111MODULE=on

# install gcc for compilation
RUN apk add --update gcc musl-dev
# Set an env var that matches github repo name
# ENV CGO_ENABLED=0
ENV SRC_DIR=${HOME}/go/src/github.com/anothrNick/send-query-result/

# Add the source code:
ADD . $SRC_DIR

# Build it:
RUN cd $SRC_DIR;\
    apk add --no-cache git;\
    go build -o app;

ENTRYPOINT ["/go/src/github.com/anothrNick/send-query-result/app"]

# alpine production environment
# copy binary for smallest image size
FROM alpine:3.7

ARG VERSION=latest
ENV VERSION=${VERSION}

RUN apk add --no-cache ca-certificates

ENV GIN_MODE=release

COPY --from=builder /go/src/github.com/anothrNick/send-query-result/app /bin/app

ENTRYPOINT ["/bin/app"]
