# Docker инструкции

## Быстрый старт

### 1. Подготовка

Убедитесь, что у вас есть файл `.env` с правильными настройками:
```bash
cp .env.example .env
# Отредактируйте .env и укажите реальный пароль MongoDB
```

### 2. Сборка образа

```bash
docker build -t schedluer:latest .
```

Или используйте Makefile:
```bash
make docker-build
```

### 3. Запуск контейнера

#### Простой запуск:
```bash
docker run --rm -p 8080:8080 --env-file .env schedluer:latest
```

#### С именем контейнера:
```bash
docker run --rm -p 8080:8080 --env-file .env --name schedluer schedluer:latest
```

#### В фоновом режиме:
```bash
docker run -d -p 8080:8080 --env-file .env --name schedluer --restart unless-stopped schedluer:latest
```

### 4. Docker Compose

```bash
docker-compose up -d
```

Просмотр логов:
```bash
docker-compose logs -f
```

Остановка:
```bash
docker-compose down
```

## Особенности Dockerfile

### Multi-stage build
- **Stage 1 (builder)**: Сборка приложения с полным Go окружением
- **Stage 2 (final)**: Минимальный Alpine образ только с бинарником

### Оптимизации
- Кэширование слоев зависимостей (go.mod/go.sum)
- Статическая сборка (CGO_ENABLED=0)
- Минимальный размер образа (~20MB)
- Непривилегированный пользователь для безопасности

### Health Check
Контейнер автоматически проверяет здоровье через `/health` endpoint каждые 30 секунд.

## Переменные окружения

Все переменные из `.env` файла передаются в контейнер. Убедитесь, что:
- `MONGODB_URI` содержит правильный пароль
- `SERVER_HOST=0.0.0.0` для доступа извне контейнера
- `SERVER_PORT=8080` соответствует EXPOSE в Dockerfile

## Публикация образа

### Тегирование
```bash
docker tag schedluer:latest your-registry/schedluer:v1.0.0
docker tag schedluer:latest your-registry/schedluer:latest
```

### Публикация
```bash
docker push your-registry/schedluer:v1.0.0
docker push your-registry/schedluer:latest
```

## Troubleshooting

### Проверка логов
```bash
docker logs schedluer
docker logs -f schedluer  # следовать за логами
```

### Проверка здоровья
```bash
curl http://localhost:8080/health
```

### Вход в контейнер
```bash
docker exec -it schedluer sh
```

### Пересборка без кэша
```bash
docker build --no-cache -t schedluer:latest .
```

