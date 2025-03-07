package database

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Database struct {
	sql *sqlx.DB
}

func New(driver string, conn string) *Database {
	db, err := sqlx.Connect(driver, conn)

	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	return &Database{
		sql: db,
	}
}

func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*5)
}

func (db *Database) Find(dest interface{}, table string, value int) (error, bool) {
	ctx, cancel := createContext()
	defer cancel()

	err := db.sql.GetContext(ctx, dest, "SELECT * FROM ? WHERE id = ?", table, value)

	if err != nil {
		return err, false
	}

	return err, true
}

func (db *Database) All(dest interface{}, table string, fields string) (error, bool) {
	ctx, cancel := createContext()
	defer cancel()

	err := db.sql.SelectContext(ctx, dest, fmt.Sprintf("SELECT %s FROM %s", fields, table))

	if err != nil {
		return err, false
	}

	return err, true
}

func (db *Database) Create(data map[string]interface{}, table string) (error, bool) {
	ctx, cancel := createContext()
	defer cancel()

	var keys []string
	var values []string
	for k := range data {
		keys = append(keys, k)
		values = append(values, ":"+k)
	}
	fields := strings.Join(keys, ",")

	_, err := db.sql.NamedExecContext(ctx, fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", table, fields, strings.Join(values, ",")), data)

	if err != nil {
		return err, false
	}

	return err, true
}

func (db *Database) Update(data map[string]interface{}, table string, id int) (error, bool) {
	ctx, cancel := createContext()
	defer cancel()

	q := fmt.Sprintf("Update %s set hit_count = $1 WHERE id = $2", table)
	_, err := db.sql.ExecContext(ctx, q, data["hit_count"], id)

	if err != nil {
		return err, false
	}

	return err, true
}

func (db *Database) Delete(table string, id int) (error, bool) {
	ctx, cancel := createContext()
	defer cancel()

	_, err := db.sql.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE id = $1", table), id)

	if err != nil {
		return err, false
	}

	return err, true
}

func (db *Database) GetRaw(query string, dest interface{}, args ...interface{}) error {
	ctx, cancel := createContext()
	defer cancel()

	return db.sql.GetContext(ctx, dest, query, args...)
}
