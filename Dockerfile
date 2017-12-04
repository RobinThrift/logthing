FROM golang:latest as builder
WORKDIR /go/src/github.com/RobinThrift/logthing/
ADD . .
RUN go get github.com/sendgrid/sendgrid-go
RUN GOOS=linux make build

FROM alpine:latest

ENV LT_BUFFER 100
ENV LT_INTERVAL "24h"
ENV LT_SENDER "PLS"
ENV LT_RECIPIENT "SET"
ENV LT_SENDGRID_API_KEY "ME"

WORKDIR /root/
COPY --from=builder /go/src/github.com/RobinThrift/logthing/logthing .
CMD ["./logthing"]
