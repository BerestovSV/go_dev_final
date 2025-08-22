FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY . .

ENV CGO_ENABLED=0

RUN go mod download
RUN GOOS=linux go build -o /todo-server

FROM alpine:latest

# Минимальные зависимости (только ca-certificates)
RUN apk add --no-cache ca-certificates

WORKDIR /root/
COPY --from=builder /todo-server .
COPY web ./web/

EXPOSE 7540
CMD ["./todo-server"]