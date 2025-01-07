GO := go

.PHONY: all
all:  build

.PHONY: tidy
tidy:
	@$(GO) mod tidy

.PHONY: build
build:
	@$(GO) build run/main.go
	@$(GO) build run/mock_client_connect.go
	@$(GO) build run/mock_client_token.go

.PHONY: client
client:
		@$(GO) build run/mock_client_connect.go
		@$(GO) build run/mock_client_token.go

.PHONY: clean
clean:
	@rm -f main.exe