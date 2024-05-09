FROM golang:1.22

LABEL authors="Polyanskiy KA"

WORKDIR /service

COPY . .

RUN go mod download
RUN go mod tidy

EXPOSE 8080

