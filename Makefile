
mock:
	mockgen -destination=internal/mocks/mock_system_service.go -package=mocks metrics/internal/core/service Pinger	

test:
	go test ./...

server:
	go run cmd/server/main.go -l debug -d "host=localhost port=45432 user=username password=password dbname=metrics sslmode=disable"

agent:
	go run cmd/agent/main.go -l debug