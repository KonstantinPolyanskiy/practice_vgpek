FROM golang:1.22

LABEL authors="Polyanskiy KA"

WORKDIR /service

COPY go.mod ./
RUN go mod download
RUN go mod tidy

COPY . .

EXPOSE 8080

