package modules

import (
	"context"

	"github.com/jackc/pgx/v5"
)

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
