# Создать новую миграцию: make migrate-create имя_миграции
ROOT_DIR:=$(dir $(realpath $(lastword $(MAKEFILE_LIST))))

migrate-create:
	docker run --rm -v $(ROOT_DIR)migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	docker run --rm -v $(ROOT_DIR)migrations:/migrations migrate/migrate -path /migrations -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up

# Линтеры и проверки
GOLANGCI_IMAGE=golangci/golangci-lint

lint:
	docker run --rm -v $(ROOT_DIR):/app -w /app $(GOLANGCI_IMAGE) golangci-lint run

lint-fix:
	docker run --rm -v $(ROOT_DIR):/app -w /app $(GOLANGCI_IMAGE) golangci-lint run --fix

# Локальный запуск без Docker (требуется установленный golangci-lint)
lint-local:
	go mod download
	golangci-lint run

lint-local-fix:
	go mod download
	golangci-lint run --fix

# Вендоризация зависимостей (альтернатива загрузке модулей)
vendor:
	go mod vendor

# Линт с вендором (Windows cmd): GOFLAGS=-mod=vendor, чтобы брать зависимости из vendor
lint-vendor: vendor
	set GOFLAGS=-mod=vendor && golangci-lint run

vet:
	go vet ./...

fmt:
	gofmt -s -w .

# Позволяет передавать имя миграции как аргумент
%: ;
