package modules

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net_http/repository"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var connFerst *pgx.Conn
var Table string
var shutdownSignal = make(chan os.Signal, 1)

func handleInterrupt(restServer *http.Server) {
	// Обработка сигнала SIGTERM для грациозного завершения сервера
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-shutdownSignal
		fmt.Printf("\n")
		InfoLog.Println("Received SIGTERM. Shutting down gracefully...")

		if err := connFerst.Close(context.Background()); err != nil {
			ErrorLog.Printf("error closing database connection: %v\n", err)
		}
		// грациозной остановки HTTP-сервера. Он позволяет серверу завершить
		// обработку уже полученных запросов и корректно закрыть все открытые сетевые соединения.
		if err := restServer.Shutdown(context.Background()); err != nil {
			ErrorLog.Printf("error shutting down server: %v\n", err)
		}
	}()
}

func Server(req *flag.FlagSet, host string, port string, connString string, table string) error {
	if req == nil {
		return fmt.Errorf("ettempt to pass nil to the 'req' variable")

	}

	Table = table

	// Создание и проверка наличия бд и таблицы, указанных пользователем
	err := repository.Control(connString, Table)
	if err != nil {
		return fmt.Errorf("error checking database existence, %v", err)

	}

	// Установка соединения с базой данных
	connFerst, err = pgx.Connect(context.Background(), connString)
	if err != nil {
		return fmt.Errorf("unable to connect to database, %v", err)

	}
	defer connFerst.Close(context.Background())

	InfoLog.Println("Server start.")
	InfoLog.Printf("Host:%s Port:%s", host, port)

	router := http.NewServeMux()

	restServer := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf("%v:%v", host, port),
		WriteTimeout: time.Hour * 3,
		ReadTimeout:  time.Hour * 3,
	}

	// Регистрация обработчиков запросов для различных путей
	router.HandleFunc("GET /items/", handleGET)
	router.HandleFunc("GET /items/{id}/", handleGETid)
	router.HandleFunc("POST /items/", handlePOST)
	router.HandleFunc("PUT /items/{id}/", handlePUT)
	router.HandleFunc("DELETE /items/{id}/", handleDELETE)

	handleInterrupt(restServer)

	// Запуск веб-сервера
	err = restServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server startup error: %v", err)
	}

	return nil
}

func handleGET(w http.ResponseWriter, r *http.Request) {

	InfoLog.Println("A GET request was received")

	// Начало транзакции
	tx, err := connFerst.BeginTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		ErrorLog.Println("error beginning transaction:", err)
		http.Error(w, "error beginning transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	// Запрос данных из таблицы
	rows, err := tx.Query(r.Context(), "SELECT * FROM "+Table)
	if err != nil {
		ErrorLog.Println("eror querying database for GET request,", err)
		http.Error(w, "error querying database", http.StatusInternalServerError)
		return

	}
	defer rows.Close()

	var items []*Item
	var id string
	var name string

	// Итерация по результатам запроса и добавление данных в  срез items
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			ErrorLog.Println("error scanning row:", err)
			http.Error(w, "error scanning row", http.StatusInternalServerError)
			return
		}
		items = append(items, &Item{ID: id, Name: name})
	}

	// Обеспечивает нужный формат возврата данных для пустой таблице
	// (Делает так, что бы вернулся не `nil`, а `{"id":"", "name":""} `)
	if items == nil {
		items = append(items, &Item{ID: "", Name: ""})
	}

	// Коммит транзакции
	err = tx.Commit(r.Context())
	if err != nil {
		ErrorLog.Println("error committing transaction:", err)
		http.Error(w, "error committing transaction", http.StatusInternalServerError)
		return
	}

	// Возвращаем список всех элементов.
	sendJSONResponse(w, items)
}

func handleGETid(w http.ResponseWriter, r *http.Request) {

	InfoLog.Println("A GET request was received")

	itemID := r.PathValue("id")

	// Запрос данных из таблицы по ID
	var name string

	// Начало транзакции
	tx, err := connFerst.BeginTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		ErrorLog.Println("error beginning transaction:", err)
		http.Error(w, "error beginning transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	err = tx.QueryRow(r.Context(), "SELECT name FROM "+Table+" WHERE id = $1", itemID).Scan(name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			ErrorLog.Println("error querying database:", err)
			http.Error(w, "error querying database", http.StatusInternalServerError)
		}
		return
	}

	// Коммит транзакции
	err = tx.Commit(r.Context())
	if err != nil {
		ErrorLog.Println("error committing transaction:", err)
		http.Error(w, "error committing transaction", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, &Item{ID: itemID, Name: name})

}

