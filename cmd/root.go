package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yayo-cli",
	Short: "Yayo-cli — локальный менеджер ENV переменных с безопасным хранилищем",
	Long:  `Yayo-cli — CLI утилита для управления ENV переменными проектов. Хранение: локальная SQLCipher база.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
