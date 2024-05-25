#
# Simple Makefile to initial, build and setup the binary
#
# @author Arash Shams <me@arashshams.com>

APP_NAME := dot-proxy
VENDOR_DIR := vendor

# Determine Operating System to build for specific platform
UNAME := $(shell uname -s)
ifeq ($(UNAME), Linux)
	OS_FAMILY := linux
endif
ifeq ($(UNAME), Darwin)
	OS_FAMILY := darwin
endif
ifeq ($(UNAME), Windows_NT)
	OS_FAMILY := windows
endif

# https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
help: ## List of Makefile targets.
	@awk -F ':|##' \
	'/^[^\t].+?:.*?##/ {\
	printf "\033[36m%-30s\033[0m %s\n", $$1, $$NF \
	}' $(MAKEFILE_LIST)

.PHONY: init
init: ## Initializes the project.
	@go mod init github.com/ara4sh/go-dot-proxy

.PHONY: get-deps
get-deps: $(VENDOR_DIR) ## Downloads dependencies.

.PHONY: $(VENDOR_DIR)
$(VENDOR_DIR):  
	@GO111MODULE=on go mod tidy
	@GO111MODULE=on go mod vendor

.PHONY: run
run: ## Runs main.go.
	@go run ./cmd/dot-proxy/main.go

.PHONY: build
build: $(VENDOR_DIR)  ## Builds the Go binary for Linux, Darwin, and Windows platforms.
	@CGO_ENABLED=0 GOOS=$(OS_FAMILY) go build -v -ldflags '-extldflags "-static"' -installsuffix cgo  -trimpath -o ${APP_NAME} -ldflags "-s -w " ./cmd/${APP_NAME}

clean: dot-proxy  ## Removes the built binary and any other generated files. 
	@rm -f ${APP_NAME}

.PHONY: format
format: ## Fixes formatting for all go files.  
	@go fmt $$(go list ./...)

.PHONY: formatcheck
formatcheck: ## Checks formatting of all go files.
	@gofmt -l -d -e $$(find . -name '*.go' | grep -v vendor)

.PHONY: test
test: $(VENDOR_DIR) formatcheck vet ## Runs test cases.
	@go test -v -timeout 5s $$(go list ./...)

.PHONY: vet
vet:
	@go vet $$(go list ./...)

.PHONY: container
container: ## Builds a Docker container with the latest tag.
	@docker build -t ${APP_NAME} -f ./build/package/Dockerfile .

.PHONY: clean-container
clean-container: ## Removes the running container.
	@docker stop ${APP_NAME}
	@docker image rm ${APP_NAME}

.PHONY: run-container
run-container: container ## Runs the container with the latest tag and default options.
	@docker run -d --rm --name ${APP_NAME} -p 8053:8053/tcp -p 8053:8053/udp ${APP_NAME}

.PHONY: log-container
log-container: ## Checks the logs of the running container.
	@docker logs -f ${APP_NAME}

.PHONY: clean-all
clean-all: clean clean-container ## Cleans both binary and container files.
	@rm -rf vendor

.PHONY: test-dig
test-dig: ## Runs test queries using dig.
	@echo "Running TCP test query: "
	@dig @127.0.0.1 -p 8053 +tcp +short www.google.com
	@echo "Running UDP test query: "
	@dig @127.0.0.1 -p 8053 +short www.google.com
