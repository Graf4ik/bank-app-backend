# Определяем переменные
CONFIG_PATH="C:\Users\user\Desktop\bank-app-backend\config\local.yaml"
GO=go

.PHONY: build run test clean

# Сборка проекта
build:
	$(GO) build -o $(BINARY_NAME) ./cmd/main.go

# Запуск
run:
	$(GO) run ./cmd/main.go

# Docker-compose
docker-compose:
	$(DOCKER) up -d $(APP_NAME):$(VERSION) .

# Тестирование
test:
	$(GO) test -v ./...

# Очистка
clean:
	rm -f $(BINARY_NAME)
	$(GO) clean