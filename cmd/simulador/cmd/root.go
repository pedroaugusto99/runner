package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd representa o comando base do simulador
var rootCmd = &cobra.Command{
	Use:   "simulador",
	Short: "Simulador do HubSaúde",
	Long:  `Interface de linha de comando para gerenciar o ciclo de vida do Simulador do HubSaúde.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("simulador v0.1.0-dev — em construção")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}