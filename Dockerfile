# Этап 1: Сборка приложения
FROM golang:1.22 AS builder

# Установка рабочей директории
WORKDIR /app

# Копируем go.mod и go.sum и загружаем зависимости
COPY go.mod go.sum ./
RUN go mod download

COPY .env .

# Копируем все файлы приложения
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp ./cmd

# Этап 2: Создание финального образа
FROM alpine:latest

# Копируем собранное приложение из этапа сборки
COPY --from=builder /app/myapp /usr/local/bin/data_processing
COPY --from=builder /app/.env /usr/local/bin/.env
COPY --from=builder /app/.env /usr/local/.env

# Указываем команду для запуска приложения
CMD ["data_processing"]

# Открываем порт, если необходимо
EXPOSE 50052
