## Инструментарий проектов webnice - единый Makefile.
## (c)Alex Geer <monoflash@gmail.com>
## Версия: 12.11.2026.

## Вычисление текущей директории проекта.
export DIR := $(strip $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST)))))

## Подключение переменных окружения проекта.
ifneq ("$(wildcard $(DIR)/.env.mk)","")
include $(DIR)/.env.mk
endif
## Подключение makefile части, содержащей специфичные для проекта команды.
ifneq ("$(wildcard $(DIR)/.prj.mk)","")
include $(DIR)/.prj.mk
endif

## Тип проекта: go - репозиторий golang, web - frontend репозиторий.
PROJECT_TYPE := $(shell if [ -f "$(DIR)/go.mod" ]; then echo "go"; elif [ -f "$(DIR)/package.json" ]; then echo "web"; fi)

GOPATH := $(GOPATH)
DATE   := $(shell date -u +%Y%m%d.%H%M%S.%Z)

####################################
## Репозитории golang.
ifeq ($(PROJECT_TYPE),go)

## Имя приложения по умолчанию - название корневой директории проекта.
APP ?= $(shell basename "$(DIR)")
## Имя go модуля проекта.
APP_PACKAGE := $(shell grep -m1 -E '^module[[:space:]]+' "$(DIR)/go.mod" 2>/dev/null | awk '{ print $$2 }')
## Директория размещения бинарных файлов.
export GOBIN := $(DIR)/bin
## Флаги для создания бинарника (производственный режим).
export LDFLAGS_MODE := $(shell if [ -n "$$GO_SERVER_URL" ]; then echo '-s -w'; fi)
## Флаги сборки с простановкой версии и даты сборки.
LDFLAGS := -X main.build=$(DATE) $(LDFLAGS_MODE)

## Пакеты кодогенерации (.gogenerate).
GOGENERATE = $(shell if [ -f "$(DIR)/.gogenerate" ]; then cat "$(DIR)/.gogenerate"; fi)
## Тестируемые пакеты (.testpackages).
TESTPACKETS = $(shell if [ -f "$(DIR)/.testpackages" ]; then cat "$(DIR)/.testpackages"; fi)
## Пакеты тестирования производительности (.benchpackages).
BENCHPACKETS = $(shell if [ -f "$(DIR)/.benchpackages" ]; then cat "$(DIR)/.benchpackages"; fi)
## Зависимости с зафиксированной версией (.packages_lock или .packages_lock_version).
PACKAGES_LOCK = $(shell for f in "$(DIR)/.packages_lock" "$(DIR)/.packages_lock_version"; do if [ -f $$f ]; then cat $$f; break; fi; done)
## Токен сервиса coveralls.io (.coveralls).
COVERALLS_REPO_TOKEN = $(shell if [ -f "$(DIR)/.coveralls" ]; then cat "$(DIR)/.coveralls"; fi)
## Локальные зависимости разработки (.localpackages).
LOCALPACKAGES = $(shell if [ -f "$(DIR)/.localpackages" ]; then cat "$(DIR)/.localpackages"; fi)

## Основной исполняемый файл проекта.
BIN01 := $(DIR)/bin/$(APP)
BINVERSION = $(shell "$(BIN01)" version 2>/dev/null; true)
VERN01 = $(shell echo "$(BINVERSION)" | awk -F '-' '{ print $$1 }')
VERB01 = $(shell echo "$(BINVERSION)" | awk -F 'build.' '{ print $$2 }')
## Файл с идентификатором процесса запущенного приложения.
PIDF01 := $(DIR)/run/$(APP).pid
PIDN01 = $(shell if [ -f "$(PIDF01)" ]; then cat "$(PIDF01)"; fi)

## Сценарий по умолчанию - отображение доступных команд.
default: help

## Подготовка проекта к сборке и тестированию.
## Переопределяется макросом PROJECT_LINK в файле .prj.mk.
link:
	$(call PROJECT_LINK)
.PHONY: link

## Создание директорий проекта и очистка устаревшей vendor директории.
dep-init:
	@rm -rf "$(DIR)/vendor" 2>/dev/null; true
	$(call PROJECT_FOLDERS)
.PHONY: dep-init

