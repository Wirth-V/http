/*
1. Придумать название для консольного приложения;
2. Приложение должно:
2.1. По команде {название_приложения} start [--port {port_number}] поднимать web-сервер доступный по адресу http://localhost:{8080 или port_number}.
Web-cервер будет хранить список структур
type Item struct {
  ID string
  Name string
}

Сервер должен обрабатывать следующие запросы:
+--------------------------+--------------------------------------------------------+---------------------+-----------------------------+
| URL                      | Описание                                               | json-формат запроса | json-формат ответа          |
+==========================+========================================================+=====================+=============================+
| GET /items/              | возвращает список item'ов                              | -                   | [{"id":"", "name":""}, ...] |
| GET /items/{item_id}/    | возвращает item у которого ID == item_id               | -                   | {"id":"", "name":""}        |
| POST /items/             | добавляет item со уникальным ID и переданным названием | {"name":"..."}      | {"id":"", "name":""}        |
| PUT /items/{item_id}/    | изменяет название item'а с соответствующим ID          | {"name":"..."}      | - или {"id":"", "name":""}  |
| DELETE /items/{item_id}/ | удаляет item                                           | -                   | -                           |
+--------------------------+--------------------------------------------------------+---------------------+-----------------------------+

Все запросы принимающие {item_id} должны возвращать NotFound (404) если item'а с таким id не существует.
Название (name) не может быть пустым. Если пустое - BadRequest (400)
Если всё хорошо OK (200). Можно также присылать Created (201) или NoContent (204) в определённых случаях

2.2. По команде {название_приложения} request [--port {port_number}] {вложенная_команда} выполнять запросы в зависимости от вложенной команды:
  - list - выполняет запрос GET /items/;
  - get {id} - выполняет GET /items/{id};
  - create --name {название} - выполняет POST /items/;
  - update --name {название} {id} - PUT /items/{id};
  - delete {id} - DELETE /items/{id};

Результаты вызовов напечатать в вывод команд.
*/

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net_http/moduls"
	"os"
	"strings"
)

/*
// Item представляет структуру данных для элементов списка.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
*/

func main() {

	/*
		nethttp.exe request --host localhost:9999 create --name test_name
		nethttp.exe one --....
		nethttp.exe two --...
	*/

	req := flag.NewFlagSet("request", flag.ExitOnError)
	host := req.String("host", "localhost", "Host")
	port := req.String("port", "8080", "Host")

	req.Parse(os.Args[2:])

	//Проверяет длинну и допустимость вводимых данных
	if moduls.Sanitize(*host) && moduls.Length(*host) && moduls.Sanitize(*port) && moduls.Length(*port) {
		return
	}

	hostPort := strings.Join([]string{*host, *port}, ":")

	switch req.Arg(0) {
	case "list":
		list := flag.NewFlagSet("list", flag.ExitOnError)

		list.Parse(req.Args()[1:])

		fmt.Println("GET request:")
		getItems("", hostPort)

	case "get":
		list := flag.NewFlagSet("list", flag.ExitOnError)
		nameList := list.String("id", "", "ID")

		list.Parse(req.Args()[1:])

		if *nameList == "" {

			if list.Args() != nil {
				if list.Arg(0) == "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}

				if list.Arg(1) != "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}

				*nameList = list.Arg(0)
			} else {
				fmt.Println("You flag is not correct:")
				os.Exit(1)
			}

		}

		if moduls.Sanitize(*nameList) && moduls.Length(*nameList) {
			return
		}

		fmt.Println("GET request:")
		getItems(*nameList, hostPort)

	case "create":
		creates := flag.NewFlagSet("create", flag.ExitOnError)
		nameCreate := creates.String("name", "", "Name")

		creates.Parse(req.Args()[1:])

		if *nameCreate == "" {

			if creates.Args() != nil {
				if creates.Arg(0) == "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}

				if creates.Arg(1) != "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}

				*nameCreate = creates.Arg(0)
			} else {
				fmt.Println("You flag is not correct:")
				os.Exit(1)
			}
		}

		if moduls.Sanitize(*nameCreate) && moduls.Length(*nameCreate) {
			return
		}

		fmt.Println("POST request:")
		newItem := moduls.Item{Name: *nameCreate}
		createItem(newItem, hostPort)

	case "update":
		update := flag.NewFlagSet("update", flag.ExitOnError)
		nameUpdate := update.String("name", "New Item", "Name")
		idName := update.String("id", "", "Name")

		update.Parse(req.Args()[1:])

		if *idName == "" {

			if update.Args() != nil {
				if update.Arg(0) == "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}

				if update.Arg(1) != "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}

				*idName = update.Arg(0)
			} else {
				fmt.Println("You flag is not correct:")
				os.Exit(1)
			}

		}

		if moduls.Sanitize(*nameUpdate) && moduls.Sanitize(*idName) && moduls.Length(*nameUpdate) && moduls.Length(*idName) {
			return
		}

		newUpdat := moduls.Item{ID: *idName, Name: *nameUpdate}

		fmt.Println(newUpdat)

		fmt.Println("PUT request:")
		updateItem(*idName, newUpdat, hostPort)

	case "delete":
		delete := flag.NewFlagSet("delete", flag.ExitOnError)
		idDelete := delete.String("id", "", "ID of delete")

		delete.Parse(req.Args()[1:])

		if *idDelete == "" {

			if delete.Args() != nil {
				if delete.Arg(0) == "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}

				if delete.Arg(1) != "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}

				*idDelete = delete.Arg(0)
			} else {
				fmt.Println("You flag is not correct:")
				os.Exit(1)
			}
		}

		if moduls.Sanitize(*idDelete) && moduls.Length(*idDelete) {
			return
		}

		fmt.Println("DELETE request:")
		deleteItem(*idDelete, hostPort)

	default:
		fmt.Println("You flag is not correct:")
		os.Exit(1)
	}
}

