# Создать новую миграцию: make migrate-create имя_миграции
ROOT_DIR:=$(dir $(realpath $(lastword $(MAKEFILE_LIST))))
migrate-create:
	docker run --rm -v $(ROOT_DIR)migrations:/migrations migrate/migrate create -ext sql -dir /migrations -seq $(filter-out $@,$(MAKECMDGOALS))

# Позволяет передавать имя миграции как аргумент
%: ;
