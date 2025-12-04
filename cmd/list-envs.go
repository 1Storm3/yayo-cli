package cmd

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	_ "github.com/mutecomm/go-sqlcipher/v4"
	"github.com/spf13/cobra"
)

var listEnvCmd = &cobra.Command{
	Use:   "list-envs",
	Short: "Выводит список ENV переменных проекта",
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("p")
		service, _ := cmd.Flags().GetString("s")

		if project == "" {
			return fmt.Errorf("необходимо указать --p")
		}

		reader := bufio.NewReader(os.Stdin)
		password, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("не удалось прочитать пароль: %w", err)
		}
		password = strings.TrimSpace(password)

		dbPath := fmt.Sprintf("%s/.yayo/%s/yayo.db", os.Getenv("HOME"), project)
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			return fmt.Errorf("база проекта '%s' не найдена", project)
		}

		sqlDB, err := sql.Open("sqlite3", fmt.Sprintf("%s?_pragma_key=%s", dbPath, password))

		if err != nil {
			return fmt.Errorf("не удалось открыть БД: %w", err)
		}
		defer sqlDB.Close()

		query := "SELECT key, value, service FROM envs"
		var argsQuery []interface{}
		if service != "" {
			query += " WHERE service = ?"
			argsQuery = append(argsQuery, service)
		}

		rows, err := sqlDB.Query(query, argsQuery...)
		if err != nil {
			return fmt.Errorf("ошибка запроса к БД: %w", err)
		}
		defer rows.Close()

		type EnvItem struct {
			Key     string `json:"key"`
			Value   string `json:"value"`
			Service string `json:"service"`
		}

		var envs []EnvItem
		for rows.Next() {
			var e EnvItem
			if err := rows.Scan(&e.Key, &e.Value, &e.Service); err != nil {
				return fmt.Errorf("ошибка чтения строки: %w", err)
			}
			envs = append(envs, e)
		}

		output, err := json.MarshalIndent(envs, "", "  ")
		if err != nil {
			return fmt.Errorf("ошибка маршалинга JSON: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	listEnvCmd.Flags().String("p", "", "Название проекта")
	listEnvCmd.Flags().String("s", "", "Фильтр по сервису (необязательно)")
	RootCmd.AddCommand(listEnvCmd)
}
