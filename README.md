Простой сервер для управления задачами с поддержкой повторяющихся событий и JWT аутентификацией.

### Функциональность

- Создание, редактирование, удаление задач
- Поддержка повторяющихся задач (ежедневно, еженедельно, ежемесячно, ежегодно)
- JWT аутентификация
- Поиск задач по тексту и дате
- RESTful API

### Список выполненных заданий со звёздочкой:
- Возможность определять путь к файлу базы данных через переменную окружения;
- Вычисление следующих дат в случае недель и месяцев;
- Реализован поиск задач по тексту заголовка/комментария или дате;

### Локальный запуск
1. Установите зависимости
    go mod download
2. Запустите сервер с переменными окружения:
    ```bash
    export TODO_PASSWORD=mysecretpassword
    export TODO_DBFILE=./scheduler.db
    export JWT_SECRET=mysupersecretpassword
    export TODO_PORT=7540
    export TOKEN_DURATION=8h

    go run main.go
    ```
    
### Пример .env файла
```bash
TODO_PASSWORD=mysecretpassword
TODO_DBFILE=./scheduler.db
JWT_SECRET=mysupersecretpassword
TODO_PORT=7540
TOKEN_DURATION=8h
```

### Доступ в браузере
Откройте: http://localhost:7540

API endpoints:
- GET /api/tasks - список задач
- POST /api/task - создать задачу
- PUT /api/task - обновить задачу
- DELETE /api/task - удалить задачу
- POST /api/signin - аутентификация
- POST /api/task/done - отметить задачу выполненной
- GET /api/nextdate - рассчитать следующую дату

### Примеры запросов
## Создание задачи:
```bash
curl -X POST http://localhost:7540/api/task \
-H "Content-Type: application/json" \
-d '{
    "date": "20231225",
    "title": "Сделать домашку",
    "comment": "По математике", 
    "repeat": "d 1"
}'
```
## Аутентификация:
```bash
curl -X POST http://localhost:7540/api/signin \
-H "Content-Type: application/json" \
-d '{"password": "mysecretpassword"}'
```

### Запуск тестов
## Всех тестов:
```bash
go test ./tests
```

## Конкретного теста:
```bash
go test -run ^наименование_теста$ ./tests
```

Для запуска тестов рекомендуется использовать следующие значения файла tests/settings.go
    Port = 7540
    DBFile = "../scheduler.db"
    FullNextDate = true
    Search = true
    Token = "token"

## Сборка и запуск проекта через докер
# Сборка и запуск
```bash
docker-compose up -d
```
# Только сборка
```bash
docker build -t todo-server .
```
# Запуск контейнера
```bash
docker run -p 7540:7540 \
  -e TODO_PASSWORD=mysecretpassword \
  -e JWT_SECRET=mysupersecretpassword \
  -e TODO_DBFILE=/app/data/scheduler.db \
  -e TODO_PORT=7540 \
  -v todo-data:/app/data \
  todo-server
```