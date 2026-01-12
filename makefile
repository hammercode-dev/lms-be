serve-http:
	@go run main.go http
migrate:
	@go run main.go migrate
migrate-up:
	@go run main.go migrate:up
migrate-down:
	@go run main.go migrate:down
migrate-create:
	@go run main.go migrate:create
migrate-fresh:
	@go run main.go migrate:fresh
.PHONY: run
run:
	@swag init --parseDependency -g main.go && air http
