install:
	go mod tidy

test:
	go clean -testcache
	go test ./...

test-v:
	go test -v ./...  # Modo verbose

test-cover:
	go clean -testcache
	go test -short -coverprofile=coverage.out ./... 2>&1
	go tool cover -func=coverage.out

build-all:
	mkdir -p bin
	# Linux build
	GOOS=linux GOARCH=amd64 go build -o bin/linux/capital-gains ./cmd/main.go
	
	# Windows build
	GOOS=windows GOARCH=amd64 go build -o bin/windows/capital-gains.exe ./cmd/main.go
	
	# macOS builds
	GOOS=darwin GOARCH=amd64 go build -o bin/mac/capital-gains-intel ./cmd/main.go    # To Intel Mac
	GOOS=darwin GOARCH=arm64 go build -o bin/mac/capital-gains-apple ./cmd/main.go    # To Apple Silicon (M1, M2, M3)
