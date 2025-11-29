# Schedluer

Современное веб-приложение для просмотра расписания занятий и экзаменов БГУИР (Белорусский государственный университет информатики и радиоэлектроники).

## Интерфейс 🎨

- Главная страница: 
![main page](docs/image/main_page.png)
- Расписание группы
![schedule](docs/image/groups_shedule.png)
- Расписание преподавателей
![schedule](docs/image/prepod_shedule.png)

## 🚀 Особенности

- 📅 Просмотр расписания занятий и экзаменов для групп и преподавателей
- ⭐ Система избранных групп с быстрым доступом
- 🔍 Полнотекстовый поиск по избранным группам (MongoDB Atlas Search)
- 📱 Адаптивный дизайн для всех устройств
- 🎨 Современный UI на базе Shadcn/ui и Tailwind CSS
- 📊 Разделение расписания на пары и экзамены
- 🔄 Автоматическое кэширование данных
- 📝 Полная API документация (Swagger)

## 🛠 Технологии

### Backend
- **Go 1.25** - Язык программирования
- **Gin** - HTTP веб-фреймворк
- **MongoDB** - База данных
- **MongoDB Atlas Search** - Полнотекстовый поиск
- **Swagger** - API документация
- **Docker** - Контейнеризация

### Frontend
- **Next.js 16** - React фреймворк с App Router
- **TypeScript** - Типизированный JavaScript
- **Tailwind CSS** - Utility-first CSS фреймворк
- **Shadcn/ui** - Высококачественные UI компоненты
- **Lucide React** - Иконки

## 📁 Структура проекта

```
schedluer/
├── server/                 # Backend приложение
│   ├── cmd/
│   │   └── schedluer/      # Точка входа приложения
│   ├── internal/
│   │   ├── config/         # Конфигурация
│   │   ├── container/      # Dependency Injection
│   │   ├── handler/        # HTTP обработчики
│   │   ├── models/         # Модели данных
│   │   ├── repository/     # Слой доступа к данным
│   │   └── service/        # Бизнес-логика
│   ├── pkg/
│   │   ├── bsuir/          # Клиент для API БГУИР
│   │   ├── converter/       # Конвертеры данных
│   │   └── database/       # Подключение к MongoDB
│   ├── docs/               # Swagger документация
│   ├── Dockerfile          # Docker образ
│   ├── docker-compose.yml  # Docker Compose для разработки
│   └── Makefile            # Команды для разработки
│
└── client/                 # Frontend приложение
    ├── app/                # Next.js App Router
    ├── components/         # React компоненты
    ├── lib/                # Утилиты и хуки
    ├── types/              # TypeScript типы
    └── public/             # Статические файлы
```

## 🚦 Быстрый старт

### Предварительные требования

- Go 1.25 или выше
- Node.js 20 или выше
- MongoDB (локально или MongoDB Atlas)
- Docker (опционально, для контейнеризации)

### Установка

1. **Клонируйте репозиторий:**
```bash
git clone <repository-url>
cd schedluer
```

2. **Настройте Backend:**

```bash
cd server

# Создайте файл .env в папке server/
cat > .env << EOF
SERVER_PORT=8080
SERVER_HOST=localhost
MONGODB_URI=mongodb://localhost:27017/schedluer
# или для MongoDB Atlas:
# MONGODB_URI=mongodb+srv://user:password@cluster.mongodb.net/schedluer?retryWrites=true&w=majority
MONGODB_DATABASE=schedluer
BSUIR_API_BASE_URL=https://iis.bsuir.by/api/v1
LOG_LEVEL=info
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
EOF

# Установите зависимости
go mod download

# Сгенерируйте Swagger документацию
make swagger
# или
swag init -g cmd/schedluer/main.go -o docs
```

3. **Настройте Frontend:**

```bash
cd ../client

# Установите зависимости
npm install

# Создайте файл .env.local
cat > .env.local << EOF
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
EOF
```

### Запуск

#### Backend

```bash
cd server

# Режим разработки
make dev
# или
go run cmd/schedluer/main.go

# Production сборка
make build
./schedluer
```

Backend будет доступен по адресу: http://localhost:8080

#### Frontend

```bash
cd client

# Режим разработки
npm run dev

# Production сборка
npm run build
npm start
```

Frontend будет доступен по адресу: http://localhost:3000

## 🐳 Docker

### Разработка

