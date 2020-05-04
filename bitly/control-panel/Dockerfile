FROM golang

ADD . /go/src/github.com/cmpe281-sshekhar93/bitly/control-panel

RUN go get github.com/catinello/base62

RUN go get github.com/go-sql-driver/mysql

RUN go get github.com/gorilla/mux

RUN go get github.com/streadway/amqp

RUN go install github.com/cmpe281-sshekhar93/bitly/control-panel

ENTRYPOINT /go/bin/control-panel

EXPOSE 8000
EXPOSE 5672