## Загрузка и обновление зависимостей проекта (продакшн режим).
dep: dep-init
	@for item in $(LOCALPACKAGES); do \
		PKGNAME=`echo "$${item}" | awk -F'=' '{ print $$1 }'`; \
		go mod edit -dropreplace "$${PKGNAME}" || true; \
	done
	@go mod download
	@go mod tidy
	@go mod vendor
	$(call PROJECT_DEPENDENCIES)
.PHONY: dep

## Загрузка и обновление зависимостей проекта (режим разработки, с локальными заменами).
dep-dev: dep-init
	@for item in $(LOCALPACKAGES); do \
		PKGNAME=`echo "$${item}" | awk -F'=' '{ print $$1 }'`; \
		REPLACE=`echo "$${item}" | awk -F'=' '{ print $$2 }'`; \
		go mod edit -replace "$${PKGNAME}=$${REPLACE}" || true; \
	done
	@go mod download
	@go mod tidy
	$(call PROJECT_DEPENDENCIES_DEVELOPMENT)
.PHONY: dep-dev

## Обновление всех зависимостей до последних версий.
update: upd
upd:
	@go clean -cache -modcache
	@go get -u ./...
	@for item in $(PACKAGES_LOCK); do \
		go get -u "$${item}" || true; \
	done
.PHONY: upd
.PHONY: update

## Кодогенерация (использовать только во время разработки).
## Все генерируемые пакеты перечислены в файле .gogenerate.
gen:
	@for PKGNAME in $(GOGENERATE); do \
		DB2STRUCT_DRV='$(MIGRATION_DRV_MYSQL)' \
		DB2STRUCT_DSN='$(MIGRATION_DSN_MYSQL)' \
		go generate -v "$${PKGNAME}" || exit 1; \
	done
.PHONY: gen

## Алиас команды gen.
generate: gen
.PHONY: generate

## Встраивание статических ресурсов в приложение.
## Переопределяется макросом PROJECT_EMBED в файле .prj.mk.
embed:
	$(call PROJECT_EMBED)
.PHONY: embed

## Компиляция проекта.
## Если в файле .prj.mk определён макрос PROJECT_BUILD - выполняется он,
## иначе собираются все main пакеты проекта в директорию bin/.
ifneq ($(strip $(PROJECT_BUILD)),)
build:
	@mkdir -p "$(DIR)/bin" "$(DIR)/pkg"; true
	$(call PROJECT_BUILD)
else
build:
	@mkdir -p "$(DIR)/bin"; true
	@MAINPACKAGES=`go list -f '{{if eq .Name "main"}}{{.ImportPath}} {{end}}' ./... 2>/dev/null`; \
	if [ -n "$${MAINPACKAGES}" ]; then \
		for PKG in $${MAINPACKAGES}; do \
			NAME=`basename "$${PKG}"`; \
			go build -o "$(DIR)/bin/$${NAME}" "$${PKG}" || exit 1; \
		done; \
	else \
		go build ./...; \
	fi
endif
.PHONY: build

## Запуск приложения в режиме разработки.
## Переопределяется макросом PROJECT_RUN_DEVELOPMENT в файле .prj.mk.
dev: clear
	$(call PROJECT_RUN_DEVELOPMENT)
.PHONY: dev

## Запуск приложения в продакшн режиме.
## Переопределяется макросом PROJECT_RUN_PRODUCTION в файле .prj.mk.
run:
	$(call PROJECT_RUN_PRODUCTION)
.PHONY: run

## Принудительное завершение процесса запущенного приложения.
kill:
	@if [ ! "$(PIDN01)x" = "x" ]; then \
		sudo kill -KILL "$(PIDN01)" 2>/dev/null; \
		if [ $$? -ne 0 ]; then \
			echo "Нет такого процесса: $(PIDN01)"; \
		else \
			rm "$(PIDF01)" 2>/dev/null; true; \
		fi; \
	else \
		echo "Файл '$(PIDF01)' отсутствует или пуст."; \
	fi
.PHONY: kill

## Вывод на экран версии приложения.
version: v
v:
	@"$(BIN01)" version 2>/dev/null; true
.PHONY: v
.PHONY: version

