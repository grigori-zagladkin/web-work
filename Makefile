# Путь к Go-сервису относительно корня репозитория
API_DIR := services/api
# Цель по умолчанию: просто `make` покажет список команд
.DEFAULT_GOAL := help
# .PHONY говорит make: это названия команд, а не файлы.
# Без этого, если в папке появится файл с именем "test",
# `make test` решит, что всё уже сделано, и ничего не запустит.
.PHONY: help up down logs ps db run test build tidy migrate admin
help: ## Показать список команд
	@echo Команды:
	@echo   make up       - собрать и запустить всё в Docker
	@echo   make down     - остановить контейнеры
	@echo   make logs     - логи Go-сервера
	@echo   make ps       - статусы контейнеров
	@echo   make db       - поднять только Postgres
	@echo   make run      - запустить API локально
	@echo   make migrate  - применить миграции локально
	@echo   make test     - запустить тесты
	@echo   make build    - проверить компиляцию
	@echo   make tidy     - обновить зависимости
	@echo   make fmt      - форматировать код
	@echo   make admin EMAIL=... PASSWORD=...  - создать админа
# ---------- Docker ----------
up: ## Собрать и запустить всё в Docker
	docker compose up -d --build
down: ## Остановить контейнеры (данные БД сохранятся)
	docker compose down
logs: ## Логи Go-сервера в реальном времени
	docker compose logs -f api
ps: ## Статусы контейнеров
	docker compose ps
db: ## Поднять только Postgres (для локального go run)
	docker compose up -d db
# ---------- Go локально ----------
run: ## Запустить API локально
	cd $(API_DIR) && go run ./cmd/api
migrate: ## Применить миграции локально
	cd $(API_DIR) && go run ./cmd/api migrate
test: ## Запустить тесты
	cd $(API_DIR) && go test ./...
build: ## Проверить, что всё компилируется
	cd $(API_DIR) && go build ./...
tidy: ## Обновить зависимости (go.mod / go.sum)
	cd $(API_DIR) && go mod tidy
fmt: ## Обновить зависимости (go.mod / go.sum)
	cd $(API_DIR) && go fmt ./...
# ---------- Админ ----------
# Использование: make admin EMAIL=you@example.com PASSWORD=очень-длинный-пароль
admin: ## Создать админа в Docker (EMAIL=... PASSWORD=...)
	docker compose run --rm api create-admin -email $(EMAIL) -password $(PASSWORD)