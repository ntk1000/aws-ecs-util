FROM golang:1.9.7-alpine

WORKDIR /go/src/app
COPY . .

RUN apk update \
      && apk add --no-cache git \
      && go get -u github.com/golang/dep/cmd/dep \
      && dep ensure \
      && go build -o gofe

