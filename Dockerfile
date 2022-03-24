# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:1.14 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

ADD . .

RUN go build -o /ckb-node-websocket-client

##
## Deploy
##
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /ckb-node-websocket-client /ckb-node-websocket-client

EXPOSE 8080

USER nonroot:nonroot

CMD /ckb-node-websocket-client
