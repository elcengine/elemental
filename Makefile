GO_BENCH_ARGS ?= -benchtime=30s
	
format:
	gofmt -w .
test:
	PARALLEL_CONVEY=false make test-lightspeed
test-lightspeed:
	go test $(GO_TEST_ARGS) --run='${run}' -v --count=1 ./tests/...
test-coverage:
	@mkdir -p ./coverage
	make test GO_TEST_ARGS="--cover -coverpkg=./cmd/...,./core/...,./plugins/...,./utils/... --coverprofile=./coverage/coverage.out"
	go tool cover -html=./coverage/coverage.out -o ./coverage/index.html
	@echo "\033[0;32mCoverage report generated at ./coverage/index.html.\033[0m"
benchmark:
	go test -bench=. $(GO_BENCH_ARGS) -benchmem -tags=benchmark ./tests/benchmarks/... 
b1x:
	make benchmark GO_BENCH_ARGS="-benchtime=1x"
b10x:
	make benchmark GO_BENCH_ARGS="-benchtime=10x"
b1s:
	make benchmark GO_BENCH_ARGS="-benchtime=1s"
b10s:
	make benchmark GO_BENCH_ARGS="-benchtime=10s"
lint:
	golangci-lint run ./...
lint-fix:
	golangci-lint run --fix ./...
install:
	go install github.com/evilmartians/lefthook@v1.11.12
	lefthook install
	@echo "\033[0;32mLefthook installed and configured successfully.\033[0m"
	go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.6
	@echo "\033[0;32mGolangCI-Lint installed successfully.\033[0m"
	@which npm > /dev/null && \
		npm install -g @commitlint/config-conventional@17.6.5 @commitlint/cli@17.6.5 && \
		echo "\033[0;32mCommitlint installed successfully.\033[0m" || \
		echo "\033[0;31mNode is not installed. Please install Node.js to use commitlint.\033[0m"
	go mod tidy
	@echo "\033[0;32mGo modules installed successfully.\033[0m"
tidy:
	@echo "\033[0;32mRunning go mod tidy...\033[0m"
	go mod tidy -v
	@echo "\033[0;32mVerifying packages...\033[0m"
	go mod verify