func handlePOST(w http.ResponseWriter, r *http.Request) {

	InfoLog.Println("A POST request was received")

	var newItem Item
	err := decodeJSONBody(r.Body, &newItem)
	if err != nil {
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	if newItem.Name == "" {
		http.Error(w, "the name cannot be empty", http.StatusBadRequest)
		return
	}

	err = check(newItem.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("when creating a new element %s", err), http.StatusBadRequest)
		return
	}

	tx, err := connFerst.BeginTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		ErrorLog.Println("error beginning transaction:", err)
		http.Error(w, "error beginning transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	// Генерация уникального ID и добавление нового элемента в карту.
	newItem.ID = GenerateID()

	_, err = tx.Exec(r.Context(), "INSERT INTO "+Table+" (id, name) VALUES ($1, $2)", newItem.ID, newItem.Name)
	if err != nil {
		ErrorLog.Println("error executing query,", err)
		http.Error(w, "error executing query", http.StatusInternalServerError)
		return
	}

	// Коммит транзакции
	err = tx.Commit(r.Context())
	if err != nil {
		ErrorLog.Println("error committing transaction:", err)
		http.Error(w, "error committing transaction", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, newItem)
}

func handlePUT(w http.ResponseWriter, r *http.Request) {

	InfoLog.Println("A PUT request was received")

	// Извлечение ID элемента из URL.
	itemID := r.PathValue("id")

	// Если элемент существует, декодирование JSON-тела запроса в обновленный элемент.
	var updatedItem Item
	err := decodeJSONBody(r.Body, &updatedItem)
	if err != nil {
		http.Error(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	if updatedItem.Name == "" {
		http.Error(w, "the name cannot be empty", http.StatusBadRequest)
		return
	}

	err = check(updatedItem.Name)
	if err != nil {
		http.Error(w, fmt.Sprintf("when updated a name of element %s", err), http.StatusBadRequest)
		return
	}

	// Начало транзакции
	tx, err := connFerst.BeginTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		ErrorLog.Println("error beginning transaction:", err)
		http.Error(w, "error beginning transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	// Проверка существования элемента в базе данных
	var count int
	err = tx.QueryRow(r.Context(), "SELECT COUNT(*) FROM "+Table+" WHERE id = $1", itemID).Scan(&count)
	if err != nil {
		ErrorLog.Println("error querying database:", err)
		http.Error(w, "error querying database", http.StatusInternalServerError)
		return
	}

	if count == Zero {
		http.NotFound(w, r)
		return
	}

	// Обновление данных элемента в базе данных
	_, err = tx.Exec(r.Context(), "UPDATE "+Table+" SET name = $1 WHERE id = $2", updatedItem.Name, itemID)
	if err != nil {
		ErrorLog.Println("error executing query:", err)
		http.Error(w, "error executing query", http.StatusInternalServerError)
		return
	}

	// Коммит транзакции
	err = tx.Commit(r.Context())
	if err != nil {
		ErrorLog.Println("error committing transaction:", err)
		http.Error(w, "error committing transaction", http.StatusInternalServerError)
		return
	}

	sendJSONResponse(w, &Item{ID: itemID, Name: updatedItem.Name})

}

func handleDELETE(w http.ResponseWriter, r *http.Request) {

	InfoLog.Println("A DELETE request was received")

	// Извлечение ID элемента из URL.
	itemID := r.PathValue("id")

	err := check(itemID)
	if err != nil {
		http.Error(w, fmt.Sprintf("when dealete element %s", err), http.StatusBadRequest)
		return
	}

	// Проверка существования элемента в базе данных
	var count int

	// Начало транзакции
	tx, err := connFerst.BeginTx(r.Context(), pgx.TxOptions{})
	if err != nil {
		ErrorLog.Println("error beginning transaction:", err)
		http.Error(w, "error beginning transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(r.Context())

	err = tx.QueryRow(r.Context(), "SELECT COUNT(*) FROM "+Table+" WHERE id = $1", itemID).Scan(&count)
	if err != nil {
		ErrorLog.Println("error querying database: ", err)
		http.Error(w, "error querying database", http.StatusInternalServerError)
		return
	}

	if count == Zero {
		http.NotFound(w, r)
		return
	}

	// Удаление элемента из базы данных
	_, err = tx.Exec(r.Context(), "DELETE FROM "+Table+" WHERE id = $1", itemID)
	if err != nil {
		ErrorLog.Println("error executing query:", err)
		http.Error(w, "error executing query", http.StatusInternalServerError)
		return
	}

	// Коммит транзакции
	err = tx.Commit(r.Context())
	if err != nil {
		ErrorLog.Println("error committing transaction:", err)
		http.Error(w, "error committing transaction", http.StatusInternalServerError)
		return
	}

}

// sendJSONResponse - устанавливает заголовки ответа и кодирует данные в формате JSON для отправки.
func sendJSONResponse(w http.ResponseWriter, data interface{}) {

	w.Header().Set("Content-Type", "application/json")

	// Кодирование данных в формат JSON и отправка в тело ответа.
	err := json.NewEncoder(w).Encode(data)

	if err != nil {
		ErrorLog.Println("error encoding JSON:", err)
		http.Error(w, "error encoding JSON", http.StatusInternalServerError)
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
