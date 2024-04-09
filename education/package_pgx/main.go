// urlExample := "postgres://username:password@localhost:5432/database_name"
// jdbc:postgresql://localhost:6667/postgres
// postgres://server:server@6667:5432/clients

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	// Строка подключения к базе данных
	connString := "host=localhost port=6667 user=server dbname=server password=198416 sslmode=disable"

	// Установка соединения с базой данных
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer conn.Close(context.Background())

	//fmt.Println(conn)

	var greeting string
	err = conn.QueryRow(context.Background(), "SELECT 'Hello, world!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(greeting)

	//INSERT INTO таблица(поле1, поле2) VALUES (значение1, значение2);
	rows, err := conn.Query(context.Background(), "INSERT INTO clients (id, name) VALUES ('666666', 'Vadim')")
	if err != nil {
		fmt.Println("Error executing query:", err)
		return
	}
	defer rows.Close()

	//fmt.Println(rows)

	// Установка соединения с базой данных
	connFerst, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer connFerst.Close(context.Background())

	//fmt.Println(connFerst)

	// Пример выполнения SQL-запроса и получения результата
	// Запрос данных из таблицы
	rowsFerst, err := connFerst.Query(context.Background(), "SELECT * FROM clients")
	if err != nil {
		fmt.Println("Error querying database:", err)
		return
	}
	defer rowsFerst.Close()

	fmt.Println(rowsFerst)

	// Итерация по результатам запроса и вывод содержимого
	for rowsFerst.Next() {
		var column1Type string
		var column2Type string // Поменяйте тип данных на соответствующий вашей таблице
		// Пример чтения данных из строки
		err := rowsFerst.Scan(&column1Type, &column2Type) // Замените переменные на соответствующие вашей таблице
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return
		}
		// Вывод содержимого строки
		fmt.Println("Column1:", column1Type, "Column2:", column2Type) // Измените вывод на соответствующий вашей таблице
	}
}
