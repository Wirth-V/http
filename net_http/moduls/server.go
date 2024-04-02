package moduls

import (
	"encoding/json"
	"flag"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

// items - глобальная переменная, представляющая соотношение элементов по их уникальным ID.
var items = make(map[string]*Item)

func Server(req *flag.FlagSet, host *string, port *string) {

	InfoLog.Println("Сервер запущен.")
	InfoLog.Printf("Хост:%s Порт:%s", *host, *port)

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
		ErrorLog.Fatal("Ошибка запуска сервера:", err)
	}
}

// handleGET - обработчик для HTTP-запросов методом GET.
func handleGET(w http.ResponseWriter, r *http.Request) {
	// Обработка запроса в зависимости от типа переданного URL.
	InfoLog.Println("Получен GET-запрос")

	slice := make([]*Item, 0)

	for _, m := range items {
		slice = append(slice, m)
	}

	//if item, ok := items[itemID]; ok {
	//sendJSONResponse(w, http.StatusOK, item)

	// Если в пути обращения GET - "/items/" , возвращаем список всех элементов.
	sendJSONResponse(w, http.StatusOK, slice)
}

// handleGET - обработчик для HTTP-запросов методом GET.
func handleGETid(w http.ResponseWriter, r *http.Request) {
	// Обработка запроса в зависимости от типа переданного URL.
	InfoLog.Println("Получен GET-запрос")

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
	InfoLog.Println("Получен POST-запрос")

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

	//Проверяет длинну и допустимость вводимых данных
	if Sanitize(newItem.Name) {
		http.Error(w, "Недопустимые символы", http.StatusBadRequest)
		return
	}

	if Length(newItem.Name) {
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
	InfoLog.Println("Получен PUT-запрос")
	// Извлечение ID элемента из URL.
	itemID := r.PathValue("id")
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

		//Проверяет длинну и допустимость вводимых данных
		if Sanitize(updatedItem.Name) {
			http.Error(w, "Недопустимые символы", http.StatusBadRequest)
			return
		}

		if Length(updatedItem.Name) {
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
	InfoLog.Println("Получен DELETE-запрос")

	// Извлечение ID элемента из URL.
	itemID := r.PathValue("id")

	if Sanitize(itemID) {
		http.Error(w, "Недопустимые символы", http.StatusBadRequest)
		return
	}

	if Length(itemID) {
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
		ErrorLog.Println("Ошибка при кодировании JSON:", err)
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
