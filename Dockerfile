# Build
FROM golang:1.12.8-alpine3.10 as builder

RUN apk update && apk upgrade && \
    apk --update add git gcc make


WORKDIR /go/src/github.com/jayvib/clean-architecture

COPY . .

ENV GO111MODULE on

RUN make engine

# Distribute
FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    mkdir /app

WORKDIR /app

EXPOSE 9090 

COPY --from=builder /go/src/github.com/jayvib/clean-architecture/engine.linux /app

CMD /app/engine.linux



