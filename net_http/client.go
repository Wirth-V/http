package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Item представляет структуру данных для элементов списка.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {

	/*
		nethttp.exe request --host localhost:9999 create --name test_name
		nethttp.exe one --....
		nethttp.exe two --...
	*/

	req := flag.NewFlagSet("request", flag.ExitOnError)
	host := req.String("host", "localhost:8080", "Host")

	req.Parse(os.Args[2:])
	fmt.Println("Host:", *host)
	fmt.Println("Host:", sanitizeHost(*host))

	fmt.Printf("\n")

	switch req.Arg(0) {

	case "list":
		list := flag.NewFlagSet("list", flag.ExitOnError)
		nameList := list.String("id", "", "ID")

		list.Parse(req.Args()[1:])

		fmt.Println(*nameList)
		fmt.Println(sanitizeInput(*nameList))
		fmt.Println("GET request:")
		getItems(sanitizeInput(*nameList), sanitizeHost(*host))

	case "create":
		creates := flag.NewFlagSet("create", flag.ExitOnError)
		nameCreate := creates.String("name", "New Item", "Name")

		creates.Parse(req.Args()[1:])

		fmt.Println("POST request:")
		fmt.Println(*nameCreate)
		newItem := Item{Name: sanitizeInput(*nameCreate)}
		createItem(newItem, sanitizeHost(*host))

	case "update":
		update := flag.NewFlagSet("update", flag.ExitOnError)
		idName := update.String("id", "1", "Name")
		nameUpdate := update.String("name", "New Item", "Name")

		update.Parse(req.Args()[1:])

		fmt.Println("PUT request:")
		updateItem(sanitizeInput(*idName), Item{Name: sanitizeInput(*nameUpdate)}, sanitizeHost(*host))

	case "delete":
		delete := flag.NewFlagSet("delete", flag.ExitOnError)
		idDelete := delete.String("id", "1", "ID of delete")

		delete.Parse(req.Args()[1:])

		fmt.Println("DELETE request:")
		deleteItem(sanitizeInput(*idDelete), sanitizeHost(*host))

	default:
		fmt.Println("You flag is not correct:")
		os.Exit(1)
	}
}

// getItems отправляет GET-запрос на сервер для получения списка элементов.
func getItems(nameList string, host string) {
	//var control string = ""
	var resp *http.Response
	var err error
	// Отправка GET-запроса на сервер по указанному URL.
	if nameList == "" {
		resp, err = http.Get(fmt.Sprintf("http://%s/items/", host))
	} else {
		resp, err = http.Get(fmt.Sprintf("http://%s/items/%s", host, nameList))
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
func createItem(item Item, host string) {
	// Кодирование структуры Item в JSON.
	itemJSON, err := json.Marshal(item)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}

	// Отправка POST-запроса на сервер с данными в формате JSON.
	resp, err := http.Post(fmt.Sprintf("http://%s/items/", host), "application/json", strings.NewReader(string(itemJSON)))
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
func updateItem(itemID string, updatedItem Item, host string) {
	// Кодирование обновленной структуры Item в JSON.
	itemJSON, err := json.Marshal(updatedItem)
	if err != nil {
		fmt.Println("Ошибка при кодировании JSON:", err)
		return
	}

	// Создание клиента для отправки PUT-запроса.
	client := &http.Client{}
	req, err := http.NewRequest(
		"PUT", fmt.Sprintf("http://%s/items/%s", host, itemID), strings.NewReader(string(itemJSON)))
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
func deleteItem(itemID string, host string) {
	// Создание клиента для отправки DELETE-запроса.
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/items/%s", host, itemID), nil)
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

//ниже реализованы 2 подхода к экранированию строк, через пакет regexp и через работу со строками

// sanitizeInput очищает строку от специальных символов
func sanitizeHost(input string) string {
	var result strings.Builder
	for _, char := range input {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || (char == ':') {
			result.WriteRune(char)
		}
	}
	return result.String()
}

// sanitizeInput очищает строку от специальных символов
func sanitizeInput(input string) string {
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	return reg.ReplaceAllString(input, "")
}
