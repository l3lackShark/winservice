package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type (
	Database interface {
		SendPayload(data []byte) (err error)
	}

	database struct {
		db *sql.DB
	}
)

func New(basePath string) (Database, error) {

	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("input path does not exist on the system")
	}

	//remove old db if exists. (non production state)
	dbAbsPath := filepath.Join(filepath.Dir(basePath), "foo.db")

	if err := os.Remove(dbAbsPath); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("os.Remove(): %w", err)
	}

	fmt.Println(dbAbsPath)
	db, err := sql.Open("sqlite3", dbAbsPath)
	if err != nil {
		return nil, fmt.Errorf("sql.Open(): %w", err)
	}

	//hardcoded table creation statement (non-production state)
	sqlStmt := `
	create table test (id integer not null primary key, date datetime, payload json);
	delete from test;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return nil, fmt.Errorf("db.Exec() on init statement: %w", err)
	}
	return &database{db: db}, nil
}

//SendPayload is non-Modular specific funtion to send payload to the database
func (w *database) SendPayload(data []byte) (err error) {
	tx, err := w.db.Begin()
	if err != nil {
		return fmt.Errorf("db.Begin(): %w", err)
	}

	stmt, err := tx.Prepare("insert into test(date, payload) values(CURRENT_TIMESTAMP, ?)")
	if err != nil {
		return fmt.Errorf("tx.Prepare(): %w", err)
	}
	defer func() {
		cerr := stmt.Close()
		if err == nil {
			err = cerr
		}
	}()

	_, err = stmt.Exec(string(data))
	if err != nil {
		return fmt.Errorf("stmt.Exec(): %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("tx.Commit(): %w", err)
	}

	return nil
}
