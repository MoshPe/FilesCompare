EXE_NAME=FileCompare
OUTPUT_FOLDER_WINDOWS=windows-build
.PHONY: build
build: ## Run tests
	go build ./...

.PHONY: test
test: ## Run tests
	gotestsum --format testname ./...

.PHONY: fmt
fmt: ## Format go files
	go fmt ./...

.PHONY: setup
setup: ## Setup the precommit hook
	@which pre-commit > /dev/null 2>&1 || (echo "pre-commit not installed see README." && false)
	@pre-commit install

install:
	go build -o $(EXE_NAME) .
	sudo mv $(EXE_NAME) /usr/local/bin
	sudo chmod +x /usr/local/bin/$(EXE_NAME)

.PHONY: windows-amd64
build-windows:
	#env GOOS=windows GOARCH=amd64 go build -o $(EXE_NAME).exe
	if [ ! -d $(OUTPUT_FOLDER_WINDOWS)  ]; then
		mkdir -p $(OUTPUT_FOLDER_WINDOWS); \
        echo "Created folder: $(OUTPUT_FOLDER)"; \
	fi

