FROM golang:1.25-alpine

WORKDIR /backend

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o backend ./cmd/main.go

EXPOSE 8080

CMD ["./backend"]