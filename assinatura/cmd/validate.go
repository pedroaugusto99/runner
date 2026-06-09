package cmd

import (
	"github.com/pedroaugusto99/runner/assinatura/internal/assinador"
	"github.com/spf13/cobra"
)

func newValidateCommand(cfg *cliConfig) *cobra.Command {
	var signaturePath string
	var local bool

	cmd := &cobra.Command{
		Use:     "validate --signature <arquivo> --local",
		Aliases: []string{"validar", "verify"},
		Short:   "Valida uma assinatura digital simulada usando o assinador.jar",
		Long: `Valida uma assinatura digital simulada a partir de um arquivo de assinatura.

O modo local invoca o assinador.jar como subprocesso:
  java -jar assinador.jar validate --signature <arquivo>

Exemplo:
  assinatura validate --signature assinatura.json --local`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !local {
				return exitError{Message: "modo HTTP ainda não está disponível neste assinador.jar; use --local para invocação direta", Code: exitUsage}
			}

			return runAssinador(cmd, assinador.LocalRequest{
				Operation:     assinador.OperationValidate,
				SignaturePath: signaturePath,
				JavaBin:       cfg.javaBin,
				JarPath:       cfg.jarPath,
				Timeout:       cfg.timeout,
			})
		},
	}

	cmd.Flags().StringVarP(&signaturePath, "signature", "s", "", "Caminho do arquivo de assinatura a validar")
	cmd.Flags().BoolVar(&local, "local", false, "Invoca o assinador.jar diretamente via CLI")
	_ = cmd.MarkFlagRequired("signature")

	return cmd
}
