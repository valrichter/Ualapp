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