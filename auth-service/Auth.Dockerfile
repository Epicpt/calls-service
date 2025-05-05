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
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -o /app/bin/auth-service ./auth-service/cmd/main.go

FROM scratch

COPY --from=builder /app/bin/auth-service /auth-service

# Запуск приложения
CMD ["./auth-service"]