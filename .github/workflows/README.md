# GitHub Actions Workflows

## docker.yml

Автоматическая сборка и публикация Docker образа в GitHub Container Registry.

**Триггеры:**
- Push в ветки: `main`, `master`, `develop`
- Создание тегов: `v*.*.*` (например, `v1.0.0`)
- Pull Request в `main`/`master` (только сборка, без публикации)

**Что делает:**
1. Собирает Docker образ с использованием Buildx
2. Публикует в `ghcr.io/YOUR_USERNAME/schedluer`
3. Поддерживает multi-platform (amd64, arm64)
4. Использует кэширование для ускорения сборки

**Теги:**
- `latest` - для основной ветки
- `main`, `master`, `develop` - для веток
- `v1.0.0`, `v1.0`, `v1` - для семантических версий
- `main-abc1234` - для коммитов

## docker-manual.yml

Ручной запуск сборки Docker образа через GitHub Actions UI.

**Использование:**
1. Перейдите в Actions → Manual Docker Build
2. Нажмите "Run workflow"
3. Укажите тег и выберите, публиковать ли образ

