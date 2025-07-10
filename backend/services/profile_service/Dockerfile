FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /profile_service ./cmd/app/main.go
CMD ["/profile_service"]