package main

import (
	"fmt"
	"net"
	"strings"
)

func handle(conn net.Conn) {

	// Буфер для чтения данных из соединения
	buffer := make([]byte, 1024)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка чтения данных:", err)
		return
	}
	defer conn.Close() // Закрываем соединение при завершении функции

	// Преобразование байтов в строку
	request := string(buffer)
	fmt.Println(request)

	// Обработка GET запроса
	if strings.HasPrefix(request, "GET") {
		response := "HTTP/1.1 200 OK\r\n\r\nПривет, это сервер! Вы сделали GET запрос."
		conn.Write([]byte(response))
	}

	// Обработка POST запроса
	if strings.HasPrefix(request, "POST") {
		response := "HTTP/1.1 200 OK\r\n\r\nПривет, это сервер! Вы сделали POST запрос."
		conn.Write([]byte(response))
	}
}

func main() {
	// Слушаем порт 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Ошибка прослушивания порта:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on :8080...")

	// Бесконечный цикл для принятия соединений
	for {
		conn, err := listener.Accept() // открываем порт
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}

		// Обработка соединения в отдельной горутине
		go handle(conn)
	}
}
