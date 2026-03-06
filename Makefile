include .env

export

run:
	go run main.go

migrate-up:
	migrate -database "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable" -path migrations up

migrate-down:
	migrate -database "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable" -path migrations down

migrate-dirty:
	migrate -database "postgres://postgres:pass@localhost:5432/postgres?sslmode=disable" -path migrations force 1
