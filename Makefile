migration-up:
	migrate -path db/migrations/ -database "postgresql://postgres:postgres@postgres-fs:5432/postgres?sslmode=disable" -verbose up

migration-down:
	migrate -path db/migrations/ -database "postgresql://postgres:postgres@postgres-fs:5432/postgres?sslmode=disable" -verbose down

migration-create:
	migrate create -ext sql -dir db/migrations/ -seq $(name)

migration-fix:
	migrate -path db/migrations/ -database "postgresql://postgres:postgres@postgres-fs:5432/postgres?sslmode=disable" force $(version)