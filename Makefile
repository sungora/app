# Инициализация
# DIR_ROOT := $(realpath -m ..)
# APP := $(shell basename $(DIR))
DIR := $(realpath -m .)
APP := github.com/sungora/app
PROJECT_LIST = sample

default: help

# Зависимости
dep:
	@if [ ! -f $(DIR)/go.mod ]; then \
		cd $(DIR); go mod init $(APP); \
	fi
	cd $(DIR); go mod tidy;
	cd $(DIR); go mod vendor;
.PHONY: dep

# Сборка
build:
	@for dir in $(PROJECT_LIST); do \
		cd $(DIR)/cmd/$${dir}; go build -i -mod vendor -o $(DIR)/bin/$${dir}; \
	done
.PHONY: build

# Запуск
sample:
	cd $(DIR)/cmd/sample; go build -i -mod vendor -o $(DIR)/bin/sample;
	$(DIR)/bin/sample;
.PHONY: run

# Help
help:
	@echo "Usage: make [target]"
	@echo "  target is:"
	@echo "    dep		- Загрузка и обновление зависимостей"
	@echo "    build	- Компиляция указанных в настройке приложений"
	@echo "    sample	- Запуск приложения sample"
.PHONY: help
