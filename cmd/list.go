package cmd

import (
	"fmt"
	"os"

	"golang.org/x/term"

	"github.com/1Storm3/yayo-cli/internal/ssh"
	"github.com/spf13/cobra"
)

var listEnvSSHCmd = &cobra.Command{
	Use:   "list",
	Short: "Выводит список ENV переменных проекта на удалённой машине через SSH",
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("host")
		project, _ := cmd.Flags().GetString("project")
		service, _ := cmd.Flags().GetString("service")

		if host == "" || project == "" {
			return fmt.Errorf("необходимо указать --host и --project")
		}

		fmt.Print("Введите пароль хранилища: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("не удалось прочитать пароль: %w", err)
		}
		password := string(passwordBytes)

		remoteCmd := fmt.Sprintf("yayo-cli list-r --project %s", project)
		if service != "" {
			remoteCmd += fmt.Sprintf(" --service %s", service)
		}

		output, err := ssh.RunSSHWithStdin(host, remoteCmd, password+"\n")
		if err != nil {
			return fmt.Errorf("ошибка при выполнении удалённой команды: %w", err)
		}

		fmt.Print(output)
		return nil
	},
}

func init() {
	listEnvSSHCmd.Flags().String("host", "", "SSH host, например root@1.2.3.4")
	listEnvSSHCmd.Flags().String("project", "", "Название проекта")
	listEnvSSHCmd.Flags().String("service", "", "Фильтр по сервису (необязательно)")
	rootCmd.AddCommand(listEnvSSHCmd)
}
