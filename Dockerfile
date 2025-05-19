FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o weatherapp ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/weatherapp .

COPY internal/web/templates/ internal/web/templates/
COPY internal/web/static/ internal/web/static/

ENV PORT=8080

EXPOSE 8080

CMD ["./weatherapp"]
