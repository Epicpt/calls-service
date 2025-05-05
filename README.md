# 📰 Calls service
Веб-приложение на Go с использованием Gin для управления заявками колл-центра. 


## 🚀 Запуск проекта

### 1. Клонирование репозитория

```sh
git clone https://github.com/Epicpt/calls-service.git
cd calls-service
```

### 2. Измените .env файл

```sh
mv .env.examples .env
```

### 3. Запуск через Docker Compose

```sh
docker-compose -f docker-compose.yml up -d --build
```

### 4. Сервер доступен:
```sh
http://localhost:8080/
```

### 📡 API Эндпоинты
#### 🔑 Аутентификация

- POST /auth/register – регистрация пользователя
- POST /auth/login – вход (возвращает JWT-токен)

#### 📞 Заявки

- POST /calls – добавление новой заявки (требуется аутентификация)
- GET /calls  – получение списка всех заявок (требуется аутентификация)
- GET /calls/:id - получение информации по конкретной заявке (требуется аутентификация)
- PATCH /calls/:id/status  - изменение статуса заявки (требуется аутентификация)
- DELETE /calls/:id  - удаление заявки (требуется аутентификация)

### 🛠 Используемые технологии

- Golang 1.24.1
- Gin
- gRPC
- PostgreSQL
- Golang-migrate
- Docker & Docker Compose
- JWT