.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/gcli .
	sudo cp ./bin/gcli /usr/bin/
	go install .


.PHONY: build-darwin-amd64
build-darwin-amd64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/gcli .
	mkdir -p ~/bin
	cp ./bin/gcli ~/bin/
	go install .

.PHONY: build-darwin-arm64
build-darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ./bin/gcli .
	mkdir -p ~/bin
	cp ./bin/gcli ~/bin/
	go install .
	
.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/gcli.exe .
	go install .
