# Schedluer Backend

Backend сервис для работы с расписанием БГУИРа на Go 1.25 + MongoDB.

## Структура проекта

```
server/
├── cmd/
│   └── schedluer/          # Точка входа приложения
│       └── main.go
├── internal/                # Внутренние пакеты приложения
│   ├── config/             # Конфигурация приложения
│   ├── handler/            # HTTP handlers
│   ├── service/            # Бизнес-логика
│   ├── repository/         # Работа с базой данных
│   └── models/             # Модели данных
├── pkg/                     # Переиспользуемые пакеты
│   ├── bsuir/              # Клиент для API БГУИРа
│   ├── database/           # MongoDB обертка
│   └── converter/          # Конвертация расписания
├── docs/                    # Swagger документация (генерируется)
├── Dockerfile              # Multi-stage Dockerfile
├── docker-compose.yml      # Docker Compose конфигурация
├── Makefile                # Команды для разработки
├── go.mod
└── README.md
```

## Зависимости

- **Gin** - HTTP веб-фреймворк
- **MongoDB Driver v2** - Драйвер для работы с MongoDB
- **Logrus** - Структурированное логирование
- **Godotenv** - Загрузка переменных окружения из .env файла
- **Swagger** - API документация

## Установка и запуск

### Локальная разработка

1. Установите зависимости:
```bash
make deps
```

2. Установите инструменты разработки:
```bash
make install-tools
```

3. Создайте файл `.env`:
```bash
cp .env.example .env
```

4. Настройте переменные окружения в `.env`:
```env
SERVER_PORT=8080
SERVER_HOST=localhost
MONGODB_URI=mongodb+srv://user:password@cluster.mongodb.net/?appName=Cluster0
MONGODB_DATABASE=schedluer
BSUIR_API_BASE_URL=https://iis.bsuir.by/api/v1
LOG_LEVEL=info
```

5. Сгенерируйте Swagger документацию:
```bash
make swagger
```

6. Запустите приложение:
```bash
make run
# или
go run cmd/schedluer/main.go
```

### Docker

#### Сборка образа:
```bash
make docker-build
# или
docker build -t schedluer:latest .
```

#### Запуск контейнера:
```bash
make docker-run
# или
docker run --rm -p 8080:8080 --env-file .env schedluer:latest
```

#### Docker Compose:
```bash
make docker-compose-up
# или
docker-compose up -d
```

Просмотр логов:
```bash
make docker-compose-logs
```

Остановка:
```bash
make docker-compose-down
```

## API Документация

После запуска приложения Swagger UI доступен по адресу:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health

## Доступные команды

Используйте `make help` для просмотра всех доступных команд:

```bash
make help          # Показать справку
make swagger       # Сгенерировать Swagger документацию
make build         # Собрать приложение
make run           # Запустить приложение
make test          # Запустить тесты
make clean         # Очистить артефакты
make docker-build  # Собрать Docker образ
make docker-run    # Запустить Docker контейнер
make dev           # Запустить в режиме разработки
```

## API Endpoints

### Расписание
- `GET /api/v1/schedule/group/:groupNumber` - Получить расписание группы
- `GET /api/v1/schedule/employee/:urlId` - Получить расписание преподавателя
- `POST /api/v1/schedule/group/:groupNumber/refresh` - Обновить расписание группы
- `POST /api/v1/schedule/employee/:urlId/refresh` - Обновить расписание преподавателя

### Группы
- `GET /api/v1/groups` - Список всех групп
- `GET /api/v1/groups/:groupNumber` - Получить группу по номеру
- `POST /api/v1/groups/refresh` - Обновить список групп

### Преподаватели
- `GET /api/v1/employees` - Список всех преподавателей
- `GET /api/v1/employees/:urlId` - Получить преподавателя по URL ID
- `POST /api/v1/employees/refresh` - Обновить список преподавателей

## Docker Best Practices

Dockerfile использует multi-stage build для:
- Минимального размера финального образа (~20MB)
- Безопасности (непривилегированный пользователь)
- Оптимизации кэширования слоев
- Статической сборки бинарника

## Разработка

1. Внесите изменения в код
2. Обновите Swagger аннотации в handlers
3. Сгенерируйте документацию: `make swagger`
4. Протестируйте локально: `make run`
5. Соберите Docker образ: `make docker-build`
6. Протестируйте в контейнере: `make docker-run`

## Лицензия

Apache 2.0
