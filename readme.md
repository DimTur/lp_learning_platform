**Before start you need to apply migrations**

    go run ./cmd/migrator --storage-path=./storage/lp.db --migrations-path=./migrations

**Start service**

    go run cmd/main.go serve --config=./config/config.yaml
