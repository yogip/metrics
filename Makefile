VERSION?=1.0.0
COMMIT=$(if $(shell git rev-parse HEAD),$(shell git rev-parse HEAD),"N/A")
DATE=$(shell date "+%Y/%m/%d %H:%M:%S")
LDFLAGS=-ldflags "-s -w -X 'main.buildVersion=$(VERSION)' -X 'main.buildDate=$(DATE)' -X 'main.buildCommit=$(COMMIT)'"

mock:
	mockgen -destination=internal/mocks/mock_system_service.go -package=mocks metrics/internal/core/service Pinger	
	mockgen -destination=internal/mocks/mock_db_store.go -package=mocks metrics/internal/core/service Store	

test:
	go test -v -coverpkg=./... -coverprofile=profile.cov.tmp ./...
	grep -Ev "mock|swagger" profile.cov.tmp > profile.cov
	go tool cover -func profile.cov

server:
	go run $(LDFLAGS) cmd/server/main.go -l debug -crypto-key private.pem -d "host=localhost port=45432 user=username password=password dbname=metrics sslmode=disable" -k SomeKey

grpc-server:
	go run $(LDFLAGS) cmd/server/main.go -b grpc -t 127.0.0.1/24 -l debug -d "host=localhost port=45432 user=username password=password dbname=metrics sslmode=disable"


agent:
	go run $(LDFLAGS) cmd/agent/main.go -v debug -k SomeKey -l 3 -crypto-key public.pem
grpc-agent:
	go run $(LDFLAGS) cmd/agent/main.go -t grpc -v debug


fmt:
	goimports -local "metrics" -w .

lint:
	go run ./cmd/staticlint ./cmd/... ./internal/...

lint-fix:
	go run ./cmd/staticlint -fix ./cmd/... ./internal/...

doc:
	godoc -http=:8000 -goroot=$(shell pwd)

swag:
	swag init -g ./cmd/server/main.go --output ./swagger


keys:
	openssl genrsa -out private.pem 4096
	openssl rsa -in private.pem -outform PEM -pubout -out public.pem


proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. \ 
		--go-grpc_opt=paths=source_relative internal/proto/metrics.proto

