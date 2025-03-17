swag:
	swag init --parseDependency -g ./internal/http/server.go


go-private:
	go env -w GOPRIVATE="bitbucket.org/play-workspace/*"

test:
	go test ./tests -v

lint:
	golangci-lint run -v

rng-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				./pkg/rng/rng.proto

overlord-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				./pkg/overlord/overlord.proto

history-proto:
	protoc --go_out=. --go_opt=paths=source_relative \
				--go-grpc_out=. --go-grpc_opt=paths=source_relative \
				./pkg/history/history.proto