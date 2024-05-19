.PHONY: tidy migrate seed run all

tidy:
	go mod tidy

migrate:
	dbmate -d sql/db/migrations up

seed:
	go run cmd/seed_borrowers.go

run:
	go run main.go

all: tidy migrate seed run