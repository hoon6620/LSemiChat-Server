FROM golang:1.14-alpine as build

WORKDIR /go/src/app

ADD . .

# realizeはgomodに対応していないらしい
RUN apk update \
    && apk add git \
    && GO111MODULE="on" go get github.com/yyh-gl/realize@master \
    && go get -v ./... \
    && go build -o ./binary/app ./cmd

FROM alpine

WORKDIR /app

COPY --from=build /go/src/app .

RUN addgroup go \
  && adduser -D -G go go \
  && chown -R go:go /app/app

CMD ["./app"]