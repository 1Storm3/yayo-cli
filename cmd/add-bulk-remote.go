package cmd

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/mutecomm/go-sqlcipher/v4"
	"github.com/spf13/cobra"
)

var addBulkCmd = &cobra.Command{
	Use:   "add-bulk-r",
	Short: "Добавляет или обновляет несколько ENV переменных через stdin",
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("project")
		password, _ := cmd.Flags().GetString("password")

		if project == "" || password == "" {
			return fmt.Errorf("необходимо указать --project и --password")
		}

		dbPath := fmt.Sprintf("%s/.yayo/%s/yayo.db", os.Getenv("HOME"), project)
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			return fmt.Errorf("база проекта '%s' не найдена", project)
		}

		sqlDB, err := sql.Open("sqlite3", fmt.Sprintf("%s?_pragma_key=%s", dbPath, password))

		if err != nil {
			return fmt.Errorf("не удалось открыть БД: %w", err)
		}
		defer sqlDB.Close()

		_, err = sqlDB.Exec(fmt.Sprintf("PRAGMA key = '%s';", password))
		if err != nil {
			return fmt.Errorf("не удалось установить пароль БД: %w", err)
		}

		var envItems []struct {
			Key     string `json:"key"`
			Value   string `json:"value"`
			Service string `json:"service"`
		}

		if err := json.NewDecoder(os.Stdin).Decode(&envItems); err != nil {
			return fmt.Errorf("не удалось прочитать stdin: %w", err)
		}

		tx, err := sqlDB.Begin()
		if err != nil {
			return err
		}

		stmt, err := tx.Prepare(`
INSERT INTO envs (key, value, service)
VALUES (?, ?, ?)
ON CONFLICT(key, service) DO UPDATE SET value = excluded.value;
		`)
		if err != nil {
			return err
		}
		defer stmt.Close()

		for _, e := range envItems {
			if _, err := stmt.Exec(e.Key, e.Value, e.Service); err != nil {
				tx.Rollback()
				return fmt.Errorf("ошибка добавления переменной '%s': %w", e.Key, err)
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		fmt.Printf("Добавлено/обновлено %d переменных\n", len(envItems))
		return nil
	},
}

func init() {
	addBulkCmd.Flags().String("project", "", "Название проекта")
	addBulkCmd.Flags().String("password", "", "Пароль SQLCipher БД")
	rootCmd.AddCommand(addBulkCmd)
}
