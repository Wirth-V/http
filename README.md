# `Создание своего клиент-серверного приложения на Golang.`

## `Сборка проекта:`
1) `cd {путь до папки с проектом}/http/net_http/postgres_db` -- в файле postgres_db.yml, лежащем по этому адресу, прописано, что БД расположена в хосте localhost на порту 6667
2) `docker-compose -f postgres_db.yml up -d` -- создание и запуск докера с psql
3) `cd ..`
4) `go build`

## `Запуск сервера:`
1) По команде `./app start [--host {host_name}] [--port {port_number}] [--user_bd {data_base_username}] [--password_bd {data_base_password}] [--host_bd {data_base_host}] [--port_bd {data_base_port}] [--db {data_base_name}] [--table {table_name}]` поднимется web-сервер доступный по адресу `http://{localhost или host_name}:{8080 или port_number}`. По умолчанию программа обращается/создает БД server с таблицей item. Можно задать свои имена. Если таблица уже есть в системе, то программа будет работать с ней. Если нет, то программа сначало создаст ее.

## `Запуск клиента:`
1) По команде `./app request [--host {host_name}] [--port {port_number}] {вложенная_команда}` поднимется клиент обращающийся к адрессу `http://{localhost или host_name}:{8080 или port_number}`.
2) Список допустимых вложенных команд:
  - `list` - выполняет запрос GET /items/
  - `get {id}` или `get --id {id}` - выполняет GET /items/{id}
  - `create --name {название}` или `create {название}`  - выполняет POST /items/
  - `update --name {название} {id}` или `update --name {название} --id {id}` - PUT /items/{id}
  - `delete {id}` или `delete -id {id}` - DELETE /items/{id}

## `Пример команд для запуска сервера`
1) `./net_http start` 
2) `./net_http start -host localhost -port 8080`
3)  `./net_http start -port 9090 -host localhost`
3)  `./net_http start -port 9090 -host localhost -db users -table items`
4) `./net_http start -port 8080 -host localhost -user_db server -password_db 198416 -host_db localhost -port_db 6667 -db user_base -table names`

## `Пример команд для клиента`
1) `./net_http request create -name Diablo`
2) `./net_http request list `
3) `./net_http request -host localhost -port 9090 create -name User`
4) `./net_http request -host localhost -port 9090 create User`
5) `./net_http request -host localhost -port 9090 list`
6) `./net_http request -host localhost -port 8080 get -id 2a58ab85`
7) `./net_http request -host localhost -port 8080 update -name USER -id cd4aec7d`
8) `./net_http request -host localhost -port 8080 update -name USER cd4aec7d`
9) `./net_http request -host localhost -port 8080 delete -id 3390b10a`
