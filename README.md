# Сервис информации об угрозах

REST API на Go для хранения информации об угрозах. Запись можно создать, получить по MD5 файла и обновить. Данные хранятся в PostgreSQL.

Сам сервис файлы не анализирует и MD5 не вычисляет — клиент передаёт готовый хеш и сведения об угрозе.

Спецификация API: [api/openapi.yaml](api/openapi.yaml). К исходной модели добавлены описание и уровень опасности.

## Запуск

Нужны Go 1.26.2 и Docker с Docker Compose. Команды выполняются из корня проекта.

Поднять PostgreSQL:

```bash
docker compose up -d db
```

Проверить, что база готова принимать подключения:

```bash
docker compose exec db pg_isready -U app -d threats
```

Когда появится `accepting connections`, можно запускать приложение:

```bash
export CONNECTION='postgres://app:local_dev_password@localhost:5433/threats?sslmode=disable'
go run ./server
```

Сервер слушает `http://localhost:9091`. Миграции применяются при запуске приложения, отдельно устанавливать `migrate` не нужно.

В Docker запускается база, Go-приложение запускается локально. Параметры базы заданы в `compose.yaml`. Данные лежат в Docker volume и сохраняются после перезапуска приложения и контейнера.

## Запросы

| Метод | Адрес | Действие | Успешный ответ |
|---|---|---|---|
| POST | `/info` | Создать запись | `201` и созданная запись |
| GET | `/info/{md5}` | Получить запись | `200` и найденная запись |
| PUT | `/info/{md5}` | Обновить запись | `200` и обновлённая запись |

Поля записи:

- `md5` — 32 шестнадцатеричных символа, уникальный ключ записи.
- `name` — непустое название угрозы.
- `description` — необязательное описание.
- `severity` — одно из значений: `safe`, `low`, `medium`, `high`.

При создании обязательны `md5`, `name` и `severity`. При обновлении MD5 берётся из адреса, а в JSON обязательны `name` и `severity`.

### Создать

```bash
curl -i -X POST http://localhost:9091/info \
  -H 'Content-Type: application/json' \
  -d '{
    "md5": "0123456789abcdef0123456789abcdef",
    "name": "Trojan.Example",
    "description": "Тестовая запись",
    "severity": "high"
  }'
```

При повторном создании записи с тем же MD5 сервер вернёт `409`.

### Получить

```bash
curl -i http://localhost:9091/info/0123456789abcdef0123456789abcdef
```

### Обновить

```bash
curl -i -X PUT http://localhost:9091/info/0123456789abcdef0123456789abcdef \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "Trojan.Example.Updated",
    "description": "Уточнённое описание",
    "severity": "medium"
  }'
```

PUT заменяет редактируемые поля целиком. MD5 остаётся прежним. Если не передать `description`, старое описание очистится.

Для некорректных данных сервер возвращает `400`, для отсутствующей записи — `404`, при внутренней ошибке — `500`. Ошибки возвращаются в JSON с полями `message` и `time`.

## Тесты

```bash
go test ./...
```

Пока есть тест проверки допустимых значений `severity`.
