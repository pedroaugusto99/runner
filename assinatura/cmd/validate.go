package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"

	"github.com/pedroaugusto99/runner/assinatura/internal/assinador"
	"github.com/pedroaugusto99/runner/assinatura/internal/output"
	"github.com/spf13/cobra"
)

func newValidateCommand(cfg *cliConfig) *cobra.Command {
	var signaturePath string
	var local bool
	var port string

	cmd := &cobra.Command{
		Use:     "validate --signature <arquivo> [--local]",
		Aliases: []string{"validar", "verify"},
		Short:   "Valida uma assinatura digital simulada usando o assinador.jar",
		Long: `Valida uma assinatura digital simulada a partir de um arquivo de assinatura.

Por padrão, o CLI tentará comunicar com o servidor em background via HTTP para menor latência.
Se o servidor não estiver ativo, ou se a flag --local for utilizada, o sistema faz fallback automático
e invoca o assinador.jar como um subprocesso tradicional:
  java -jar assinador.jar validate --signature <arquivo>

Exemplos:
  assinatura validate --signature assinatura.json
  assinatura validate --signature assinatura.json --local`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !local && isServerUp(port) {
				slog.Info("Servidor detetado. Iniciando requisição REST POST /validate...")

				exitCode, httpErr := validateViaHTTP(cmd, port, signaturePath)
				if httpErr == nil {
					if exitCode != 0 {
						return exitError{Code: exitCode, Silent: true}
					}
					return nil
				}
				slog.Warn("Falha na comunicação HTTP com o servidor, prosseguindo para fallback", "erro", httpErr)
			} else if local {
				slog.Info("Modo Local (CLI) forçado explicitamente via flag --local.")
			} else {
				slog.Info("Servidor inativo. Iniciando fallback automático para Modo Local (subprocesso).")
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
	cmd.Flags().BoolVar(&local, "local", false, "Força a invocação do assinador.jar diretamente via subprocesso")
	cmd.Flags().StringVarP(&port, "port", "p", "8080", "Porta do servidor HTTP a ser consultada")

	// A obrigatoriedade do parâmetro é validada pelo assinador.jar (autoridade
	// única), não pelo CLI — ver docs/adr/0001-contrato-cli-jar.md.

	return cmd
}

func validateViaHTTP(cmd *cobra.Command, port, signaturePath string) (int, error) {
	payload := map[string]string{
		"signaturePath": signaturePath,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("falha ao preparar payload JSON: %w", err)
	}

	url := fmt.Sprintf("http://localhost:%s/api/v1/validate", port)
	resp, err := signerHTTPClient.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	res, ok := output.Render(cmd.OutOrStdout(), body)
	if !ok {
		return 0, fmt.Errorf("resposta HTTP não reconhecida (status %d)", resp.StatusCode)
	}
	return res.ExitCode(), nil
}
