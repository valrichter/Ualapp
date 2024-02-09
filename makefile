create_migration:
	# Create a new migration
	migrate create -ext sql -dir backend/db/migrations -seq $(name)