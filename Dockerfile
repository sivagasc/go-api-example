FROM golang:1.13 as mybuildstage

RUN mkdir /go-example

WORKDIR /go-example

COPY . /go-example/

# ADD . /go/src/github.com/sivagasc/go-api-example 

# RUN go install github.com/sivagasc/go-api-example 

RUN make build

# ENTRYPOINT  /go-example/bin/server

# EXPOSE 8090

FROM ubuntu

COPY --from=mybuildstage /go-example/bin/server .
COPY --from=mybuildstage /go-example/.env .

CMD [ "./server" ]

EXPOSE 8090