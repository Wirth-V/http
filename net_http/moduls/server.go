package moduls

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

// items - глобальная переменная, представляющая соотношение элементов по их уникальным ID.
var items = make(map[string]*Item)
var connFerst *pgx.Conn
var Table string

func Server(req *flag.FlagSet, host *string, port *string, db *string, table *string) {
	connString := "postgres://server:198416@localhost:6667/" + *db
	Table = *table

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

	// Пример выполнения SQL-запроса и получения результата
	// Запрос данных из таблицы
	rowsFerst, err := connFerst.Query(context.Background(), "SELECT * FROM "+Table)
	if err != nil {
		fmt.Println("Error querying database:", err)
		return
	}
	defer rowsFerst.Close()

	var id string
	var name string
	var newItem Item
	// Итерация по результатам запроса и вывод содержимого
	for rowsFerst.Next() {
		// Поменяйте тип данных на соответствующий вашей таблице
		// Пример чтения данных из строки
		err_bd := rowsFerst.Scan(&id, &name) // Замените переменные на соответствующие вашей таблице
		if err_bd != nil {
			fmt.Println("Error scanning row:", err_bd)
			return
		}
		newItem.ID = id
		newItem.Name = name
		// Вывод содержимого строки
		items[id] = &newItem // Измените вывод на соответствующий вашей таблице
	}

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
	err_bd := http.ListenAndServe(strings.Join([]string{*host, *port}, ":"), nil)
	if err_bd != nil {
		ErrorLog.Fatal("Ошибка запуска сервера:", err_bd)
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

	rows, err := connFerst.Query(context.Background(), "INSERT INTO "+Table+" (id, name) VALUES ($1, $2)", newItem.ID, newItem.Name)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

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

		rows, err := connFerst.Query(context.Background(), "UPDATE $3 SET name = $1 WHERE id = $2", updatedItem.Name, itemID, Table)
		if err != nil {
			fmt.Println("Error executing query:", err)
			return
		}
		defer rows.Close()

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

		rows, err := connFerst.Query(context.Background(), "DELETE FROM $2 WHERE id = $1", itemID, Table)
		if err != nil {
			fmt.Println("Error executing query:", err)
			return
		}
		defer rows.Close()

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

// Проверяет наличие бд, если его нет, то создет нужное бд и таблицу
func db_control(connString string, Table string) error {
	//разбора строки, возвращает структуру
	connConfig, err := pgx.ParseConfig(connString)
	//pgx.ConnConfig, содержащую параметры соединения
	if err != nil {
		return err
	}

	//извлекается имя подключаемой бд
	dbname := connConfig.Database

	//меняет название подключаемой бд
	connConfig.Database = "postgres"
	//Соединение с бд, чтобы проверить ее существование
	conn_db, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return err
	}
	defer conn_db.Close(context.Background())

	//проверка существования бд
	var exists_db bool
	//Выполняется запрос к системной таблице pg_database, чтобы проверить существует ли бд
	err = conn_db.QueryRow(context.Background(), "SELECT EXISTS (SELECT FROM pg_database WHERE datname = $1)", dbname).Scan(&exists_db)
	//conn.QueryRow Выполнения запроса о существовании, возвращает true, если существует, или false, если нет.
	if err != nil {
		return err
	}

	//создание бд
	if !exists_db {
		//сама команда создания отсутсвующей бд
		_, err = conn_db.Exec(context.Background(), "CREATE DATABASE "+dbname)

		if err != nil {
			return err
		}

		// Установка подключения к созданной бд
		connConfig.Database = dbname
		conn_db, err = pgx.ConnectConfig(context.Background(), connConfig)
		if err != nil {
			return err
		}
		defer conn_db.Close(context.Background())

		// Создание таблицы
		_, err = conn_db.Exec(context.Background(), "CREATE TABLE "+Table+"(id VARCHAR(8), name VARCHAR(30))")
		if err != nil {
			return err
		}

	} else {
		connConfig.Database = dbname
		conn_table, err := pgx.ConnectConfig(context.Background(), connConfig)
		if err != nil {
			return err
		}
		defer conn_table.Close(context.Background())

		var exists_table bool
		//Выполняется запрос к системной таблице pg_database, чтобы проверить существует ли бд
		err = conn_table.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", Table).Scan(&exists_table)
		//conn.QueryRow Выполнения запроса о существовании, возвращает true, если существует, или false, если нет.
		if err != nil {
			return err
		}

		if !exists_table {

			// Создание таблицы
			_, err = conn_table.Exec(context.Background(), "CREATE TABLE "+Table+"(id VARCHAR(8), name VARCHAR(30))")
			if err != nil {
				return err
			}
		}

	}
	return nil
}
