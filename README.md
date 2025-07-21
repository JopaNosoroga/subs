# Сервис Управления Подписками

RESTful API для управления подписками пользователей. Позволяет создавать, удалять, обновлять и просматривать подписки, а также рассчитывать их стоимость за указанный период.

## Технологический стек
- **Язык**: Go 1.24.4
- **База данных**: PostgreSQL 15
- **Роутер**: Gorilla Mux
- **Миграции**: `golang-migrate`
- **Контейнеризация**: Docker
- **Документация**: Swagger (OpenAPI 2.0)

## Запуск проекта
1. Склонируйте репозиторий:
```bash
git clone <URL_репозитория>
cd <директория_проекта>
```

2. Создайте файл config.env в корне проекта

```env
# Настройки подключения к БД
USER=postgres
PASSWORD=postgres
HOST=db
PORT=5432
SSL_MODE=disable
DBNAME=subs
```

3. Запустите проект
```bash
docker-compose up --build --force-recreate
```
Приложение будет доступно по адрессу http://localhost:8080

EndPoints
## Работа проекта
POST	/subscriptions	Создать подписку
GET	/subscriptions	Получить подписки (с фильтрами)
GET	/subscriptions/{id}	Получить подписку по ID
PUT	/subscriptions/{id}	Обновить подписку
DELETE	/subscriptions/{id}	Удалить подписку
POST	/sum	Рассчитать стоимость подписок

-Примеры запросов
```bash
curl -X POST http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "123e4567-e89b-12d3-a456-426614174000",
    "service_name": "Netflix",
    "price": 1000,
    "start_date": "2023-01-15T00:00:00Z",
    "end_date": "2023-02-15T00:00:00Z"
  }'
```

```bash
curl -X GET http://localhost:8080/subscriptions \
  -H "Content-Type: application/json" \
  -d '{"user_id": "123e4567-e89b-12d3-a456-426614174000"}'
```

```bash
curl -X POST http://localhost:8080/sum \
  -H "Content-Type: application/json" \
  -d '{
    "start_date": "2023-01-01T00:00:00Z",
    "end_date": "2023-01-31T00:00:00Z"
  }'
```
## Структура проекта
```text
├── docker-compose.yml
├── go.mod
├── main.go
├── .env (создается пользователем)
├── pkg
│   ├── dbwork
│   │   ├── dbwork.go
│   │   └── migrations
│   │       ├── 1_subscriptions.down.sql
│   │       └── 1_subscriptions.up.sql
│   ├── handlers
│   │   └── handlers.go
│   └── models
│       └── models.go
└── swagger.yaml
```
