FROM golang:1.13.1-alpine3.10 as builder

RUN apk update && apk upgrade && \
    apk --update add git gcc make

WORKDIR /go/src/github.com/jayvib/app

COPY . .

ENV GO111MODULE on

RUN make front-end

FROM alpine:latest

RUN apk update && apk upgrade && \
    apk --update --no-cache add tzdata && \
    apk add --no-cache bash && \
    mkdir /app

WORKDIR /app

EXPOSE 8080

COPY --from=builder /go/src/github.com/jayvib/app/front-end.linux /app

COPY web/ /app/web

CMD /app/front-end.linux
