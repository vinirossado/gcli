.PHONY: build-linux
build-linux:
	go mod verify
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	Golangci-lint run ./...
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/gcli .
	sudo cp ./bin/gcli /usr/bin/
	go install .

.PHONY: build-darwin-amd64
build-darwin-amd64:	
	go mod verify
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	Golangci-lint run ./...
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/gcli .
	mkdir -p ~/bin
	cp ./bin/gcli ~/bin/
	go install .

.PHONY: build-darwin-arm64
build-darwin-arm64:
	go mod verify
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	Golangci-lint run ./...
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ./bin/gcli .
	mkdir -p ~/bin
	cp ./bin/gcli ~/bin/
	go install .
	
.PHONY: build-windows
build-windows:
	go mod verify
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest -checks=all,-ST1000,-U1000 ./...
	Golangci-lint run ./...
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/gcli.exe .
	go install .
