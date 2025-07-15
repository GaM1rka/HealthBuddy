FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /auth_service ./cmd/app/main.go
CMD ["/auth_service"]