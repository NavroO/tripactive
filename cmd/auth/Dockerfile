FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY cmd/auth/.env /app/.env


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /auth ./cmd/auth

FROM golang:1.25
WORKDIR /app
COPY --from=builder /auth .
EXPOSE 8081
ENTRYPOINT ["/app/auth"]
