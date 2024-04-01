# `Создание своего клиент-серверного приложения на Golang.`

## `Сборка и запуск сервера:`
1) `cd {путь до папки с проектом}/http/net_http/server`
2) `go build .`
3) По команде `./server start [--host {host_name}] [--port {port_number}]` поднимется web-сервер доступный по адресу `http://{localhost или host_name}:{8080 или port_number}`

## `Сборка и запуск клиента:`

1) `cd {путь до папки с проектом}/http/net_http/client`
2) `go build .`
3) По команде `./client request [--host {host_name}] [--port {port_number}] {вложенная_команда}` поднимется клиент обращающийся к адрессу `http://{localhost или host_name}:{8080 или port_number}`.
4) Список допустимых вложенных команд:
  - `list` - выполняет запрос GET /items/
  - `get {id}` или `get --id {id}` - выполняет GET /items/{id}
  - `create --name {название}` или `create {название}`  - выполняет POST /items/
  - `update --name {название} {id}` или `update --name {название} --id {id}` - PUT /items/{id}
  - `delete {id}` или `delete -id {id}` - DELETE /items/{id}


## `Пример команд для запуска сервера`
1) `./server start  ` 
2) `./server start -host localhost -port 8080`
3)  `./server start -port 9090 -host localhost`

## `Пример команд для клиента`
1) `./client request  create -name Diablo`
2) ` ./client request list `
3) `./client request -host localhost -port 9091 create -name User`
4) `./client request -host localhost -port 9091 create User`
5) `./client request -host localhost -port 9091 list`
6) `./client request -host localhost -port 8080 get -id 2a58ab85`
7) `./client request -host localhost -port 8080 update -name USER -id cd4aec7d`
8) `./client request -host localhost -port 8080 update -name USER cd4aec7d`
9) `./client request -host localhost -port 8080 delete -id 3390b10a`

