package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mutecomm/go-sqlcipher/v4"
)

func InitSQLCipher(path, password string) error {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_pragma_key=%s", path, password))
	if err != nil {
		return err
	}
	defer db.Close()

	return nil
}

func Migrations(path, password string) error {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("PRAGMA key = '%s';", password))
	if err != nil {
		return err
	}

	schema := `
CREATE TABLE IF NOT EXISTS envs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	key TEXT NOT NULL,
	value TEXT NOT NULL,
	service TEXT NOT NULL,
	UNIQUE(key, service)
);
	`

	_, err = db.Exec(schema)
	return err
}
