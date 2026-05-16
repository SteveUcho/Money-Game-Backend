# Code gen with sqlc and goose

## Setup

Install goose on macOS
```bash
brew install goose
```

Usage: `goose DRIVER DBSTRING [OPTIONS] COMMAND`
OR
Set environment variables in .env:
- `GOOSE_DRIVER`
- `GOOSE_DBSTRING`
- `GOOSE_MIGRATION_DIR`

- for use with CockroachDB use `postgres` driver

example no env vars:
```bash
goose sqlite3 ./theDB.db -dir=./sql/migrations create init sql
```

example with env vars:
```bash
goose create init sql
```
