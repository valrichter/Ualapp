server:
	# Run HTTP server for Ualapp
	go run cmd/main.go

create_migration:
	# Create a new migration
	migrate create -ext sql -dir db/migrations -seq $(name)

app_up:
	# Start Postgres service
	docker compose up -d

app_down:
	# Stop Postgres service
	docker compose down
	docker rmi ualapp_api:latest

db_up:
	# Database up
	docker exec -it ualapp_postgres createdb --username=root --owner=root ualapp

db_down:
	# Database down
	docker exec -it ualapp_postgres dropdb --username=root ualapp

# Run migrations
DB_SOURCE=postgresql://root:secret@localhost:5432/ualapp?sslmode=disable
migrate_up:
	# Migrate up
	migrate -path db/migrations -database "$(DB_SOURCE)" -verbose up

migrate_down:
	# Migrate down
	migrate -path db/migrations -database "$(DB_SOURCE)" -verbose down $(count)

sqlc:
	# Run sqlc
	sqlc generate

docker_app:
	docker build -t ualapp_api .
	docker run -p 8080:8080 ualapp_api

test:
	# Run tests
	 go test -v -cover ./...