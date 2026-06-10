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

			if !local {
				slog.Debug("Verificando se o servidor está ativo para invocação HTTP", "porta", port)

				if isServerUp(port) {
					slog.Info("Servidor detetado. Iniciando requisição REST POST /validate...")

					httpErr := validateViaHTTP(cmd, port, signaturePath)

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

			cmd.Println("⚙️ A processar validação no Modo Local (CLI JVM)...")

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

	_ = cmd.MarkFlagRequired("signature")

	return cmd
}

func validateViaHTTP(cmd *cobra.Command, port, signaturePath string) error {
	payload := map[string]string{
		"signaturePath": signaturePath,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("falha ao preparar payload JSON: %w", err)
	}

	url := fmt.Sprintf("http://localhost:%s/api/v1/validate", port)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusOK {
		slog.Info("Validação concluída via HTTP")
		cmd.Println("✅ Validação concluída com sucesso via Modo Servidor (HTTP)!")
		cmd.Printf("Resposta do Servidor: %s\n", string(body))
		return nil
	}

	slog.Warn("Servidor retornou erro HTTP", "status", resp.StatusCode)
	cmd.PrintErrf("⚠️ O Servidor HTTP processou o pedido mas retornou um erro (%d): %s\n", resp.StatusCode, string(body))
	return fmt.Errorf("falha na validação de negócio via HTTP: status %d", resp.StatusCode)
}
