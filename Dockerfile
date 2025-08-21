# Билд этап
FROM golang:1.24-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git

# Создаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для скачивания зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарник с оптимизациями
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o todo-server .

# Финальный этап
FROM alpine:3.19

# Устанавливаем зависимости для runtime
RUN apk add --no-cache ca-certificates tzdata

# Создаем пользователя для безопасности
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Создаем директории
RUN mkdir -p /app/data && chown -R appuser:appgroup /app

# Переключаемся на непривилегированного пользователя
USER appuser

# Копируем бинарник из билд этапа
COPY --from=builder --chown=appuser:appgroup /app/todo-server /app/
COPY --from=builder --chown=appuser:appgroup /app/web /app/web/

# Рабочая директория
WORKDIR /app

# Экспортируем порт
EXPOSE 7540

# Переменные окружения по умолчанию
ENV TODO_PORT=7540
ENV TODO_DBFILE=/app/data/scheduler.db
ENV TODO_PASSWORD=""

# Точка входа
ENTRYPOINT ["/app/todo-server"]