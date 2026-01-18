sqlc:
	sqlc generate


add_migration:
	migrate create -ext sql -dir internal/db/migration -seq $(name)

.PHONY: sqlc add_migration