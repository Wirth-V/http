package modules

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Проверяет наличие бд, если его нет, то создет нужное бд и таблицу
func db_control(connString string, Table string) error {
	//разбора строки, возвращает структуру
	connConfig, err := pgx.ParseConfig(connString)
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
	err = conn_db.QueryRow(context.Background(), "SELECT EXISTS (SELECT FROM pg_database WHERE datname = $1)", dbname).Scan(&exists_db)
	if err != nil {
		return err
	}

	//создание бд
	if !exists_db {
		_, err = conn_db.Exec(context.Background(), "CREATE DATABASE "+dbname)

		if err != nil {
			return err
		}
	}

	connConfig.Database = dbname
	conn_table, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return err
	}
	defer conn_table.Close(context.Background())

	var exists_table bool
	err = conn_table.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = $1)", Table).Scan(&exists_table)

	if err != nil {
		return err
	}

	if !exists_table {
		// Создание таблицы
		_, err = conn_table.Exec(context.Background(), "CREATE TABLE "+Table+"(id VARCHAR(8) NOT NULL, name VARCHAR(30) NOT NULL)")
		if err != nil {
			return err
		}
	}

	return nil
}
