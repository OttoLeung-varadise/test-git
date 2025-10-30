# This makefile just run in local, only for development
DB_HOST=localhost
DB_PWD ?=123456
DB_NAME ?=bookstore
DB_PORT ?=5432
DB_USER ?=postgres
DB_URL ?="postgres://${DB_USER}:${DB_PWD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

# run local development server
run-dev:
	DB_HOST=localhost DB_PORT=${DB_PORT} DB_USER=${DB_USER} DB_PASSWORD=${DB_PWD} DB_NAME=${DB_NAME} go run main.go

# run database migrations
migrate-up:
	migrate -verbose -path ./migrations -database ${DB_URL} up

# run database migrations down all
migrate-clearup:
	migrate -verbose -path ./migrations -database ${DB_URL} down -all

# rebuild swagger json&files
swag-init:
	swag init -g main.go -o ./docs