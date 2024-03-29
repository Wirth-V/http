# `Создание своего клиент-серверного приложения на Golang.`

По команде `{название_приложения} start [--port {port_number}]` поднимать web-сервер доступный по адресу `http://localhost:{8080 или port_number}`.

## `Сборка:`
1) `cd $HOME/http/net_http/server`
2) `go build .`
3) `cd ../client`
4) `go build .`
5) `cd ~`

## `Запуск программы:`
1) `cd {корень репощитория}/http/net_http`-- переходим в директрорию, в которой содержится проект
2) `./server` -- запускаем сервер (сервер работает на порту 8080)
3) По команде `./client request [--port {port_number}] {вложенная_команда}` выполнять запросы в зависимости от вложенной команды:
  - `list` - выполняет запрос GET /items/
  - `get {id}` - выполняет GET /items/{id}
  - `create --name {название}` - выполняет POST /items/
  - `update --name {название} {id}` - PUT /items/{id}
  - `delete {id}` - DELETE /items/{id}

## `Пример запросов:`
1) Создать запись c именем User в списке item: `./client request --host localhost:8080/ list --name User  ` 
2) Вывести список всех записей в item:  `./client request --host localhost:8080/ list`
3) Обновить обновить имя поьзователя с заданным  ID: `./client request --host localhost:8080/ update --name NewNameUser --id 1`
