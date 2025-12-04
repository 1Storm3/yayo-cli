package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/1Storm3/yayo-cli/internal/config"
	"github.com/1Storm3/yayo-cli/internal/ssh"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Загружает локальный .env на удалённую машину через SSH",
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		service, _ := cmd.Flags().GetString("service")
		project, _ := cmd.Flags().GetString("project")
		filePath, _ := cmd.Flags().GetString("file")

		fmt.Print("Введите пароль хранилища: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("не удалось прочитать пароль: %w", err)
		}
		password := string(passwordBytes)

		envMap, err := config.LoadEnvFile(filePath)
		if err != nil {
			return err
		}

		type EnvItem struct {
			Key     string `json:"key"`
			Value   string `json:"value"`
			Service string `json:"service"`
		}

		envItems := make([]EnvItem, 0, len(envMap))
		for key, value := range envMap {
			envItems = append(envItems, EnvItem{
				Key:     key,
				Value:   value,
				Service: service,
			})
		}

		data, err := json.Marshal(envItems)
		if err != nil {
			return err
		}

		remoteCmd := fmt.Sprintf(
			"yayo-cli add-bulk --project %s --password '%s'",
			project, password,
		)
		fmt.Println("JSON для отправки:", string(data))

		output, err := ssh.RunSSHWithStdin(host, remoteCmd, string(data))
		if err != nil {
			return fmt.Errorf("ошибка при пуше: %w", err)
		}

		fmt.Println(output)
		fmt.Println("Все переменные успешно загружены!")
		return nil
	},
}

func init() {
	RootCmd.AddCommand(pushCmd)
	pushCmd.Flags().String("host", "", "SSH host, например root@1.2.3.4")
	pushCmd.Flags().String("project", "", "Название проекта")
	pushCmd.Flags().String("service", "", "Название сервиса")
	pushCmd.Flags().String("file", "", "Путь к .env файлу")
}
