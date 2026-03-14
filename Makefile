COMPILER=go1.24.13
SERVER_PATH=./cmd/server
SERVER_BIN=server
AGENT_PATH=./cmd/agent
AGENT_BIN=agent

help:
	@echo "Доступные команды:"
	@echo "    make clean"
	@echo "        Удаление сгенерированных файлов."
	@echo "    make test"
	@echo "        Тестирование приложения."
	@echo "    make run"
	@echo "        Запуск приложения."
	@echo "    make build"
	@echo "        Сборка приложения."

clean:
	rm -f $(SERVER_PATH)/$(SERVER_BIN)
	rm -f $(AGENT_PATH)/$(AGENT_BIN)
	rm -f coverage.out

test:
	@echo "Форматирование кода:"
	$(COMPILER) fmt ./...
	@echo "Статический анализ:"
	$(COMPILER) vet ./...
	@echo "Тесты приложения:"
	$(COMPILER) test -v -count 1 ./...
	@echo "Покрытие тестами:"
	$(COMPILER) test -coverprofile=coverage.out ./...
	$(COMPILER) tool cover -func=coverage.out

run: clean test
	$(COMPILER) run $(SERVER_PATH)/main.go
	$(COMPILER) run $(AGENT_PATH)/main.go

build: clean test
	$(COMPILER) build -o $(SERVER_PATH)/$(SERVER_BIN) $(SERVER_PATH)/*.go
	$(COMPILER) build -o $(AGENT_PATH)/$(AGENT_BIN) $(AGENT_PATH)/*.go
