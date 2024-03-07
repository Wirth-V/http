package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	// GET запрос
	response, err := http.Get("http://localhost:8080")
	if err != nil {
		fmt.Println("Ошибка при выполнении GET запроса:", err)
		return
	}
	defer response.Body.Close()

	// Чтение ответа
	body, err := ioutil.ReadAll(response.Body) // читаем ответ
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа:", err)
		return
	}
	fmt.Println("Ответ на GET запрос:", string(body))

	// POST запрос
	url := "http://localhost:8080/post"
	payload := strings.NewReader("Пример POST запроса")

	request, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println("Ошибка при создании POST запроса:", err)
		return
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Отправка POST запроса
	postResponse, err := http.DefaultClient.Do(request)
	if err != nil {
		fmt.Println("Ошибка при выполнении POST запроса:", err)
		return
	}
	defer postResponse.Body.Close()

	// Чтение ответа на POST запрос
	postBody, err := ioutil.ReadAll(postResponse.Body)
	if err != nil {
		fmt.Println("Ошибка при чтении тела ответа на POST запрос:", err)
		return
	}
	fmt.Println("Ответ на POST запрос:", string(postBody))
}
