package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd representa o comando base do assinador CLI
var rootCmd = &cobra.Command{
	Use:   "assinatura",
	Short: "CLI para o Sistema Runner de Assinatura Digital",
	Long: `O Sistema Runner facilita a execução de operações de assinatura 
e validação digital através do assinador.jar, ocultando a complexidade 
de configuração do ambiente Java.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}