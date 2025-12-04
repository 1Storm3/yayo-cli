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

	return nil
}
