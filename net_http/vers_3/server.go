package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Item представляет структуру данных для элементов списка.
type Item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// items - глобальная переменная, представляющая соотношение элементов по их уникальным ID.
var items = make(map[string]*Item)

func main() {
	// Регистрация обработчика запросов для пути "/items/".
	http.HandleFunc("/items/", handleRequest)

	// Запуск веб-сервера на порту 8080.
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

// handleRequest - обработчик входящих HTTP-запросов.
func handleRequest(w http.ResponseWriter, r *http.Request) {
	// Обработка запроса в зависимости от метода HTTP.
	switch r.Method {
	case http.MethodGet:
		handleGET(w, r)
	case http.MethodPost:
		handlePOST(w, r)
	case http.MethodPut:
		handlePUT(w, r)
	case http.MethodDelete:
		handleDELETE(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// handleGET - обработчик для HTTP-запросов методом GET.
func handleGET(w http.ResponseWriter, r *http.Request) {

	// Обработка запроса в зависимости от типа переданного URL.
	if r.URL.Path == "/items/" {
		// Если в пути обращения GET - "/items/" , возвращаем список всех элементов.
		sendJSONResponse(w, http.StatusOK, items)
	} else {
		// Если в пути обращения GET - "/items/{item_id}/", возвращаем соответствующий элемент.
		itemID := r.URL.Path[len("/items/"):]
		if item, ok := items[itemID]; ok {
			sendJSONResponse(w, http.StatusOK, item)
		} else {
			// Если элемент с указанным ID не существует, возвращаем ошибку "Not Found".
			http.NotFound(w, r)
		}
	}
}

// handlePOST - обработчик для HTTP-запросов методом POST.
func handlePOST(w http.ResponseWriter, r *http.Request) {

	// Декодирование JSON-тела запроса в новый элемент.
	var newItem Item
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

	// Генерация уникального ID и добавление нового элемента в карту.
	newItem.ID = generateID()
	items[newItem.ID] = &newItem

	// Отправка JSON-ответа с созданным элементом и статусом "Created".
	sendJSONResponse(w, http.StatusCreated, newItem)
}

// handlePUT - обработчик для HTTP-запросов методом PUT.
func handlePUT(w http.ResponseWriter, r *http.Request) {

	// Извлечение ID элемента из URL.
	itemID := r.URL.Path[len("/items/"):]
	if item, ok := items[itemID]; ok {
		// Если элемент существует, декодирование JSON-тела запроса в обновленный элемент.
		var updatedItem Item
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

	// Извлечение ID элемента из URL.
	itemID := r.URL.Path[len("/items/"):]
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
		log.Println("Ошибка при кодировании JSON:", err)
		http.Error(w, "Ошибка при кодировании JSON", http.StatusInternalServerError)
		return
	}
}

// decodeJSONBody - декодирует JSON-тело запроса в структуру данных.
func decodeJSONBody(body io.Reader, v interface{}) error {
	// Используется json.NewDecoder для декодирования JSON из тела запроса.
	return json.NewDecoder(body).Decode(v)
}

// generateID - генерирует уникальный ID для элемента на основе текущего количества элементов.
func generateID() string {
	return fmt.Sprintf("%d", len(items)+1)
}
