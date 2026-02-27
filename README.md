REST API для управления иерархической структурой отделов и сотрудников компании.

Структуру проекта можно посмотреть в файле "ТЗ Go - API организационной структуры.pdf"

## Технологии
- **Go 1.25**
- **PostgreSQL 16**
- **GORM** — ORM
- **Goose** — миграции БД
- **Docker / Docker Compose**

## Запуск

### 1. Клонировать репозиторий

```bash
git clone https://github.com/gtimofej0303/TZ-hitalent
cd TZ-hitalent
```

### 2. Создать .env файл
```bash
cp .env.example .env
```

### 3. Заполнить .env своими значениями:
```
DB_HOST=postgres
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_db
SERVER_PORT=8080
```

### 3. Запустить через Docker Compose
```Bash
docker compose up --build
```

Сервер поднимется на http://localhost:8080.
Миграции применяются автоматически при старте.

### Примеры запросов:
```
==Создание подразделения==
curl -X POST http://localhost:8080/departments/ \
    -H "Content-Type: application/json" \
    -d '{"name": "Engineering", "parent_id": null}'

==Создание сотрудника==
curl -X POST http://localhost:8080/departments/1/employees/ \
  -H "Content-Type: application/json" \
  -d '{
    "fullname": "Ivan Petrov",
    "position": "Developer",
    "hired_at": "2024-01-15T00:00:00Z"
  }'

==Обновление подразделения==
curl -X PATCH http://localhost:8080/departments/2 \
  -H "Content-Type: application/json" \
  -d '{"name": "Math", "parent_id": 1}'

==Получение информации о подразделении==
curl "http://localhost:8080/departments/1?depth=2&include_employees=true"

==Удаление (cascade)==
curl -X DELETE "http://localhost:8080/departments/1?mode=cascade"

==Удаление (reassign)==
curl -X DELETE "http://localhost:8080/departments/2?mode=reassign&reassign_to=1"
```