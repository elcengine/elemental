# elemental

Golang ODM for MongoDB

## Installation

```bash
go get github.com/elcengine/elemental
```

## Development Setup

### Prerequisites
 - [Go 1.21 or later](https://golang.org/dl) - The Go programming language
 - [Node (optional)](https://nodejs.org/en) - If you want to make use of [commitlint](https://commitlint.js.org)


### Getting started

- Run `make install` to download all dependencies and install the required tools. This is required only once. Afterwards you could use the traditional `go mod tidy` for dependency management.
- Run `make test` to run all tests suites.
- Run `make test-lightspeed` to run the same above tests cost faster at the cost of readability.
- Run `make test-coverage` to run all test suites and generate a coverage report. Executes with `make test-lightspeed` under the hood.
- Run `make lint` to run the linter.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.