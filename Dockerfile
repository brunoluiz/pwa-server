#
# BUILD
#
FROM golang:1-alpine AS build

## Add Maintainer Info
LABEL maintainer="Bruno Luiz Silva <contact@brunoluiz.net>"

## Set the Current Working Directory inside the container
WORKDIR /app

## Copy go mod and sum files
COPY go.mod go.sum ./

## Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

## Copy the source from the current directory to the Working Directory inside the container
COPY . .

## Build the Go app
RUN go build -o main ./cmd

# Command to run the executable
CMD ["./main"]

#
# RUNTIME
#
FROM alpine:3

RUN apk --no-cache add ca-certificates

COPY --from=build /app/main /app

EXPOSE 80

ENTRYPOINT exec /app $0 $@
