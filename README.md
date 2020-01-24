# Gouthy - Simple SSO authenticator for personal projects

## Installation

TBD


## Development

Development notes

### Create a new migrate

This project is using the [go migrate tool](https://github.com/golang-migrate/migrate).
CLI tool documentation can be found [here](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)


#### Create a new migration

```bash
migrate create -ext sql -dir db/migrations/psql <name>
```

#### Run the migration on the testing env

```bash
migrate -database 'postgres://postgres:postgres@localhost:5432/gouthy?sslmode=disable' -path db/migrations/psql up
```

#### Discard all the migrations

```bash
migrate -database 'postgres://postgres:postgres@localhost:5432/gouthy?sslmode=disable' -path db/migrations/psql down
```