server:
	# Run HTTP server for Ualapp
	cd backend && go run cmd/main.go

create_migration:
	# Create a new migration
	migrate create -ext sql -dir backend/db/migrations -seq $(name)

postgres_up:
	# Start Postgres service
	docker compose up -d

postgres_down:
	# Stop Postgres service
	docker compose down

db_up:
	# Database up
	docker exec -it ualapp_postgres createdb --username=root --owner=root ualapp

db_down:
	# Database down
	docker exec -it ualapp_postgres dropdb --username=root ualapp

# Run migrations
DB_URL=postgresql://root:secret@localhost:5432/ualapp?sslmode=disable
migrate_up:
	# Migrate up
	migrate -path backend/db/migrations -database "$(DB_URL)" -verbose up

migrate_down:
	# Migrate down
	migrate -path backend/db/migrations -database "$(DB_URL)" -verbose down

sqlc:
	# Run sqlc
	cd backend && sqlc generate

test:
	# Run tests
	 cd backend && go test -v -cover ./...