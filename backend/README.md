# Code gen with sqlc and goose

To run the server:
```bash
go run .
```

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

example no env vars:
```bash
goose sqlite3 ./theDB.db -dir=./sql/migrations create init sql
```

example with env vars:
```bash
goose create init sql
```
