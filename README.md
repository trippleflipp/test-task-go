# Subscription Service

REST-сервис для агрегации данных об онлайн подписках пользователей.

## Cтек
- **Язык:** Go (Gin Framework)
- **База данных:** PostgreSQL
- **Библиотеки:** sqlx, logrus, swaggo
- **Документация:** Swagger 2.0
- **Контейнеризация:** Docker & Docker Compose

## Основные возможности
- **CRUDL:** Создание, получение (по ID и списком), обновление и удаление подписок.
- **Агрегация:** Ручка для подсчета суммарных затрат пользователя на сервисы за выбранный период.
- **Миграции:** Автоматическое создание таблиц и индексов при старте приложения.
- **Логирование:** Все операции и ошибки логируются в формате JSON.

---

## Требования
Для запуска сервиса необходимы:
- Docker
- Docker Compose

## Быстрый старт

1. **Клонируйте репозиторий:**
   
```bash
   git clone https://github.com/trippleflipp/test-task-go
   cd test-task-go
```

2. **Настройте переменные окружения**
- Создайте файл .env на основе примера:
```bash
cp .env.example .env
```
3. **Запустите проект**
```bash
sudo docker compose up --build
```

Сервис будет доступен по адресу http://localhost:8080.

## Документация API

После запуска проекта документация Swagger доступна по адресу:
http://localhost:8080/swagger/index.html

- Пример создания подписки (POST /api/subscriptions):
```bash
curl -X 'POST' \
  'http://localhost:8080/api/subscriptions' \
  -H 'Content-Type: application/json' \
  -d '{
  "service_name": "Yandex Plus",
  "price": 400,
  "user_id": "60601fee-2bf1-4721-ae6f-7636e79a0cba",
  "start_date": "07-2025"
}'
```
- Пример расчета стоимости (GET /api/subscriptions/total):
```bash
curl 'http://localhost:8080/api/subscriptions/total?user_id=60601fee-2bf1-4721-ae6f-7636e79a0cba&from=01-2025&to=12-2025'
```

