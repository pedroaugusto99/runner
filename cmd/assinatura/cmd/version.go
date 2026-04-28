package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Exibe a versão atual da aplicação",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("assinatura v0.1.0-dev")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}