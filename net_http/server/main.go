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
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net_http/moduls"
	"os"
	"strings"

	"github.com/google/uuid"
)

/*
// Item представляет структуру данных для элементов списка.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
*/

// items - глобальная переменная, представляющая соотношение элементов по их уникальным ID.
var items = make(map[string]*moduls.Item)

func main() {

	req := flag.NewFlagSet("start", flag.ExitOnError)
	host := req.String("host", "localhost", "Port")
	port := req.String("port", "8080", "Port")

	req.Parse(os.Args[2:])

	moduls.InfoLog.Println("Сервер запущен.")
	moduls.InfoLog.Printf("Хост:%s Порт:%s", *host, *port)

	// Регистрация обработчика запросов для пути "/items/".
	http.HandleFunc("GET /items/", handleGET)
	// Регистрация обработчика запросов для пути "/items/".
	http.HandleFunc("GET /items/{id}/", handleGETid)
	// Регистрация обработчика запросов для пути "/items/".
	http.HandleFunc("POST /items/", handlePOST)
	// Регистрация обработчика запросов для пути "/items/".
	http.HandleFunc("PUT /items/{id}/", handlePUT)
	// Регистрация обработчика запросов для пути "/items/".
	http.HandleFunc("DELETE /items/{id}/", handleDELETE)

	// Запуск веб-сервера на порту 8080.
	err := http.ListenAndServe(strings.Join([]string{*host, *port}, ":"), nil)
	if err != nil {
		moduls.ErrorLog.Fatal("Ошибка запуска сервера:", err)
	}
}

// handleGET - обработчик для HTTP-запросов методом GET.
func handleGET(w http.ResponseWriter, r *http.Request) {
	// Обработка запроса в зависимости от типа переданного URL.
	moduls.InfoLog.Println("Получен GET-запрос")

	buf := make(map[string]string)

	for k, m := range items {
		buf[k] = m.Name
	}

	fmt.Println(buf)
	// Если в пути обращения GET - "/items/" , возвращаем список всех элементов.
	sendJSONResponse(w, http.StatusOK, buf)
}

// handleGET - обработчик для HTTP-запросов методом GET.
func handleGETid(w http.ResponseWriter, r *http.Request) {
	// Обработка запроса в зависимости от типа переданного URL.
	moduls.InfoLog.Println("Получен GET-запрос")

	// Если в пути обращения GET - "/items/{item_id}/", возвращаем соответствующий элемент.
	itemID := r.PathValue("id")

	if item, ok := items[itemID]; ok {
		sendJSONResponse(w, http.StatusOK, item)
	} else {
		// Если элемент с указанным ID не существует, возвращаем ошибку "Not Found".
		http.NotFound(w, r)
	}
}

// handlePOST - обработчик для HTTP-запросов методом POST.
func handlePOST(w http.ResponseWriter, r *http.Request) {
	moduls.InfoLog.Println("Получен POST-запрос")

	// Декодирование JSON-тела запроса в новый элемент.
	var newItem moduls.Item
	err := decodeJSONBody(r.Body, &newItem)
	if err != nil {
		// Если произошла ошибка при декодировании JSON, возвращаем ошибку "Bad Request".
		http.Error(w, "Некорректный формат JSON", http.StatusBadRequest)
		return
	}

	// Проверка, что имя нового элемента не пустое.
	if newItem.Name == "" {
		// Если имя пустое, возвращаем ошибку "Bad Request".
		http.Error(w, "Название не может быть пустым", http.StatusBadRequest)
		return
	}

	//Проверяет длинну и допустимость вводимых данных
	if moduls.Sanitize(newItem.Name) {
		http.Error(w, "Недопустимые символы", http.StatusBadRequest)
		return
	}

	if moduls.Length(newItem.Name) {
		http.Error(w, "Недопустимая длинна (более 20 символов)", http.StatusBadRequest)
		return
	}
	// Генерация уникального ID и добавление нового элемента в карту.
	newItem.ID = GenerateID()
	items[newItem.ID] = &newItem

	// Отправка JSON-ответа с созданным элементом и статусом "Created".
	sendJSONResponse(w, http.StatusCreated, newItem)
}

// handlePUT - обработчик для HTTP-запросов методом PUT.
func handlePUT(w http.ResponseWriter, r *http.Request) {
	moduls.InfoLog.Println("Получен PUT-запрос")
	// Извлечение ID элемента из URL.
	itemID := r.PathValue("id")
	if item, ok := items[itemID]; ok {
		// Если элемент существует, декодирование JSON-тела запроса в обновленный элемент.
		var updatedItem moduls.Item
		err := decodeJSONBody(r.Body, &updatedItem)
		if err != nil {
			// Если произошла ошибка при декодировании JSON, возвращаем ошибку "Bad Request".
			http.Error(w, "Некорректный формат JSON", http.StatusBadRequest)
			return
		}

		// Проверка, что имя обновленного элемента не пустое.
		if updatedItem.Name == "" {
			// Если имя пустое, возвращаем ошибку "Bad Request".
			http.Error(w, "Название не может быть пустым", http.StatusBadRequest)
			return
		}

		//Проверяет длинну и допустимость вводимых данных
		if moduls.Sanitize(updatedItem.Name) {
			http.Error(w, "Недопустимые символы", http.StatusBadRequest)
			return
		}

		if moduls.Length(updatedItem.Name) {
			http.Error(w, "Недопустимая длинна (более 20 символов)", http.StatusBadRequest)
			return
		}

		// Обновление имени элемента и отправка JSON-ответа с обновленным элементом.
		item.Name = updatedItem.Name
		sendJSONResponse(w, http.StatusOK, item)
	} else {
		// Если элемент с указанным ID не существует, возвращаем ошибку "Not Found".
		http.NotFound(w, r)
	}
}

// handleDELETE - обработчик для HTTP-запросов методом DELETE.
func handleDELETE(w http.ResponseWriter, r *http.Request) {
	moduls.InfoLog.Println("Получен DELETE-запрос")

	// Извлечение ID элемента из URL.
	itemID := r.PathValue("id")

	if moduls.Sanitize(itemID) {
		http.Error(w, "Недопустимые символы", http.StatusBadRequest)
		return
	}

	if moduls.Length(itemID) {
		http.Error(w, "Недопустимая длинна (более 20 символов)", http.StatusBadRequest)
		return
	}

	if item, ok := items[itemID]; ok {
		// Если элемент существует, удаление элемента из карты.
		delete(items, item.ID)
		// Возвращение статуса "No Content" (204) в ответе.
		w.WriteHeader(http.StatusNoContent)
	} else {
		// Если элемент с указанным ID не существует, возвращаем ошибку.
		http.NotFound(w, r)
	}
}

// sendJSONResponse - устанавливает заголовки ответа и кодирует данные в формате JSON для отправки.
func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	// Установка заголовка "Content-Type" как "application/json".
	w.Header().Set("Content-Type", "application/json")
	// Установка кода состояния ответа.
	w.WriteHeader(statusCode)

	// Кодирование данных в формат JSON и отправка в тело ответа.
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		// Если произошла ошибка при кодировании JSON, возвращаем ошибку
		moduls.ErrorLog.Println("Ошибка при кодировании JSON:", err)
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
		return
	}
}

// decodeJSONBody - декодирует JSON-тело запроса в структуру данных.
func decodeJSONBody(body io.Reader, v interface{}) error {
	// Используется json.NewDecoder для декодирования JSON из тела запроса.
	return json.NewDecoder(body).Decode(v)
}

// GenerateID - генерирует уникальный ID для элемента
func GenerateID() string {
	return uuid.New().String()[:8]
}
