package main

import (
	"fmt"
	"net"
)

func main() {
	// Подключаемся к серверу
	conn, err := net.Dial("tcp", "localhost:8080") //Функция net.Dial возвращает тип net.Conn, который реализует интерфейсы
	//io.Reader и io.Writer , что позволяет в данный объект записывать и читать данные
	if err != nil {
		fmt.Println("Ошибка подключения к серверу:", err)
		return
	}
	defer conn.Close()

	// Отправляем GET запрос
	getRequest := "GET / HTTP/1.1\r\n\r\n"
	conn.Write([]byte(getRequest))

	// Буфер для чтения ответа
	buffer := make([]byte, 1024)
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка чтения данных:", err)
		return
	}

	// Выводим ответ от сервера
	fmt.Printf("Ответ на GET запрос: %s\n\n", string(buffer))

	// Отправляем POST запрос
	postRequest := "POST / HTTP/1.1\r\n\r\n"
	conn.Write([]byte(postRequest))

	// Считываем ответ на POST запрос
	_, err = conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка чтения данных POST:", err)
		return
	}

	// Выводим ответ от сервера
	fmt.Println("Ответ на POST запрос:", string(buffer))
}
