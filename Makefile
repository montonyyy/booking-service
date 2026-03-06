include .env

export

run:
	go run main.go

migrate-last:
	migrate -database ${CONN_STRING} -path migrations up

migrate-first:
	migrate -database ${CONN_STRING} -path migrations down

migrate-force:
	migrate -database ${CONN_STRING} -path migrations force 1
