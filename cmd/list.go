package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"

	"golang.org/x/term"

	"github.com/1Storm3/yayo-cli/internal/ssh"
	"github.com/spf13/cobra"
)

var listEnvSSHCmd = &cobra.Command{
	Use:   "list",
	Short: "Выводит список ENV переменных проекта на удалённой машине через SSH",
	RunE: func(cmd *cobra.Command, args []string) error {
		host, _ := cmd.Flags().GetString("h")
		project, _ := cmd.Flags().GetString("p")
		service, _ := cmd.Flags().GetString("s")

		if host == "" || project == "" {
			return fmt.Errorf("необходимо указать --h и --p")
		}

		fmt.Print("Введите пароль хранилища: ")
		passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return fmt.Errorf("не удалось прочитать пароль: %w", err)
		}
		password := string(passwordBytes)

		remoteCmd := fmt.Sprintf("yayo-cli list-envs --p %s", project)
		if service != "" {
			remoteCmd += fmt.Sprintf(" --s %s", service)
		}

		output, err := ssh.RunSSHWithStdin(host, remoteCmd, password+"\n")
		if err != nil {
			return fmt.Errorf("ошибка при выполнении удалённой команды: %w", err)
		}

		var envItems []struct {
			Key     string `json:"key"`
			Value   string `json:"value"`
			Service string `json:"service"`
		}

		if err := json.Unmarshal([]byte(output), &envItems); err != nil {
			return fmt.Errorf("не удалось распарсить JSON: %w", err)
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "KEY\tVALUE\tSERVICE")
		fmt.Fprintln(w, "----\t-----\t-------")

		for _, e := range envItems {
			fmt.Fprintf(w, "%s\t%s\t%s\n", e.Key, e.Value, e.Service)
		}

		w.Flush()
		return nil
	},
}

func init() {
	listEnvSSHCmd.Flags().String("h", "", "SSH host, например root@1.2.3.4")
	listEnvSSHCmd.Flags().String("p", "", "Название проекта")
	listEnvSSHCmd.Flags().String("s", "", "Фильтр по сервису (необязательно)")
	RootCmd.AddCommand(listEnvSSHCmd)
}
