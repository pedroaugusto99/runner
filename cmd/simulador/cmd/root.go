package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "simulador",
	Short: "Simulador do HubSaúde",
	Long:  `Interface de linha de comando para gerenciar o ciclo de vida do Simulador do HubSaúde.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("simulador v%s — em construção\n", Version)
	},
}

func Execute(v string) {
	Version = v
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
