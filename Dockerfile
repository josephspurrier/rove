FROM golang:latest

ADD . /go/src/github.com/josephspurrier/rove

WORKDIR /go

RUN go build github.com/josephspurrier/rove/cmd/rove

CMD ["/go/rove", "migrate", "all", "testdata/success.sql"]