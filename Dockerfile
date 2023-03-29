FROM golang:1.20-alpine as build

RUN mkdir -p /go/identity
WORKDIR /go/identity

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o /opt/identity main.go

FROM alpine:latest

COPY --from=build /opt/identity /opt/identity

ENTRYPOINT ["/opt/identity"]