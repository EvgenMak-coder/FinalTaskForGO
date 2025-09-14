# Этап сборки
FROM golang:1.24.5-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app ./cmd

# Этап запуска
FROM alpine:latest

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/todo-app /app/todo-app

# Копируем статические файлы
COPY web /app/web

# Устанавливаем переменные окружения по умолчанию
ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

# Создаем директорию для базы данных
RUN mkdir -p /data

# Открываем порт
EXPOSE 7540

# Запускаем приложение
CMD ["/app/todo-app"]