FROM golang

ADD . /go/src/github.com/ChrisMcKenzie/andrew/fulfillment

RUN go install github.com/ChrisMcKenzie/andrew/fulfillment

ENTRYPOINT /go/bin/fulfillment

EXPOSE 8080
