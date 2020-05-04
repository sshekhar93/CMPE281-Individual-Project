FROM golang

ADD . /go/src/github.com/cmpe281-sshekhar93/bitly/link-redirect-server

RUN go get github.com/catinello/base62

RUN go get github.com/go-sql-driver/mysql

RUN go get github.com/gorilla/mux

RUN go get github.com/streadway/amqp

RUN go install github.com/cmpe281-sshekhar93/bitly/link-redirect-server

ENTRYPOINT /go/bin/link-redirect-server

EXPOSE 8000
