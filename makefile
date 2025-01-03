serve-http:
	@go run main.go http
migrate:
	@go run main.go migrate
.PHONY: run
run: 
	@swag init --parseDependency -g main.go && air http
