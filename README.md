# Файлы для итогового задания

В директории `tests` находятся тесты для проверки API, которое должно быть реализовано в веб-сервере.

Директория `web` содержит файлы фронтенда.

# Сборка и запуск:
    # Собрать образ
    docker build -t todo-server .

    # Запустить контейнер
    docker run -p 7540:7540 \
    -e TODO_PASSWORD="mysecretpassword" \
    -v todo-data:/app/data \
    todo-server

    # Или с docker-compose
    docker-compose up -d