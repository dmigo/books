ARG GO_VERSION=1.16.3-buster

# Image to compile the application
FROM golang:${GO_VERSION} AS build

WORKDIR /books
COPY src/ .
RUN pwd

RUN go mod download
RUN go mod verify
RUN go build -o api .

# Image to run the application
FROM golang:${GO_VERSION} AS main

COPY --from=build /books/api ./api
ENTRYPOINT ["./api"]
