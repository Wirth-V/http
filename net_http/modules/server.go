package modules

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
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
	err := Control(connString)
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

	items, err := GetItem(r.Context())

	if err != nil {
		ErrorLog.Println("error in GET request:", err)
		http.Error(w, "error", http.StatusInternalServerError)
		return
	}

	// Возвращаем список всех элементов.
	sendJSONResponse(w, items)
}

func handleGETid(w http.ResponseWriter, r *http.Request) {

	InfoLog.Println("A GET request was received")

	itemID := r.PathValue("id")

	// Запрос данных из таблицы по ID
	items, err := GetItemId(r.Context(), itemID)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.NotFound(w, r)
		} else {
			ErrorLog.Println("error querying database:", err)
			http.Error(w, "error querying database", http.StatusInternalServerError)
		}
		return
	}

	sendJSONResponse(w, items)

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

	err = PostItem(r.Context(), &newItem)

	if err != nil {
		ErrorLog.Println("error in POST request,", err)
		http.Error(w, "error", http.StatusInternalServerError)
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

	var contol bool
	err, contol = PutItem(r.Context(), &updatedItem, itemID)

	if contol {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		ErrorLog.Println("error in PUT request:", err)
		http.Error(w, "error", http.StatusInternalServerError)
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

	var contol bool
	err, contol = DeleteItem(r.Context(), itemID)

	if contol {
		http.NotFound(w, r)
		return
	}

	if err != nil {
		ErrorLog.Println("error in PUT request:", err)
		http.Error(w, "error", http.StatusInternalServerError)
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
