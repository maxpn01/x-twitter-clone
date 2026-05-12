include .env
export

run:
	go run cmd/main.go

migrate-up:
	migrate -path migrations -database $(DATABASE_URL) up

migrate-down:
	migrate -path migrations -database $(DATABASE_URL) down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)