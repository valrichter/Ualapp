server:
	# Run HTTP server for Ualapp
	go run cmd/main.go

create_migration:
	# Create a new migration
	migrate create -ext sql -dir db/migrations -seq $(name)

postgres_up:
	# Start Postgres service
	docker compose up -d

postgres_down:
	# Stop Postgres service
	docker compose down

db_up:
	# Database up
	docker exec -it ualapp_postgres createdb --username=root --owner=root ualapp
	docker exec -it ualapp_postgres_test createdb --username=root --owner=root ualapp

db_down:
	# Database down
	docker exec -it ualapp_postgres dropdb --username=root ualapp
	docker exec -it ualapp_postgres_test dropdb --username=root ualapp

# Run migrations
DB_SOURCE=postgresql://root:secret@localhost:5432/ualapp?sslmode=disable
DB_SOURCE_TEST=postgresql://root:secret@localhost:5433/ualapp?sslmode=disable
migrate_up:
	# Migrate up
	migrate -path db/migrations -database "$(DB_SOURCE)" -verbose up
	migrate -path db/migrations -database "$(DB_SOURCE_TEST)" -verbose up

migrate_down:
	# Migrate down
	migrate -path db/migrations -database "$(DB_SOURCE)" -verbose down $(count)
	migrate -path db/migrations -database "$(DB_SOURCE_TEST)" -verbose down $(count)

sqlc:
	# Run sqlc
	sqlc generate

test:
	# Run tests
	 go test -v -cover ./...