FROM golang

ADD . /go/src/github.com/cmpe281-sshekhar93/bitly/counter-service

RUN go get github.com/gorilla/mux

RUN go install github.com/cmpe281-sshekhar93/bitly/counter-service

ENTRYPOINT /go/bin/counter-service

EXPOSE 3000
