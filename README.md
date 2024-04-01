# `Создание своего клиент-серверного приложения на Golang.`

## `Сборка проекта:`
1) `cd {путь до папки с проектом}/http/net_http`
2) `go mod init app`
3) `go get github.com/google/uuid`
4) `go build`

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

## `Пример запросов:`
1) Создать запись c именем User в списке item: `./client request --host localhost:8080/ list --name User  ` 
2) Вывести список всех записей в item:  `./client request --host localhost:8080/ list`
3) Обновить обновить имя поьзователя с заданным  ID: `./client request --host localhost:8080/ update --name NewNameUser --id 1`