```bash
cd server
docker-compose up -d
```

### Production

```bash
cd server

# Соберите образ
make docker-build
# или
docker build -t schedluer:latest .

# Запустите контейнер
docker run -d \
  --name schedluer \
  -p 8080:8080 \
  -e MONGODB_URI=your_mongodb_uri \
  -e MONGODB_DATABASE=schedluer \
  schedluer:latest
```

## 📚 API Документация

После запуска backend, Swagger UI доступен по адресу:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health

### Основные эндпоинты

#### Расписание
- `GET /api/v1/schedule/group/:groupNumber` - Получить расписание группы
- `GET /api/v1/schedule/employee/:urlId` - Получить расписание преподавателя
- `POST /api/v1/schedule/group/:groupNumber/refresh` - Обновить расписание группы
- `POST /api/v1/schedule/employee/:urlId/refresh` - Обновить расписание преподавателя

#### Группы
- `GET /api/v1/groups` - Список всех групп
- `GET /api/v1/groups/:groupNumber` - Получить группу по номеру
- `POST /api/v1/groups/refresh` - Обновить список групп

#### Преподаватели
- `GET /api/v1/employees` - Список всех преподавателей
- `GET /api/v1/employees/:urlId` - Получить преподавателя по URL ID
- `POST /api/v1/employees/refresh` - Обновить список преподавателей

#### Избранные группы
- `GET /api/v1/favorites?user_id=default` - Получить все избранные группы
- `GET /api/v1/favorites/search?user_id=default&query=350` - Поиск по избранным группам
- `POST /api/v1/favorites/:groupNumber?user_id=default` - Добавить группу в избранное
- `DELETE /api/v1/favorites/:groupNumber?user_id=default` - Удалить группу из избранного
- `GET /api/v1/favorites/:groupNumber/check?user_id=default` - Проверить, является ли группа избранной

## 🔍 MongoDB Atlas Search Index

Для использования полнотекстового поиска по избранным группам необходимо создать Search Index в MongoDB Atlas.

### Создание Search Index

1. Откройте MongoDB Atlas Dashboard
2. Перейдите в ваш кластер → вкладка "Search"
3. Нажмите "Create Search Index" → выберите "JSON Editor"
4. Вставьте следующую конфигурацию:

```json
{
  "name": "default",
  "type": "search",
  "definition": {
    "mappings": {
      "dynamic": false,
      "fields": {
        "group_number": {
          "type": "string",
          "analyzer": "lucene.standard",
          "searchAnalyzer": "lucene.standard"
        },
        "user_id": {
          "type": "string"
        }
      }
    }
  }
}
```

5. Выберите базу данных: `schedluer`
6. Выберите коллекцию: `favorite_groups`
7. Нажмите "Create Search Index"

После создания индекса, API endpoint `/api/v1/favorites/search` будет использовать полнотекстовый поиск. Если индекс не создан, система автоматически переключится на regex поиск.

## 🛠 Разработка

### Backend команды

```bash
cd server

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

### Frontend команды

```bash
cd client

npm run dev        # Запустить в режиме разработки
npm run build      # Собрать production версию
npm start          # Запустить production версию
npm run lint       # Проверить код линтером
```

### Переменные окружения

#### Backend (.env)

```env
SERVER_PORT=8080
SERVER_HOST=localhost
MONGODB_URI=mongodb://localhost:27017/schedluer
MONGODB_DATABASE=schedluer
BSUIR_API_BASE_URL=https://iis.bsuir.by/api/v1
LOG_LEVEL=info
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
```

#### Frontend (.env.local)

```env
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

## 📝 Особенности реализации

### Backend

- **Архитектура**: Clean Architecture с разделением на слои (handler → service → repository)
- **Кэширование**: Автоматическое кэширование данных групп и преподавателей в MongoDB
- **Обработка ошибок**: Централизованная обработка ошибок с логированием
- **CORS**: Настраиваемые разрешенные источники для фронтенда
- **Health Check**: Эндпоинт для проверки состояния приложения

### Frontend

- **Оптимистичные обновления**: Мгновенное обновление UI при добавлении/удалении избранных групп
- **Адаптивный дизайн**: Полная поддержка мобильных устройств
- **Фильтрация**: Фильтр по подгруппам для расписания
- **Дедупликация**: Автоматическое удаление дубликатов в расписании
- **Табы**: Разделение расписания на пары и экзамены
