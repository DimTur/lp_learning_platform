Схема БД: https://dbdiagram.io/d/learning_platform-66f26210a0828f8aa6d9b69b

Для работы с базой данных я использую библиотеку github.com/jackc/pgx/v5. 
Она используется для более гибкой работы с чистыми SQL запросами и их пониманием, а также для ускорения работы приложения. 

В рамках ДЗ реализованы транзакции, джоины, а также пагинация. Последняя реализована для channels, lessons, plans, потому как их может быть достаточно много и лучше получать данные пачками. Для остальных данных это не нужно, исходя из моей бизнес логики прохожедения контента.

P.S. В целом работа над сервисом не закончена: не доконца реализована бизнес логика. Под неё также будут написаны дополнительные SQL запросы к БД. 

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
    