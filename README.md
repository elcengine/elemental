<br/>

<p align="center">
  <img src="https://github.com/user-attachments/assets/5aa4922a-ded3-4d51-8c63-94cd2fe09127" width="200" height="200" alt="Elemental Logo"/>
</p>

<p align="center">
  <a aria-label="License" href="https://github.com/elcengine/elemental/blob/main/LICENSE">
    <img alt="" src="https://img.shields.io/badge/License-MIT-yellow.svg">
  </a>
  <a aria-label="CI Tests" href="https://github.com/elcengine/elemental/actions/workflows/tests.yml">
    <img alt="" src="https://github.com/elcengine/elemental/actions/workflows/tests.yml/badge.svg">
  </a>
</p>

<hr/>

<br/>

Elemental is inspired by multiple ODMs and ORMs such as [Mongoose](https://mongoosejs.com), [TypeORM](https://typeorm.io), and [Eloquent](https://laravel.com/docs/12.x/eloquent) and its primary purpose is to improve developer experience without loss of performance or extensibility

## [Documentation](https://elcengine.github.io/docs/intro)

## Installation

```bash
go get github.com/elcengine/elemental
```

## CLI Installation
Elemental also provides a CLI to help you with migrations and seeding your database.

```bash
go install github.com/elcengine/elemental@latest
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
- Run `make benchmark` to run all benchmarks.
- Run `make lint` to run the linter.
- Run `make format` to format all files.

## Contributing

Contributions are more than welcome, as well as any suggestions / things you would differently to improve developer experience, etc...

Just open an issue or pull request and I'll surely go through it when time permits

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
