package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/1Storm3/yayo-cli/internal/config"
	"github.com/1Storm3/yayo-cli/internal/db"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Инициализирует локальное зашифрованное хранилище ENV переменных",
	RunE: func(cmd *cobra.Command, args []string) error {

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Введите имя проекта: ")
		projectName, _ := reader.ReadString('\n')
		projectName = config.Clean(projectName)

		fmt.Print("Введите пароль для хранилища: ")
		password, _ := reader.ReadString('\n')
		password = config.Clean(password)

		projectDir := filepath.Join(config.BaseDir(), projectName)
		dbPath := filepath.Join(projectDir, "yayo.db")

		if _, err := os.Stat(projectDir); !os.IsNotExist(err) {
			return fmt.Errorf("проект '%s' уже инициализирован", projectName)
		}

		if err := os.MkdirAll(projectDir, 0755); err != nil {
			return err
		}

		if err := db.InitSQLCipher(dbPath, password); err != nil {
			return err
		}

		if err := config.Save(projectName, dbPath); err != nil {
			return err
		}

		cfg := map[string]string{
			"project": projectName,
			"db_path": dbPath,
		}

		cfgFile := filepath.Join(config.BaseDir(), projectName, "config.json")
		file, err := os.Create(cfgFile)
		if err != nil {
			return err
		}

		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(cfg); err != nil {
			return err
		}

		fmt.Println("Хранилище успешно инициализировано!!")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(initCmd)
}
