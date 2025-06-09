# Go Translate Service

Сервис для изучения иностранных слов с поддержкой MongoDB и Redis кэширования.

## Технологии

- Go 1.22+
- MongoDB 7.0
- Redis 7.0
- Docker & Docker Compose
- JWT для авторизации
- Swagger для API документации

## Быстрый старт

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd go_translate
```

2. Запустите сервисы через Docker Compose:
```bash
docker-compose up -d --build
```

3. Проверьте статус сервисов:
- API: http://localhost:8080/swagger/index.html
- Mongo Express: http://localhost:8081
  - Логин: admin
  - Пароль: strongpassword

## API Endpoints

### Аутентификация

1. Регистрация нового пользователя:
```bash
curl -X POST http://localhost:8080/api/auth/sign-up \
  -H "Content-Type: application/json" \
  -d '{"name":"User Name","username":"user","password":"password"}'
```

2. Получение токена:
```bash
curl -X POST http://localhost:8080/api/auth/sign-in \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"password"}'
```

### Работа с переводами

1. Добавление нового перевода:
```bash
curl -X POST http://localhost:8080/api/translations \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"phrase":"hello","expected_translation":"привет"}'
```

2. Получение всех переводов:
```bash
curl -X GET http://localhost:8080/api/translations \
  -H "Authorization: Bearer <TOKEN>"
```

### Тестирование производительности

1. Добавить 100 тестовых переводов:
```bash
curl -X POST http://localhost:8080/api/translations/post100 \
  -H "Authorization: Bearer <TOKEN>"
```

2. Добавить 100 000 тестовых переводов:
```bash
curl -X POST http://localhost:8080/api/translations/post100k \
  -H "Authorization: Bearer <TOKEN>"
```

3. Удалить все переводы пользователя:
```bash
curl -X DELETE http://localhost:8080/api/translations/delete_all \
  -H "Authorization: Bearer <TOKEN>"
```

## Архитектура

Проект следует чистой архитектуре и разделен на слои:
- `cmd/` - точка входа приложения
- `pkg/handler/` - обработчики HTTP запросов
- `pkg/service/` - бизнес-логика
- `pkg/repository/` - работа с БД
- `pkg/tables/` - модели данных

## Кэширование

Redis используется для кэширования:
- Списка переводов пользователя
- Отдельных переводов по ID
Время жизни кэша: 60 секунд

## Разработка

1. Установите зависимости:
```bash
go mod download
```

2. Запустите локально:
```bash
go run cmd/main.go
```

## Остановка сервисов

```bash
docker-compose down
```

## Лицензия

MIT
