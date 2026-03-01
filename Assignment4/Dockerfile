# Этап 1: Сборка приложения
FROM golang:1.25-alpine AS builder 
WORKDIR /app

# Копируем модули и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь остальной код
COPY . .

# Собираем бинарный файл (называем его myapp)
RUN go build -o myapp ./cmd/api/main.go

# Этап 2: Создание финального легкого образа (Multi-stage)
FROM alpine:latest
WORKDIR /app

# Копируем готовый файл из первого этапа
COPY --from=builder /app/myapp .

EXPOSE 8080

# Команда для запуска сервера
CMD ["./myapp"]