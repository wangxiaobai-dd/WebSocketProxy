GO := go

.PHONY: all
all: tidy build

.PHONY: tidy
tidy:
	@$(GO) mod tidy

.PHONY: build
build:
	@$(GO) build run/main.go

.PHONY: clean
clean:
	@rm -f main.exe