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


## `Пример запросов:`
1) Создать запись c именем User в списке item: `./client request --host localhost:8080/ list --name User  ` 
2) Вывести список всех записей в item:  `./client request --host localhost:8080/ list`
3) Обновить обновить имя поьзователя с заданным  ID: `./client request --host localhost:8080/ update --name NewNameUser --id 1`
