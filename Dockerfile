FROM golang:1.13

RUN mkdir /go-example

WORKDIR /go-example

COPY . /go-example/

# ADD . /go/src/github.com/sivagasc/go-api-example 

# RUN go install github.com/sivagasc/go-api-example 

RUN make build

ENTRYPOINT  /go-example/bin/server

EXPOSE 8090