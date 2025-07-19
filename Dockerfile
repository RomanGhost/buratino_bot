# Базовый образ
FROM golang:1.24 as builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы проекта
COPY go.mod go.sum ./
#ENV GOPROXY=direct
RUN go mod download

# Копируем весь код
COPY . .

# Сборка приложения
RUN go build -o app ./cmd

# Финальный образ
FROM debian:bookworm-slim

# Устанавливаем доверенные корневые сертификаты
RUN apt-get update && apt-get install -y ca-certificates && update-ca-certificates && rm -rf /var/lib/apt/lists/*

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем скомпилированное приложение из предыдущего контейнера
COPY --from=builder /app/app .
COPY --from=builder /app/.env .
# Запускаем приложение
CMD ["./app"]
