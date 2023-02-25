.PHONY: gen
gen:
	protoc --proto_path=internal/controller/grpcapi/proto/blacklist internal/controller/grpcapi/proto/blacklist/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/controller/grpcapi/proto/whitelist internal/controller/grpcapi/proto/whitelist/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/controller/grpcapi/proto/bucket internal/controller/grpcapi/proto/bucket/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import
	protoc --proto_path=internal/controller/grpcapi/proto/authorization internal/controller/grpcapi/proto/authorization/*.proto  --go_out=. --go_opt=paths=import --go-grpc_out=. --go-grpc_opt=paths=import

.PHONY: clean
clean:
	rm -f internal/controller/grpcapi/blacklistpb/*
	rm -f internal/controller/grpcapi/whitelistpb/*
	rm -f internal/controller/grpcapi/bucketpb/*
	rm -f internal/controller/grpcapi/authorization/*

.PHONY: build.docker
build.docker:
	docker build --tag abf --  .

.PHONY: run.docker
run.docker:
	docker run -p 8080:8080 -it --name abf_container abf

.PHONY: build
build:
	docker-compose build

.PHONY: build.bin
build.bin:
	go build -o ./build/anti_bruteforce_app/anti_bruteforce_service ./cmd/anti_bruteforce_app

.PHONY: run
run: build
	docker-compose up

.PHONY: stop
stop:
	docker-compose down

.PHONY: lint
lint:
	golangci-lint run

.PHONY: migrate
migrate:
	migrate -version $(version)

.PHONY: migrate.down
migrate.down:
	migrate -source file://migrations -database postgres://localhost:5433/anti-bruteforce-service-database?sslmode=disable down

.PHONY: migrate.up
migrate.up:
	migrate -source file://migrations -database postgres://localhost:5433/anti-bruteforce-service-database?sslmode=disable up
