package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/1Storm3/yayo-cli/internal/config"
	"github.com/spf13/cobra"
)

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Удаляет локальное зашифрованное хранилище проекта и все связанные файлы",
	RunE: func(cmd *cobra.Command, args []string) error {
		project, _ := cmd.Flags().GetString("p")
		if project == "" {
			return fmt.Errorf("необходимо указать --p")
		}

		projectDir := filepath.Join(config.BaseDir(), project)

		if _, err := os.Stat(projectDir); os.IsNotExist(err) {
			return fmt.Errorf("проект '%s' не найден", project)
		}

		if err := os.RemoveAll(projectDir); err != nil {
			return fmt.Errorf("не удалось удалить проект: %w", err)
		}

		fmt.Printf("Проект '%s' и все данные успешно удалены!\n", project)
		return nil
	},
}

func init() {
	destroyCmd.Flags().String("p", "", "Название проекта для удаления")
	RootCmd.AddCommand(destroyCmd)
}
