FROM golang:1.20-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/zepif/EtherUSDC
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/EtherUSDC /go/src/github.com/zepif/EtherUSDC


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/EtherUSDC /usr/local/bin/EtherUSDC
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["EtherUSDC"]