## Сборка RPM пакета.
## Имя пакета определяется по spec-файлу, расположенному в директории misc/.
RPM_OS ?= leap
RPM_SPEC = $(shell ls "$(DIR)/misc"/*.spec 2>/dev/null | head -1)
RPM_APP  = $(shell basename "$(RPM_SPEC)" .spec 2>/dev/null)
RPMBIN := $(DIR)/bin/$(RPM_APP)
RPMVERSION = $(shell "$(RPMBIN)" version 2>/dev/null; true)
RPMVERN = $(shell echo "$(RPMVERSION)" | awk -F '-' '{ print $$1 }')
RPMVERB = $(shell echo "$(RPMVERSION)" | awk -F 'build.' '{ print $$2 }')
RPMVERN = 2.0.0
RPMVERB = "build.$(DATE)"

ifeq ($(strip $(RPM_APP)),)
rpm:
	@echo "Сборка RPM пакета не настроена: поместите spec файл в директорию misc/."
else
rpm:
	@## Подготовка папок для утилиты создания RPM пакета.
	@mkdir -p "$(DIR)/rpmbuild/"{BUILD,BUILDROOT,RPMS,SOURCES,SPECS,SRPMS}; true
	@## Копирование файлов, описанных в переменной PROJECT_RPM_BUILD_SOURCE файла .prj.mk.
	@for item in $(PROJECT_RPM_BUILD_SOURCE); do \
		SRC=`echo "$${item}" | awk -F':' '{ print $$1 }'`; \
		DST=`echo "$${item}" | awk -F':' '{ print $$2 }'`; \
		cp -v "$(DIR)/$${SRC}" "$(DIR)/rpmbuild/$${DST}" || true; \
	done
	@## Копирование spec файла и бинарного файла приложения.
	@if [ -f "$(RPM_SPEC)" ]; then cp -v "$(RPM_SPEC)" "$(DIR)/rpmbuild/SPECS/$(RPM_APP).spec"; fi
	@if [ -f "$(RPMBIN)" ]; then cp -v "$(RPMBIN)" "$(DIR)/rpmbuild/SOURCES/$(RPM_APP)"; fi
	@## Выполнение проектных команд подготовки данных.
	$(call PROJECT_RPM_BUILD)
	@## Обновление секции changelog из git log.
	@if command -v changelogmaker >/dev/null 2>&1; then \
		mv "$(DIR)/rpmbuild/SPECS/$(RPM_APP).spec" "$(DIR)/rpmbuild/SPECS/src.spec"; \
		cd "$(DIR)" && changelogmaker -s "$(DIR)/rpmbuild/SPECS/src.spec" > "$(DIR)/rpmbuild/SPECS/$(RPM_APP).spec"; \
	fi
	@## Сборка RPM пакета.
	@RPM_OS="$(RPM_OS)" rpmbuild \
		--target x86_64 \
		--define "debug_package %{nil}" \
		--define "_topdir $(DIR)/rpmbuild" \
		--define "_app_version_number $(RPMVERN)" \
		--define "_app_version_build $(RPMVERB)" \
		-bb "$(DIR)/rpmbuild/SPECS/$(RPM_APP).spec"
endif
.PHONY: rpm

## Инструменты миграций баз данных.
## Используемые базы данных и подключения описываются в файле .env.mk.
####################################
COMMANDS  = up create down status redo version
MTARGETS := $(shell \
for cmd in $(COMMANDS); do \
	for drv in $(MIGRATION_DRIVER); do \
		echo "m-$${drv}-$${cmd}"; \
	done; \
done)

ifneq ($(strip $(MTARGETS)),)
MIGRATION_TMP := $(shell mktemp)
endif

## Создание директорий миграций.
migration-mkdir:
	@for dir in $$(echo "$(MIGRATION_DRIVER)"); do \
		mkdir -p "$(DIR)/migration/$${dir}" || true; \
	done
.PHONY: migration-mkdir

## Получение параметров подключения из переменных окружения.
MIGRATION_DIR = ${$(shell echo $(shell echo "${@}" | sed -e 's/^m-\(.*\)-\(.*\)$$/\1/') | awk '{ print "MIGRATION_DIR_"toupper($$0) }')}
MIGRATION_DRV = ${$(shell echo $(shell echo "${@}" | sed -e 's/^m-\(.*\)-\(.*\)$$/\1/') | awk '{ print "MIGRATION_DRV_"toupper($$0) }')}
MIGRATION_DSN = ${$(shell echo $(shell echo "${@}" | sed -e 's/^m-\(.*\)-\(.*\)$$/\1/') | awk '{ print "MIGRATION_DSN_"toupper($$0) }')}
MIGRATION_CMD = $(shell echo $(shell echo "${@}" | sed -e 's/^m-\(.*\)-\(.*\)$$/\2/'))

## Команды инструментов миграций.
migration-commands: $(MTARGETS)
$(MTARGETS): migration-mkdir
	@if [ "$(MIGRATION_CMD)" = "create" ]; then \
		read -p "Введите название миграции: " MGRNAME && \
		if [ "$${MGRNAME}" = "" ]; then MGRNAME="new"; fi && \
		echo "$${MGRNAME}" > "$(MIGRATION_TMP)"; \
	fi
	@if [ "`cat $(MIGRATION_TMP) 2>/dev/null`" != "" ] && [ "$(MIGRATION_CMD)" = "create" ]; then \
		GOOSE_DIR='$(MIGRATION_DIR)' GOOSE_DRV='$(MIGRATION_DRV)' GOOSE_DSN='$(MIGRATION_DSN)' gsmigrate create "`cat $(MIGRATION_TMP)`"; \
	else \
		GOOSE_DIR='$(MIGRATION_DIR)' GOOSE_DRV='$(MIGRATION_DRV)' GOOSE_DSN='$(MIGRATION_DSN)' gsmigrate "$(MIGRATION_CMD)"; \
	fi
	@if [ -f "$(MIGRATION_TMP)" ]; then rm "$(MIGRATION_TMP)"; fi
.PHONY: migration-commands $(MTARGETS)
####################################

## Тестирование одного или нескольких пакетов с отчётом о проценте покрытия кода тестами.
## Все тестируемые пакеты перечислены в файле .testpackages.
test: link
	@echo "mode: count" > "$(DIR)/coverage.log"
	@for PACKET in $(TESTPACKETS); do \
		touch "$(DIR)/coverage-tmp.log"; \
		go test -v -covermode=count -coverprofile="$(DIR)/coverage-tmp.log" $${PACKET} || exit 1; \
		tail -n +2 "$(DIR)/coverage-tmp.log" | sort -r | awk '{ if ($$1 != last) { print $$0; last=$$1 } }' >> "$(DIR)/coverage.log"; \
		rm -f "$(DIR)/coverage-tmp.log"; true; \
	done
	@if [ -n "$(COVERALLS_REPO_TOKEN)" ] && command -v goveralls >/dev/null 2>&1; then \
		echo "Загрузка на coveralls.io процента покрытия кода тестами."; \
		goveralls -repotoken "$(COVERALLS_REPO_TOKEN)" >/dev/null 2>&1 || true; \
	fi
.PHONY: test

## Отображение в браузере покрытия тестируемого кода в HTML-отчёте (запускать только во время разработки).
cover: test
	@go tool cover -html="$(DIR)/coverage.log"
.PHONY: cover

## Тестирование производительности.
## Все тестируемые пакеты перечислены в файле .benchpackages.
bench:
	@for PACKET in $(BENCHPACKETS); do \
		go test -race -bench=. -benchmem $${PACKET}; \
	done
.PHONY: bench

## Очистка экрана консоли.
clear:
	@clear 2>/dev/null || true
.PHONY: clear

## Очистка директории проекта от временных файлов.
clean:
	@if [ -d "$(DIR)/pkg" ]; then \
		chown -R `whoami` "$(DIR)/pkg/" 2>/dev/null; true; \
		chmod -R 0777 "$(DIR)/pkg/" 2>/dev/null; true; \
	fi
	@rm -rf "$(DIR)/bin" 2>/dev/null; true
	@rm -rf "$(DIR)/pkg" 2>/dev/null; true
	@rm -rf "$(DIR)/rpmbuild" 2>/dev/null; true
	@rm -rf "$(DIR)/run"/*.pid 2>/dev/null; true
	@rm -rf "$(DIR)/log"/*.log 2>/dev/null; true
	@rm -rf "$(DIR)"/*.log 2>/dev/null; true
	@rm -rf "$(DIR)"/*.lock 2>/dev/null; true
	@rm -rf "$(DIR)"/*.db 2>/dev/null; true
.PHONY: clean

## Помощь по командам.
help:
	@echo "Usage: make [target]"
	@echo "  target is:"
	@echo "    dep-init             - Создание директорий проекта и очистка vendor."
	@echo "    dep                  - Загрузка и обновление зависимостей проекта."
	@echo "    dep-dev              - Загрузка зависимостей с локальными заменами (режим разработки)."
	@echo "    upd или update       - Обновление всех зависимостей до последних версий."
	@echo "    gen                  - Кодогенерация с использованием go generate."
	@echo "    embed                - Встраивание статических ресурсов в приложение."
	@echo "    build                - Компиляция приложения."
	@echo "    dev                  - Запуск приложения в режиме разработки."
	@echo "    run                  - Запуск приложения в продакшн режиме."
	@echo "    kill                 - Принудительное завершение процесса приложения."
	@echo "    version              - Вывод на экран версии приложения."
	@echo "    rpm                  - Создание RPM пакета."
	@echo "    m-[driver]-[command] - Работа с миграциями базы данных."
	@echo "                           Драйверы: mysql, clickhouse, sqlite3, postgres, redshift, tidb."
	@echo "                           Команды: up, down, create, status, redo, version."
	@echo "    test                 - Запуск тестов проекта с отчётом о покрытии кода."
	@echo "    cover                - Показ процента покрытия кода тестами в браузере."
	@echo "    bench                - Запуск тестов производительности."
	@echo "    clear                - Очистка экрана консоли."
	@echo "    clean                - Очистка директории проекта от временных файлов."
	@$(call PROJECT_HELP)
.PHONY: help

endif
####################################
## Репозитории frontend.
ifeq ($(PROJECT_TYPE),web)

## Базовый путь, на котором разворачивается приложение.
BASE_HREF ?= /

## Сценарий по умолчанию - отображение доступных команд.
default: help

## Загрузка и обновление зависимостей проекта.
dep:
	@npm install --force
.PHONY: dep

## Запуск приложения в режиме разработки.
dev:
	@clear 2>/dev/null || true
	@ng cache clean
	NODE_OPTIONS="--max-http-header-size=10485760 --openssl-legacy-provider" ng serve \
		--hmr \
		--serve-path="$(BASE_HREF)" \
		--host backend \
		--port 8990 \
		--disable-host-check \
		--proxy-config proxy.conf.json \
		--allowed-hosts="all" \
		--live-reload
.PHONY: dev

## Сборка и запуск автосборки для библиотеки '@webnice/ngx'.
dev-lib: clean-pre-build
	@clear 2>/dev/null || true
	@ng cache clean
	NODE_OPTIONS="--max-http-header-size=10485760 --openssl-legacy-provider" ng build \
		@webnice/ngx \
		--watch
.PHONY: dev-lib

## Сборка библиотеки '@webnice/ngx' в директорию dist.
build-lib: clean-pre-build
	@clear 2>/dev/null || true
	@ng cache clean
	@ng build @webnice/ngx
.PHONY: build-lib

## Публикация библиотеки '@webnice/ngx' в репозиторий npm.
pub-lib: clean-pre-build
	@clear 2>/dev/null || true
	@ng cache clean
	@npm run lib:publish
.PHONY: pub-lib

## Очистка перед сборкой.
clean-pre-build:
	@rm -rf "$(DIR)/dist" 2>/dev/null; true
	@rm -rf "$(DIR)/npm-debug.log" 2>/dev/null; true
	@rm -rf "$(DIR)/ng-warnings.log" 2>/dev/null; true
	@rm -rf "$(DIR)/ng-errors.log" 2>/dev/null; true
.PHONY: clean-pre-build

## Очистка директории проекта от временных файлов.
clean: clean-pre-build
	@rm -rf "$(DIR)/node_modules" 2>/dev/null; true
	@rm -rf "$(DIR)/dist" 2>/dev/null; true
	@rm -rf "$(DIR)/nohup.out" 2>/dev/null; true
	@npm cache clean --force 2>/dev/null || true
.PHONY: clean

## Помощь по командам.
help:
	@echo "Usage: make [target]"
	@echo "  target is:"
	@echo "    dep                  - Загрузка и обновление зависимостей проекта."
	@echo "    dev                  - Запуск приложения в режиме разработки."
	@echo "    dev-lib              - Сборка и запуск автосборки библиотеки '@webnice/ngx'."
	@echo "    build-lib            - Сборка библиотеки '@webnice/ngx' в директорию dist."
	@echo "    pub-lib              - Публикация библиотеки '@webnice/ngx' в репозиторий npm."
	@echo "    clean                - Очистка директории проекта от временных файлов."
.PHONY: help

endif
