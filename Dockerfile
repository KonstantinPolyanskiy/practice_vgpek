FROM golang:1.22

LABEL authors="Polyanskiy KA"

ENV GOPATH=/

COPY ./ ./

RUN go mod download

RUN go mod tidy

RUN go build -o backend ./cmd/app/main.go

CMD ["./backend"]

EXPOSE 8080

