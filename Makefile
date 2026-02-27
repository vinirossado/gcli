DEBUG_DIR := $(CURDIR)/debug
DEBUG_BIN := /tmp/gcli-debug
NAME      ?= user
ARGS      ?= create handler $(NAME)

# ─── Help ────────────────────────────────────────────────────────────────────

.PHONY: help
help: ## Show this help message
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-28s\033[0m %s\n", $$1, $$2}'

# ─── Debug: run ──────────────────────────────────────────────────────────────

.PHONY: debug-build debug debug-init debug-reset debug-create-handler debug-create-service \
        debug-create-repository debug-create-model debug-create-all \
        debug-dlv debug-clean debug-test

debug-build: ## Compile the CLI to /tmp/gcli-debug (with debug symbols)
	go build -gcflags="all=-N -l" -o $(DEBUG_BIN) .

debug: debug-build ## Build and run with custom args  [ARGS="create handler user"]
	cd $(DEBUG_DIR) && $(DEBUG_BIN) $(ARGS)

debug-init: ## Pull the debug submodule (run after cloning gcli for the first time)
	git submodule update --init debug

debug-reset: ## Discard all changes in debug/ and restore the template to its original state
	cd $(DEBUG_DIR) && git checkout -- . && git clean -fd source/

debug-create-handler: debug-build ## Create a handler in debug/  [NAME=user]
	cd $(DEBUG_DIR) && $(DEBUG_BIN) create handler $(NAME)

debug-create-service: debug-build ## Create a service in debug/  [NAME=user]
	cd $(DEBUG_DIR) && $(DEBUG_BIN) create service $(NAME)

debug-create-repository: debug-build ## Create a repository in debug/  [NAME=user]
	cd $(DEBUG_DIR) && $(DEBUG_BIN) create repository $(NAME)

debug-create-model: debug-build ## Create a model in debug/  [NAME=user]
	cd $(DEBUG_DIR) && $(DEBUG_BIN) create model $(NAME)

debug-create-all: debug-build ## Create handler+service+repository+model in debug/  [NAME=user]
	cd $(DEBUG_DIR) && $(DEBUG_BIN) create all $(NAME)

debug-dlv: ## Start Delve headless on :2345 (attach from IDE)  [ARGS="create handler user"]
	dlv debug --headless --listen=:2345 --api-version=2 --wd $(DEBUG_DIR) . -- $(ARGS)

debug-test: ## Run repository integration tests in debug/ (requires Docker)  [NAME=order]
	cd $(DEBUG_DIR) && go test ./source/repository/ -v -run "Test$(shell echo $(NAME) | awk '{print toupper(substr($$0,1,1)) substr($$0,2)}')" -timeout 120s

debug-clean: ## Remove only generated (untracked) files and restore all modified template files
	cd $(DEBUG_DIR) && git clean -fd \
		source/handler \
		source/service \
		source/repository \
		source/router \
		source/model
	cd $(DEBUG_DIR) && git checkout -- \
		source/handler \
		source/service \
		source/repository \
		source/router/http.go \
		source/model/model.go \
		source/cmd/server/wire.go \
		source/cmd/migration/wire.go

# ─── Build ───────────────────────────────────────────────────────────────────

.PHONY: build-linux
build-linux: ## Build for Linux (amd64) and copy to /usr/bin/
	go mod verify
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	Golangci-lint run ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/gcli .
	sudo cp ./bin/gcli /usr/bin/
	go install .

.PHONY: build-darwin-amd64
build-darwin-amd64: ## Build for macOS (amd64) and copy to ~/bin/
	go mod verify
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	Golangci-lint run ./...
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/gcli .
	mkdir -p ~/bin
	cp ./bin/gcli ~/bin/
	go install .

.PHONY: build-darwin-arm64
build-darwin-arm64: ## Build for macOS (arm64 / Apple Silicon) and copy to ~/bin/
	go mod verify
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	Golangci-lint run ./...
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ./bin/gcli .
	mkdir -p ~/bin
	cp ./bin/gcli ~/bin/
	go install .
	
.PHONY: build-windows
build-windows: ## Build for Windows (amd64)
	go mod verify
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	Golangci-lint run ./...
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/gcli.exe .
	go install .
