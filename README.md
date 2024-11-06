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
- Just

### Onboarding script

```bash
docker compose up -d
just migrate
just dev
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
PUBLIC_URL="http://localhost:4000"
FRONTEND_URL="http://localhost:3000"
GOOGLE_KEY="<REDACTED>"
GOOGLE_SECRET="<REDACTED>"
GITHUB_KEY="<REDACTED>"
GITHUB_SECRET="<REDACTED>"
MAIL_SMTP="<REDACTED>"
MAIL_PASSWORD="<REDACTED>"
SMTP_PORT="<REDACTED>"
SMTP_HOST="<REDACTED>"
APP_KEY="<REDACTED>"
ENV="development"
```

4. Run the server in development mode

```bash
just dev
# or
air
```

## Configuration

All the configuration is done through env variables :

- `PORT` : the port on which your web server runs
- `PG_USER` : the postgres user
- `PG_PASSWORD` : the postgres password
- `PG_PORT` : the postgres port
- `PG_HOST` : the postgres host
- `PG_DATABASE` : the postgres database
- `PUBLIC_URL` : The URL on which your application is available. If you use a reverse proxy to make your app available, you need to provide its URL thanks to this variable. It is mainly use as the redirection URL used during the authentication flow.
- `FRONTEND_URL` : The URL on which your frontend is available. It is mainly use to set the cookie after the authentication flow
- `GOOGLE_KEY` : your google client key for OAUTH2
- `GOOGLE_SECRET` : your google client secret for OAUTH2
- `GITHUB_KEY` : your github client key for OAUTH2
- `GITHUB_SECRET` : your github client secret for OAUTH2
- `MAIL_SMTP` : your smtp server
- `MAIL_PASSWORD` : your smtp password
- `SMTP_PORT` : your smtp port
- `SMTP_HOST` : your smtp host
- `APP_KEY` : your app key, which is a random string that is used to encrypt access tokens
- `ENV`: the environment in which the app is running (development, production)

## Tests

To run tests, execute:

```bash
just test <your-test-name>
# or to run all tests
just test-all
```

## Migrations

For migrations we use [migrate](https://github.com/golang-migrate/migrate).

To create a migration run the following command :

```bash
migrate create -ext=sql -dir=migrations -seq init
```

To run the existing migrations locally :

```bash
just migrate
# or
migrate -path=migrations -database "postgresql://postgres:postgres@0.0.0.0:5432/gists?sslmode=disable" -verbose up
```

To rollback a migration :

```bash
migrate -path=migrations -database "postgresql://postgres:postgres@localhost:5432/gists?sslmode=disable" -verbose down 1
```
