# TODO-List

REST API для управления задачами (TODO-лист) на Go с использованием Fiber и PostgreSQL.

## 📋 Содержание

- [Описание](#описание)
- [Требования](#требования)
- [Установка](#установка)
  - [Через Docker](#через-docker)
  - [Локальный запуск](#локальный-запуск)
- [Миграции](#миграции)
- [Переменные окружения](#переменные-окружения)
- [Запуск приложения](#запуск-приложения)
- [API эндпоинты](#api-эндпоинты)
- [Примеры запросов](#примеры-запросов)
- [Очистка и остановка](#очистка-и-остановка)

---

## Описание

Это учебный проект для демонстрации CRUD-операций над задачами (tasks) через HTTP API.

- Язык: Go 1.21+    
- Веб-фреймворк: [Fiber](https://github.com/gofiber/fiber)
- БД: PostgreSQL (pgx)
- Миграции: SQL-скрипты из папки `migrations/`


## Требования

- Docker & Docker Compose
- Go 1.21 и выше (для локального запуска)


## Установка

### Через Docker

1. Скопировать `.env.example` в `.env` и при необходимости изменить значения.
2. Запустить контейнеры:
   ```bash
   docker-compose up --build -d
   ```
3. Проверить логи (не обязательно):
   ```bash
   docker-compose logs -f app
   ```

Контейнер PostgreSQL автоматически инициализирует базу из `migrations/`.

### Локальный запуск

1. Установить зависимости:
   ```bash
   go mod download
   ```
2. Скопировать `.env.example` в `.env` и заполнить параметры доступа к локальной БД.
3. Выполнить миграцию (например, через `psql` или любой мигратор):
   ```bash
   psql "host=$DB_HOST port=$DB_PORT user=$DB_USER password=$DB_PASSWORD dbname=$DB_NAME sslmode=$DB_SSLMODE" -f migrations/0001_create_tasks_table.sql
   ```
4. Запустить приложение:
   ```bash
   go run cmd/server/main.go
   ```

---

## Миграции

SQL-файл `migrations/0001_create_tasks_table.sql` создаёт таблицу:

```sql
CREATE TABLE IF NOT EXISTS tasks (
  id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  status TEXT CHECK (status IN ('new','in_progress','done')) DEFAULT 'new',
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMP DEFAULT now()
);
```

---

## Переменные окружения

```ini
DB_HOST=postgres         # адрес хоста PostgreSQL
DB_PORT=5432             # порт БД
DB_USER=postgres         # пользователь БД
DB_PASSWORD=password     # пароль
DB_NAME=todo             # имя БД
DB_SSLMODE=disable       # sslmode
APP_PORT=8080            # порт приложения
```

---

## Запуск приложения

- Docker: `docker-compose up --build`
- Локально: `go run cmd/server/main.go`

После старта сервер слушает на `http://localhost:${APP_PORT}`.

---

## API эндпоинты

| Метод | URL            | Описание                      |
|-------|----------------|-------------------------------|
| POST  | `/tasks`       | Создать новую задачу          |
| GET   | `/tasks`       | Получить список всех задач    |
| PUT   | `/tasks/:id`   | Обновить задачу по ID         |
| DELETE| `/tasks/:id`   | Удалить задачу по ID          |

---

## Примеры запросов

### Создание задачи

```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy milk",
    "description": "2 liters of milk",
    "status": "new"
  }'
```

### Список задач

```bash
curl http://localhost:8080/tasks
```

### Обновление задачи

```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Buy milk and bread",
    "description": "2 liters of milk, 1 loaf bread",
    "status": "in_progress"
  }'
```

### Удаление задачи

```bash
curl -X DELETE http://localhost:8080/tasks/1
```

---

## Очистка и остановка

- Остановить контейнеры:
  ```bash
  docker-compose down
  ```
- Удалить том с данными PostgreSQL (если нужно):
  ```bash
  docker-compose down -v
  ```

---

