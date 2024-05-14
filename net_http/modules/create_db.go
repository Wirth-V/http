package modules

import (
	"context"

	"github.com/jackc/pgx/v5"
)

// Проверяет наличие бд, если его нет, то создет нужное бд и таблицу
func Control(connString string) error {
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

func GetItem(ctx context.Context) ([]*Item, error) {

	InfoLog.Println("A GET request was received")

	// Начало транзакции
	tx, err := connFerst.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Запрос данных из таблицы
	rows, err := tx.Query(ctx, "SELECT * FROM "+Table)
	if err != nil {
		return nil, err

	}
	defer rows.Close()

	var items []*Item
	var id string
	var name string

	// Итерация по результатам запроса и добавление данных в  срез items
	for rows.Next() {
		err := rows.Scan(&id, &name)
		if err != nil {
			return nil, err
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
		return nil, err
	}

	return items, nil
}

func GetItemId(ctx context.Context, itemID string) (*Item, error) {
	// Запрос данных из таблицы по ID
	var name string

	// Начало транзакции
	tx, err := connFerst.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "SELECT name FROM "+Table+" WHERE id = $1", itemID).Scan(&name)
	if err != nil {
		return nil, err
	}

	// Коммит транзакции
	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &Item{ID: itemID, Name: name}, nil
}

func PostItem(ctx context.Context, newItem *Item) error {

	tx, err := connFerst.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Генерация уникального ID и добавление нового элемента в карту.
	newItem.ID = GenerateID()

	_, err = tx.Exec(ctx, "INSERT INTO "+Table+" (id, name) VALUES ($1, $2)", newItem.ID, newItem.Name)
	if err != nil {
		return err
	}

	// Коммит транзакции
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func PutItem(ctx context.Context, updatedItem *Item, itemID string) (error, bool) {
	// Начало транзакции
	tx, err := connFerst.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err, false
	}
	defer tx.Rollback(ctx)

	// Проверка существования элемента в базе данных
	var count int
	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM "+Table+" WHERE id = $1", itemID).Scan(&count)
	if err != nil {
		return err, false
	}

	if count == Zero {
		return nil, true
	}

	// Обновление данных элемента в базе данных
	_, err = tx.Exec(ctx, "UPDATE "+Table+" SET name = $1 WHERE id = $2", updatedItem.Name, itemID)
	if err != nil {
		return err, false
	}

	// Коммит транзакции
	err = tx.Commit(ctx)
	if err != nil {
		return err, false
	}

	return nil, false
}

func DeleteItem(ctx context.Context, itemID string) (error, bool) {
	// Проверка существования элемента в базе данных
	var count int

	// Начало транзакции
	tx, err := connFerst.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err, false
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "SELECT COUNT(*) FROM "+Table+" WHERE id = $1", itemID).Scan(&count)
	if err != nil {
		return err, false
	}

	if count == Zero {
		return nil, true
	}

	// Удаление элемента из базы данных
	_, err = tx.Exec(ctx, "DELETE FROM "+Table+" WHERE id = $1", itemID)
	if err != nil {
		return err, false
	}

	// Коммит транзакции
	err = tx.Commit(ctx)
	if err != nil {
		return err, false
	}
	return nil, false
}
