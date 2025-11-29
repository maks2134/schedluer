.PHONY: build run test clean swagger docker-build docker-run help

# Переменные
APP_NAME=schedluer
DOCKER_IMAGE=schedluer:latest
DOCKER_CONTAINER=schedluer

help: ## Показать справку
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'

swagger: ## Сгенерировать Swagger документацию
	@echo "Генерация Swagger документации..."
	@swag init -g cmd/schedluer/main.go -o docs
	@echo "Swagger документация сгенерирована в docs/"

build: ## Собрать приложение
	@echo "Сборка приложения..."
	@go build -o $(APP_NAME) ./cmd/schedluer
	@echo "Приложение собрано: $(APP_NAME)"

run: ## Запустить приложение
	@echo "Запуск приложения..."
	@go run ./cmd/schedluer/main.go

test: ## Запустить тесты
	@echo "Запуск тестов..."
	@go test -v ./...

clean: ## Очистить артефакты сборки
	@echo "Очистка..."
	@rm -f $(APP_NAME)
	@rm -rf docs/
	@go clean

docker-build: ## Собрать Docker образ
	@echo "Сборка Docker образа..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker образ собран: $(DOCKER_IMAGE)"

docker-run: ## Запустить Docker контейнер
	@echo "Запуск Docker контейнера..."
	@docker run --rm -p 8080:8080 --env-file .env --name $(DOCKER_CONTAINER) $(DOCKER_IMAGE)

docker-stop: ## Остановить Docker контейнер
	@echo "Остановка Docker контейнера..."
	@docker stop $(DOCKER_CONTAINER) || true
	@docker rm $(DOCKER_CONTAINER) || true

docker-compose-up: ## Запустить через docker-compose
	@echo "Запуск через docker-compose..."
	@docker-compose up -d

docker-compose-down: ## Остановить docker-compose
	@echo "Остановка docker-compose..."
	@docker-compose down

docker-compose-logs: ## Показать логи docker-compose
	@docker-compose logs -f

deps: ## Установить зависимости
	@echo "Установка зависимостей..."
	@go mod download
	@go mod tidy

install-tools: ## Установить инструменты разработки
	@echo "Установка инструментов..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "Инструменты установлены"

dev: swagger run ## Запустить в режиме разработки (с генерацией swagger)

