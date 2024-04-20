package modules

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var connFerst *pgx.Conn
var Table string

func Server(req *flag.FlagSet, host string, port string, db string, table string) {
	if req == nil {
		fmt.Println("Attempt to pass nil to the 'req' variable")
		return
	}

	connString := "postgres://server:198416@localhost:6667/" + db
	Table = table

	err := db_control(connString, Table)
	if err != nil {
		fmt.Println("Error checking database existence:", err)
		return
	}

	// postgres://server:198416@localhost:6667/
	// host=localhost port=6667 user=server dbname=server password=198416 sslmode=disable

	// Установка соединения с базой данных
	connFerst, err = pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer connFerst.Close(context.Background())

	InfoLog.Println("Сервер запущен.")
	InfoLog.Printf("Хост:%s Порт:%s", host, port)

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
	err_bd := http.ListenAndServe(strings.Join([]string{host, port}, ":"), nil)
	if err_bd != nil {
		ErrorLog.Fatal("Ошибка запуска сервера:", err_bd)
	}
}

// handleGET - обработчик для HTTP-запросов методом GET.
func handleGET(w http.ResponseWriter, r *http.Request) {
	// Обработка запроса в зависимости от типа переданного URL.
	InfoLog.Println("Получен GET-запрос")

	// Запрос данных из таблицы
	rows, err := connFerst.Query(context.Background(), "SELECT * FROM "+Table)
	if err != nil {
		fmt.Println("Error querying database:", err)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []*Item
	var id string
	var name string
	// Итерация по результатам запроса и добавление данных в массив items
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		items = append(items, &Item{ID: id, Name: name})
	}

	// Если в пути обращения GET - "/items/" , возвращаем список всех элементов.
	sendJSONResponse(w, http.StatusOK, items)
}

// handleGET - обработчик для HTTP-запросов методом GET.
func handleGETid(w http.ResponseWriter, r *http.Request) {
	// Обработка запроса в зависимости от типа переданного URL.
	InfoLog.Println("Получен GET-запрос")

	// Если в пути обращения GET - "/items/{item_id}/", возвращаем соответствующий элемент.
	itemID := r.PathValue("id")

	// Запрос данных из таблицы по ID
	var name string
	err := connFerst.QueryRow(context.Background(), "SELECT name FROM "+Table+" WHERE id = $1", itemID).Scan(&name)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.NotFound(w, r)
		} else {
			fmt.Println("Error querying database:", err)
			http.Error(w, "Error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Отправка JSON-ответа с данными из базы данных
	sendJSONResponse(w, http.StatusOK, &Item{ID: itemID, Name: name})

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

	_, err = connFerst.Exec(context.Background(), "INSERT INTO "+Table+" (id, name) VALUES ($1, $2)", newItem.ID, newItem.Name)
	if err != nil {
		fmt.Println("Error executing query:", err)
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}

	// Отправка JSON-ответа с созданным элементом и статусом "Created".
	sendJSONResponse(w, http.StatusCreated, newItem)
}

// handlePUT - обработчик для HTTP-запросов методом PUT.
func handlePUT(w http.ResponseWriter, r *http.Request) {
	InfoLog.Println("Получен PUT-запрос")
	// Извлечение ID элемента из URL.
	itemID := r.PathValue("id")

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

	// Проверка существования элемента в базе данных
	var count int
	err = connFerst.QueryRow(context.Background(), "SELECT COUNT(*) FROM "+Table+" WHERE id = $1", itemID).Scan(&count)
	if err != nil {
		fmt.Println("Error querying database:", err)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.NotFound(w, r)
		return
	}

	// Обновление данных элемента в базе данных
	_, err = connFerst.Exec(context.Background(), "UPDATE "+Table+" SET name = $1 WHERE id = $2", updatedItem.Name, itemID)
	if err != nil {
		fmt.Println("Error executing query:", err)
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}

	// Отправка JSON-ответа с обновленным элементом.
	sendJSONResponse(w, http.StatusOK, &Item{ID: itemID, Name: updatedItem.Name})

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

	// Проверка существования элемента в базе данных
	var count int
	err := connFerst.QueryRow(context.Background(), "SELECT COUNT(*) FROM "+Table+" WHERE id = $1", itemID).Scan(&count)
	if err != nil {
		fmt.Println("Error querying database:", err)
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}

	if count == 0 {
		http.NotFound(w, r)
		return
	}

	// Удаление элемента из базы данных
	_, err = connFerst.Exec(context.Background(), "DELETE FROM "+Table+" WHERE id = $1", itemID)
	if err != nil {
		fmt.Println("Error executing query:", err)
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}

	// Возвращение статуса "No Content" (204) в ответе.
	w.WriteHeader(http.StatusNoContent)
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
