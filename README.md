# Gists.app API - Code snippets and scripts safe

[![Go Report Card](https://goreportcard.com/badge/github.com/gistsapp/api)](https://goreportcard.com/report/github.com/kubernetes/kubernetes) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/gistsapp/api?sort=semver)

## Usage

Check the [API documentation](http://localhost:4000) for more information (for now this page is hosted locally, so please run the project before accessing the documentation).

## Quick Start

### Pre-requisites

- Go (version 1.22+)
- Air
- docker compose
- migrate

### Onboarding script

```bash
docker compose up -d
migrate -path=migrations -database "postgresql://postgres:postgres@0.0.0.0:5432/gists?sslmode=disable" -verbose up
air
```

## Installation

1. Make sure that you have `air`

```bash
go install github.com/air-verse/air@latest
```

2. Clone the repository
3. Provide the following environment variables in the `.env` file

```bash
# for now none needed

PORT="4000"
PG_USER="postgres"
PG_PASSWORD="postgres"
PG_PORT="5432"
PG_HOST="0.0.0.0"
PG_DATABASE="gists"
```

4. Run the server

```bash
air
```

## Tests

To run tests, execute:

```bash
go test .
```

## Migrations

For migrations we use [migrate](https://github.com/golang-migrate/migrate).

To create a migration run the following command :

```bash
migrate create -ext=sql -dir=migrations -seq init
```

To run the existing migrations locally :

### With bash

```bash
migrate -path=migrations -database "postgresql://postgres:postgres@0.0.0.0:5432/gists?sslmode=disable" -verbose up
```

### With go

```bash
go build main.go
./main migrate
```
