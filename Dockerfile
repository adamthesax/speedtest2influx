FROM golang:alpine AS builder
ENV CGO_ENABLED=0
RUN apk update && apk add --no-cache git ca-certificates
WORKDIR $GOPATH/src/github.com/adamthesax/speedtest2influx

COPY . .
RUN go get -d -v
RUN go build -o /go/bin/speedtest2influx

FROM scratch
COPY --from=builder /go/bin/speedtest2influx /go/bin/speedtest2influx
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/speedtest2influx"]
