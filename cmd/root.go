package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "yayo",
	Short: "Yayo — локальный менеджер ENV переменных с безопасным хранилищем",
	Long:  `Yayo — CLI утилита для управления ENV переменными проектов. Хранение: локальная SQLCipher база.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