// getItems отправляет GET-запрос на сервер для получения списка элементов.
func getItems(nameList string, hostPort string) {
	//var control string = ""
	var resp *http.Response
	var err error
	// Отправка GET-запроса на сервер по указанному URL.
	if nameList == "" {
		resp, err = http.Get(fmt.Sprintf("http://%s/items/", hostPort))
	} else {
		resp, err = http.Get(fmt.Sprintf("http://%s/items/%s/", hostPort, nameList))
	}
	if err != nil {
		fmt.Println("Ошибка при отправке GET-запроса:", err)
		return
	}
	defer resp.Body.Close() // всегда сначало дефери, а потом уже что-то делай

	// Читаем и конвертируем тело ответа в байты
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// Выводим содержимое тела ответа
	fmt.Println(string(bytesResp))

	// Обработка ответа сервера.
	printResponse(resp)
}

// createItem отправляет POST-запрос на сервер для создания нового элемента.
func createItem(item moduls.Item, hostPort string) {
	// Кодирование структуры Item в JSON.
	itemJSON, err := json.Marshal(item)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}

	// Отправка POST-запроса на сервер с данными в формате JSON.
	resp, err := http.Post(fmt.Sprintf("http://%s/items/", hostPort), "application/json", strings.NewReader(string(itemJSON)))
	if err != nil {
		fmt.Println("Ошибка при отправке POST-запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Читаем и конвертируем тело ответа в байты
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// Выводим содержимое тела ответа
	fmt.Println(string(bytesResp))

	// Обработка ответа сервера.
	printResponse(resp)
}

// updateItem отправляет PUT-запрос на сервер для обновления элемента с указанным ID.
func updateItem(itemID string, updatedItem moduls.Item, hostPort string) {
	// Кодирование обновленной структуры Item в JSON.
	itemJSON, err := json.Marshal(updatedItem)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}

	// Создание клиента для отправки PUT-запроса. bytes.NewBuffer(itemJSON)
	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://%s/items/%s/", hostPort, itemID), bytes.NewBuffer(itemJSON))
	if err != nil {
		fmt.Println("Ошибка при создании PUT-запроса:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправка PUT-запроса на сервер.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке PUT-запроса:", err)
		return
	}
	defer resp.Body.Close()

	// Читаем и конвертируем тело ответа в байты
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	// Выводим содержимое тела ответа
	fmt.Println(string(bytesResp))

	// Обработка ответа сервера.
	printResponse(resp)
}

// deleteItem отправляет DELETE-запрос на сервер для удаления элемента с указанным ID.
func deleteItem(itemID string, hostPort string) {
	// Создание клиента для отправки DELETE-запроса.
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/items/%s/", hostPort, itemID), nil)
	if err != nil {
		fmt.Println("Ошибка при создании DELETE-запроса:", err)
		return
	}

	// Отправка DELETE-запроса на сервер.
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Ошибка при отправке DELETE-запроса:", err)
		return
	}

	// Читаем и конвертируем тело ответа в байты
	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	// Выводим содержимое тела ответа
	fmt.Println(string(bytesResp))

	// Обработка ответа сервера.
	printResponse(resp)
}

// printResponse выводит информацию о статусе и теле ответа сервера.
func printResponse(resp *http.Response) {
	// Чтение тела ответа в байтовый массив.
	var bodyBytes []byte
	_, err := resp.Body.Read(bodyBytes)
	if err != nil && err != io.EOF {
		fmt.Println("Ошибка при чтении ответа:", err)
		return
	}

	// Вывод информации о статусе и теле ответа.
	fmt.Printf("Status Code: %d\n", resp.StatusCode)
	//fmt.Printf("Response Body: %s\n", bodyBytes) Переделай, или через readAll или через цыкл, так же разбери потоки
	fmt.Printf("---------------\n\n")

}
