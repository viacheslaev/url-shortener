# URL Shortener

## DB Migrations

This project uses **golang-migrate** to manage PostgreSQL schema changes.  
Migrations are stored in:

```
migrations/
```

### Install migrate

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
Verify:

```bash
migrate -version
```
## Apply migrations

From the project root:

```bash
migrate -path migrations -database "postgres://USER:PASSWORD@localhost:5432/DBNAME?sslmode=disable" up
```

## Roll back last migration

```bash
migrate -path migrations -database "postgres://USER:PASSWORD@localhost:5432/DBNAME?sslmode=disable" down 1
```

## Check current version

```bash
migrate -path migrations -database "postgres://USER:PASSWORD@localhost:5432/DBNAME?sslmode=disable" version
