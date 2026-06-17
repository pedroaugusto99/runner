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

func newSignCommand(cfg *cliConfig) *cobra.Command {
	var inputPath string
	var outputPath string
	var local bool
	var port string

	cmd := &cobra.Command{
		Use:     "sign --input <arquivo> --output <arquivo> [--local]",
		Aliases: []string{"assinar", "criar"},
		Short:   "Cria uma assinatura digital simulada usando o assinador.jar",
		Long: `Cria uma assinatura digital simulada a partir de um arquivo de entrada.

Por padrão, o CLI tentará comunicar com o servidor em background via HTTP para menor latência.
Se o servidor não estiver ativo, ou se a flag --local for utilizada, o sistema faz fallback automático
e invoca o assinador.jar como um subprocesso tradicional:
  java -jar assinador.jar sign --input <arquivo> --output <arquivo>

Exemplos:
  assinatura sign --input bundle.json --output assinatura.json
  assinatura sign --input bundle.json --output assinatura.json --local`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !local && isServerUp(port) {
				slog.Info("Servidor detetado. Iniciando requisição REST POST /sign...")

				exitCode, httpErr := signViaHTTP(cmd, port, inputPath, outputPath)
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
	cmd.Flags().BoolVar(&local, "local", false, "Força a invocação do assinador.jar diretamente via subprocesso")
	cmd.Flags().StringVarP(&port, "port", "p", "8080", "Porta do servidor HTTP a ser consultada")


	return cmd
}

func signViaHTTP(cmd *cobra.Command, port, inputPath, outputPath string) (int, error) {
	payload := map[string]string{
		"inputPath":  inputPath,
		"outputPath": outputPath,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("falha ao preparar payload JSON: %w", err)
	}

	url := fmt.Sprintf("http://localhost:%s/api/v1/sign", port)
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
