check: vet test build386 buildlinuxarm # Perform basic checks and tests

vet: ## Run go vet
	go vet ./...

tidy:	## Run go mod tidy
	go mod tidy

test: ## Run all tests
	go test ./...

buildall: build386 buildlinuxarm ## Build all platforms

build386: ## Build for linux/386
	GOOS=linux GOARCH=386 go install .

buildlinuxarm: ## Build for linux/arm
	GOOS=linux GOARCH=arm go install .