FROM golang:1.24.1-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o /todo-server

FROM alpine:latest
RUN apk --no-cache add ca-certificates libc6-compat
WORKDIR /root/
COPY --from=builder /todo-server .
COPY web ./web/

EXPOSE 7540
CMD ["./todo-server"]