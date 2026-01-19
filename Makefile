CONTAINER_NAME := Moludb
DB_USER        := molu
DB_NAME        := sms
NOW            := $(shell date +%Y%m%d_%H%M%S)
BACKUP_NAME    := backup_$(NOW).dump
CONTAINER_BACKUP_PATH := /tmp/$(BACKUP_NAME)
HOST_BACKUP_DIR       := ./db_backups
RESTORE_FILE   ?= latest.dump  # Default now matches the symlink name

createdb:
	docker exec -it $(CONTAINER_NAME) \
		createdb --username=$(DB_USER) --owner=$(DB_USER) $(DB_NAME)

dropdb:
	docker exec -it $(CONTAINER_NAME) \
		dropdb -U $(DB_USER) $(DB_NAME)


sqlc:
	sqlc generate


add_migration:
	migrate create -ext sql -dir internal/db/migration -seq $(name)

start:
	docker start $(CONTAINER_NAME)

migrateup:
	migrate -path internal/db/migration \
		-database "postgresql://$(DB_USER):incorrect@localhost:5432/$(DB_NAME)?sslmode=disable" \
		-verbose up

migratedown:
	migrate -path internal/db/migration \
		-database "postgresql://$(DB_USER):incorrect@localhost:5432/$(DB_NAME)?sslmode=disable" \
		-verbose down

.PHONY: sqlc add_migration createdb dropdb start migrateup migratedown