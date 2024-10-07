## !!!Attention!!! 

**All .env and config files for demonstration purpose**

## Start service

Start docker container with postgresql

    docker compose up -d

Start APP

    go run cmd/main.go serve --config=./config/config.yml

Stop docker container with postgresql if you needed

    docker compose down -v

If you need apply some new migration use command

    docker compose up migrator
    