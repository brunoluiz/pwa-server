#
# BUILD
#
FROM golang:1 AS build

## Set the Current Working Directory inside the container
WORKDIR /app

## Copy go mod and sum files
COPY go.mod go.sum ./

## Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

## Copy the source from the current directory to the Working Directory inside the container
COPY . .

RUN make install-tools test lint

## Build the Go app
RUN go build -o main ./cmd

#
# RUNTIME
#
FROM alpine:3

## Add Maintainer Info
LABEL maintainer="Bruno Luiz Silva <contact@brunoluiz.net>"

RUN apk --no-cache add ca-certificates

COPY --from=build /app/main /app

EXPOSE 80

ENTRYPOINT exec /app $0 $@
