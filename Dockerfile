FROM golang:1.22-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/github.com/zepif/EtherUSDC
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/EtherUSDC /go/src/github.com/zepif/EtherUSDC


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/EtherUSDC /usr/local/bin/EtherUSDC
COPY start.sh /usr/local/bin/start.sh
RUN apk add --no-cache ca-certificates wget iputils

RUN chmod +x /usr/local/bin/start.sh

# ENTRYPOINT ["EtherUSDC"]
CMD ["EtherUSDC"]
