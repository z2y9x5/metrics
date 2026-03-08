COMPILER=go1.24.13
SERVER_PATH=./cmd/server
SERVER_BIN=server

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

test:
	@echo "Форматирование кода:"
	$(COMPILER) fmt ./...
	@echo "Статический анализ:"
	$(COMPILER) vet ./...
	@echo "Тесты приложения:"
	$(COMPILER) test -v -count 1 ./...
	@echo "Покрытие тестами:"
	$(COMPILER) test -count 1 -cover ./...

run: clean test
	$(COMPILER) run $(SERVER_PATH)/main.go

build: clean test
	$(COMPILER) build -o $(SERVER_PATH)/$(SERVER_BIN) $(SERVER_PATH)/*.go
