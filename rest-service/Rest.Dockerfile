FROM golang:1.24.2 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o /app/bin/rest-service ./rest-service/cmd/main.go

# Финальный минимальный образ
FROM scratch

COPY --from=builder /app/bin/rest-service /rest-service

# Запуск приложения
CMD ["./rest-service"]