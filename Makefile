
mock:
	mockgen -destination=internal/mocks/mock_system_service.go -package=mocks metrics/internal/core/service Pinger	
	mockgen -destination=internal/mocks/mock_db_store.go -package=mocks metrics/internal/core/service Store	


test:
	go test -v -coverpkg=./... -coverprofile=profile.cov.tmp ./...
	grep -Ev "mock|swagger" profile.cov.tmp > profile.cov
	go tool cover -func profile.cov

server:
	go run cmd/server/main.go -l debug -d "host=localhost port=45432 user=username password=password dbname=metrics sslmode=disable" -k SomeKey

agent:
	go run cmd/agent/main.go -v debug -k SomeKey -l 3

fmt:
	goimports -local "metrics" -w .

doc:
	godoc -http=:8000 -goroot=$(shell pwd)

swag:
	swag init -g ./cmd/server/main.go --output ./swagger