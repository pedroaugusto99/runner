package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/pedroaugusto99/runner/assinatura/internal/assinador"
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

			if !local {
				slog.Debug("Verificando se o servidor está ativo para invocação HTTP", "porta", port)

				if isServerUp(port) {
					slog.Info("Servidor detetado. Iniciando requisição REST POST /sign...")

					httpErr := signViaHTTP(cmd, port, inputPath, outputPath)

					if httpErr == nil {
						return nil
					}

					slog.Warn("Falha na comunicação HTTP com o servidor, prosseguindo para fallback", "erro", httpErr)
				} else {
					slog.Info("Servidor inativo. Iniciando fallback automático para Modo Local (subprocesso).")
				}
			} else {
				slog.Info("Modo Local (CLI) forçado explicitamente via flag --local.")
			}

			cmd.Println("⚙️ A processar operação no Modo Local (CLI JVM)...")

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

	_ = cmd.MarkFlagRequired("input")
	_ = cmd.MarkFlagRequired("output")

	return cmd
}

func signViaHTTP(cmd *cobra.Command, port, inputPath, outputPath string) error {
	payload := map[string]string{
		"inputPath":  inputPath,
		"outputPath": outputPath,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("falha ao preparar payload JSON: %w", err)
	}

	url := fmt.Sprintf("http://localhost:%s/api/v1/sign", port)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		slog.Info("Assinatura concluída via HTTP")
		cmd.Println("✅ Assinatura concluída com sucesso via Modo Servidor (HTTP)!")
		cmd.Printf("Resposta do Servidor: %s\n", string(body))
		return nil
	}

	slog.Warn("Servidor retornou erro HTTP", "status", resp.StatusCode)
	cmd.PrintErrf("⚠️ O Servidor HTTP processou o pedido mas retornou um erro (%d): %s\n", resp.StatusCode, string(body))
	return fmt.Errorf("falha na validação de negócio via HTTP: status %d", resp.StatusCode)
}
