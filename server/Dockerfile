# 1. Билд-образ: golang на базе Debian
FROM golang:1.23.1 AS builder

# Устанавливаем необходимые зависимости (компилятор и библиотеки)
RUN apt-get update && apt-get install -y gcc g++ make libwebp-dev

WORKDIR /app

# Копируем зависимости и проект
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем бинарник с включенным CGO
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main ./cmd/offliner/main.go

# 2. Финальный образ: минимальный Debian (Bookworm)
FROM debian:bookworm-slim

# Устанавливаем необходимые библиотеки для запуска
RUN apt-get update && apt-get install -y libwebp7 ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Копируем бинарник из билд-образа
COPY --from=builder /app/main /app/main

# Устанавливаем права на выполнение
RUN chmod +x /app/main

EXPOSE 8080

# Проверяем содержимое и запускаем приложение
RUN ls -l /app
CMD ["./main"]
