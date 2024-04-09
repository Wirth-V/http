# `Создание своего клиент-серверного приложения на Golang.`

## `Сборка проекта:`
1) `cd {путь до папки с проектом}/http/net_http/postgresbd` -- в файле прописано, что БД расположена в хосте localhost на порту 6667
2) `docker-compose -f postgresbd.yml up -d`
3) Зайти в только что созданную базу `postgres:12.1-alpine` и создать таблицу `clients` с полями `id` и `name`
4) `cd ..`
5) `go mod init app`
6) `go get github.com/google/uuid`
7) `go get github.com/jackc/pgx/v5`
8) 
9) `go build`

## `Запуск сервера:`
1) По команде `./app start [--host {host_name}] [--port {port_number}]` поднимется web-сервер доступный по адресу `http://{localhost или host_name}:{8080 или port_number}`

## `Запуск клиента:`
1) По команде `./app request [--host {host_name}] [--port {port_number}] {вложенная_команда}` поднимется клиент обращающийся к адрессу `http://{localhost или host_name}:{8080 или port_number}`.
2) Список допустимых вложенных команд:
  - `list` - выполняет запрос GET /items/
  - `get {id}` или `get --id {id}` - выполняет GET /items/{id}
  - `create --name {название}` или `create {название}`  - выполняет POST /items/
  - `update --name {название} {id}` или `update --name {название} --id {id}` - PUT /items/{id}
  - `delete {id}` или `delete -id {id}` - DELETE /items/{id}

## `Пример команд для запуска сервера`
1) `./app start  ` 
2) `./app start -host localhost -port 8080`
3)  `./app start -port 9090 -host localhost`

## `Пример команд для клиента`
1) `./app request create -name Diablo`
2) `./app request list `
3) `./app request -host localhost -port 9091 create -name User`
4) `./app request -host localhost -port 9091 create User`
5) `./app request -host localhost -port 9091 list`
6) `./app request -host localhost -port 8080 get -id 2a58ab85`
7) `./app request -host localhost -port 8080 update -name USER -id cd4aec7d`
8) `./app request -host localhost -port 8080 update -name USER cd4aec7d`
9) `./app request -host localhost -port 8080 delete -id 3390b10a`
