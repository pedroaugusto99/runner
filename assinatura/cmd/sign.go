package cmd

import (
	"github.com/pedroaugusto99/runner/assinatura/internal/assinador"
	"github.com/spf13/cobra"
)

func newSignCommand(cfg *cliConfig) *cobra.Command {
	var inputPath string
	var outputPath string
	var local bool

	cmd := &cobra.Command{
		Use:     "sign --input <arquivo> --output <arquivo> --local",
		Aliases: []string{"assinar", "criar"},
		Short:   "Cria uma assinatura digital simulada usando o assinador.jar",
		Long: `Cria uma assinatura digital simulada a partir de um arquivo de entrada.

O modo local invoca o assinador.jar como subprocesso:
  java -jar assinador.jar sign --input <arquivo> --output <arquivo>

Exemplo:
  assinatura sign --input bundle.json --output assinatura.json --local`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !local {
				return exitError{Message: "modo HTTP ainda não está disponível neste assinador.jar; use --local para invocação direta", Code: exitUsage}
			}

			return runAssinador(cmd, assinador.LocalRequest{
				Operation:  assinador.OperationSign,
				InputPath:  inputPath,
				OutputPath: outputPath,
				JavaBin:    cfg.javaBin,
				JarPath:    cfg.jarPath,
				Timeout:    cfg.timeout,
			})
		},
	}

	cmd.Flags().StringVarP(&inputPath, "input", "i", "", "Caminho do arquivo FHIR/JSON de entrada")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Caminho onde a assinatura simulada será gravada")
	cmd.Flags().BoolVar(&local, "local", false, "Invoca o assinador.jar diretamente via CLI")
	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}
