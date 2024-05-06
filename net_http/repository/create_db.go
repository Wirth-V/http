package repository

import (
	"context"
	"net/http"

	"github.com/jackc/pgx/v5"
)

// Проверяет наличие бд, если его нет, то создет нужное бд и таблицу
func Db_control(connString string, Table string) error {
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

func DbHandleGET(tx pgx.Tx, r *http.Request, Table string) (pgx.Rows, error) {
	rows, err := tx.Query(r.Context(), "SELECT * FROM "+Table)
	return rows, err
}

func DbHandleGETid(tx pgx.Tx, r *http.Request, Table string, itemID string, name *string) error {
	err := tx.QueryRow(r.Context(), "SELECT name FROM "+Table+" WHERE id = $1", itemID).Scan(name)
	return err
}

func DbHandlePOST(tx pgx.Tx, r *http.Request, Table string, ID string, Name string) error {
	_, err := tx.Exec(r.Context(), "INSERT INTO "+Table+" (id, name) VALUES ($1, $2)", ID, Name)
	return err
}

func DbСheck(tx pgx.Tx, r *http.Request, Table string, itemID string, count *int) error {
	err := tx.QueryRow(r.Context(), "SELECT COUNT(*) FROM "+Table+" WHERE id = $1", itemID).Scan(count)
	return err
}

func DbHandlePUT(tx pgx.Tx, r *http.Request, Table string, Name string, ID string) error {
	_, err := tx.Exec(r.Context(), "UPDATE "+Table+" SET name = $1 WHERE id = $2", Name, ID)
	return err
}

func DbHandleDELETE(tx pgx.Tx, r *http.Request, Table string, itemID string) error {
	_, err := tx.Exec(r.Context(), "DELETE FROM "+Table+" WHERE id = $1", itemID)
	return err
}
