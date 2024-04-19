//данный модуль реазлизует клиентскую часть приложения

package moduls

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func Client(req *flag.FlagSet, host *string, port *string) {

	//Проверяет длинну и допустимость вводимых данных
	if Sanitize(*host) && Length(*host) && Sanitize(*port) && Length(*port) {
		return
	}

	//Объяденяет хост и порт в одну строку
	hostPort := strings.Join([]string{*host, *port}, ":")

	//Определяет функционал вложенных команд
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

		//Определяет поведение кода, когда вместо флага команды `-id {id}`
		//вводится аргумент команды `{id}` без флага
		if *nameList == "" {

			if list.Args() != nil {

				//ничего не ввили
				if list.Arg(0) == "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}

				//ввели что-то лишнее
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

		//Контроль длинны и символов в водимой строке
		if Sanitize(*nameList) && Length(*nameList) {
			return
		}

		fmt.Println("GET request:")
		getItems(*nameList, hostPort)

	case "create":
		creates := flag.NewFlagSet("create", flag.ExitOnError)
		nameCreate := creates.String("name", "", "Name")

		creates.Parse(req.Args()[1:])

		//Определяет поведение кода, когда вместо флага команды `-name {Имя}`
		//вводится аргумент команды `{Имя}` без флага
		if *nameCreate == "" {

			if creates.Args() != nil {
				//ничего не ввели
				if creates.Arg(0) == "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}

				//ввели что-то лишнее
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

		//Контроль длинны и символов в водимой строке
		if Sanitize(*nameCreate) && Length(*nameCreate) {
			return
		}

		fmt.Println("POST request:")
		newItem := Item{Name: *nameCreate}
		createItem(newItem, hostPort)

	case "update":
		update := flag.NewFlagSet("update", flag.ExitOnError)
		nameUpdate := update.String("name", "New Item", "Name")
		idName := update.String("id", "", "Name")

		update.Parse(req.Args()[1:])

		//Определяет поведение кода, когда вместо флагов команды `-name {Имя} -id {id}`
		//вводится флаг и аргумент команды `-name {Имя} {id}`

		//При этом игнарирование флага `-name' не допустимо
		if *idName == "" {

			if update.Args() != nil {
				//ничего не ввели
				if update.Arg(0) == "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}

				//ввели что-то лишнее
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

		//проверка введенных значений
		if Sanitize(*nameUpdate) && Sanitize(*idName) && Length(*nameUpdate) && Length(*idName) {
			return
		}

		newUpdat := Item{ID: *idName, Name: *nameUpdate}

		fmt.Println(newUpdat)

		fmt.Println("PUT request:")
		updateItem(*idName, newUpdat, hostPort)

	case "delete":
		delete := flag.NewFlagSet("delete", flag.ExitOnError)
		idDelete := delete.String("id", "", "ID of delete")

		delete.Parse(req.Args()[1:])

		//Определяет поведение кода, когда вместо флага команды `-id {id}`
		//вводится аргумент команды `{id}` без флага
		if *idDelete == "" {

			if delete.Args() != nil {
				//ничего не ввели
				if delete.Arg(0) == "" {
					fmt.Println("You flag is not correct:")
					os.Exit(1)
				}
				//ввели что-то лишнее
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

		//проверка введенных значений
		if Sanitize(*idDelete) && Length(*idDelete) {
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
func createItem(item Item, hostPort string) {
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
func updateItem(itemID string, updatedItem Item, hostPort string) {
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
