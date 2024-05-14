package modules

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

// Проверяет наличие бд, если его нет, то создет нужное бд и таблицу
func Control(connString string, Table string) error {
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

func HandleGET(ctx context.Context, Table string, w http.ResponseWriter) ([]*Item, error) {

	var items []*Item
	var id string
	var name string

	// Начало транзакции
	tx, err := conn.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return items, fmt.Errorf("error in transaction: %v", err)
	}
	defer tx.Rollback(ctx)

	// Запрос данных из таблицы
	rows, err := tx.Query(ctx, "SELECT * FROM "+Table)
	if err != nil {
		return items, fmt.Errorf("eror querying database for GET request, %v", err)

	}
	defer rows.Close()

	// Итерация по результатам запроса и добавление данных в массив items
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return items, fmt.Errorf("error scanning row: %v", err)
		}
		items = append(items, &Item{ID: id, Name: name})
	}

	// Обеспечивает нужный формат возврата данных для пустой таблице
	// (Делает так, что бы вернулся не `nil`, а `{"id":"", "name":""} `)
	if items == nil {
		items = append(items, &Item{ID: "", Name: ""})
	}

	// Коммит транзакции
	err = tx.Commit(ctx)
	if err != nil {
		return items, fmt.Errorf("error committing transaction: %v", err)
	}

	// Возвращаем список всех элементов.
	sendJSONResponse(w, items)

	return items, err
}

func HandleGETid(ctx context.Context, Table string, itemID string, name *string) error {
	err := conn.QueryRow(ctx, "SELECT name FROM "+Table+" WHERE id = $1", itemID).Scan(name)
	return err
}

func HandlePOST(ctx context.Context, Table string, ID string, Name string) error {
	_, err := conn.Exec(ctx, "INSERT INTO "+Table+" (id, name) VALUES ($1, $2)", ID, Name)
	return err
}

func Сheck(ctx context.Context, Table string, itemID string, count *int) error {
	err := conn.QueryRow(ctx, "SELECT COUNT(*) FROM "+Table+" WHERE id = $1", itemID).Scan(count)
	return err
}

func HandlePUT(ctx context.Context, Table string, Name string, ID string) error {
	_, err := conn.Exec(ctx, "UPDATE "+Table+" SET name = $1 WHERE id = $2", Name, ID)
	return err
}

func HandleDELETE(ctx context.Context, Table string, itemID string) error {
	_, err := conn.Exec(ctx, "DELETE FROM "+Table+" WHERE id = $1", itemID)
	return err
}
