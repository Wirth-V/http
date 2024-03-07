package main

import (
	"fmt"
	"net/http"
)

func main() {
	// Обработчик для GET запросов
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Привет, это GET запрос!")
	})

	// Обработчик для POST запросов
	http.HandleFunc("/post", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			fmt.Fprintf(w, "Привет, это POST запрос!")
		} else {
			http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed) // вот так должно ругаться на get и post, которые обращаются не к своим items
		}
	})

	// Запуск HTTP сервера на порту 8080
	http.ListenAndServe(":8080", nil)
	fmt.Println("Listening on :8080...")
}
