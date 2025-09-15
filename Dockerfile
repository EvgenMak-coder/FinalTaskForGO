# Этап сборки
FROM golang:1.24.5-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/todo-app /app/todo-app
COPY web /app/web

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/scheduler.db

RUN mkdir -p /data


EXPOSE 7540

CMD ["/app/todo-app"]