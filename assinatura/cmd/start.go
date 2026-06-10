package cmd

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"

	"github.com/pedroaugusto99/runner/assinatura/internal/java"
)

func newStartCommand(cfg *cliConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Inicia o assinador.jar em background (Modo Servidor)",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetString("port")

			slog.Info("Iniciando verificação de status do servidor", "porta", port)
			if isServerUp(port) {
				cmd.Printf("✅ O servidor já está em execução e saudável na porta %s!\n", port)
				return nil
			}

			slog.Debug("Servidor inativo. Preparando para acionar o provisionador Java.")

			javaBin, err := java.EnsureJRE()
			if err != nil {
				slog.Warn("Falha no provisionador, usando fallback", "erro", err)
				cmd.PrintErrf("⚠️ Aviso do provisionador: %v\n", err)
				javaBin = "java"
			}

			jarPath := filepath.Join("..", "assinador", "target", "assinador.jar")
			if _, err := os.Stat(jarPath); os.IsNotExist(err) {
				return fmt.Errorf("assinador.jar não encontrado em %s. Execute 'mvn package' primeiro", jarPath)
			}

			serverCmd := exec.Command(javaBin, "-jar", jarPath)

			slog.Debug("Iniciando processo da JVM", "binario", javaBin, "jar", jarPath)
			cmd.Println("A iniciar a Máquina Virtual Java...")

			if err := serverCmd.Start(); err != nil {
				slog.Error("Falha fatal ao iniciar JVM", "erro", err)
				return fmt.Errorf("falha ao iniciar o servidor: %w", err)
			}

			pid := serverCmd.Process.Pid
			if err := savePID(pid); err != nil {
				slog.Error("Falha ao guardar ficheiro PID", "pid", pid, "erro", err)
				cmd.PrintErrf("⚠️ Servidor iniciado (PID: %d), mas falhou ao guardar o ficheiro PID: %v\n", pid, err)
			} else {
				slog.Info("Máquina Virtual Java acionada em background", "pid", pid)
				cmd.Printf("🚀 Servidor iniciado com sucesso em background! (PID: %d)\n", pid)
			}

			cmd.Println("A aguardar a inicialização do Spring Boot...")
			time.Sleep(3 * time.Second)

			if isServerUp(port) {
				slog.Info("Servidor HTTP confirmou inicialização", "porta", port)
				cmd.Println("✅ Servidor pronto para receber requisições HTTP!")
			} else {
				slog.Warn("Servidor não respondeu ao health check após timeout inicial", "porta", port)
				cmd.Println("⏳ O servidor foi iniciado, mas ainda não está a responder. Verifique os logs.")
			}

			return nil
		},
	}

	cmd.Flags().StringP("port", "p", "8080", "Porta do servidor HTTP")

	return cmd
}

func isServerUp(port string) bool {
	client := http.Client{
		Timeout: 2 * time.Second,
	}
	url := fmt.Sprintf("http://localhost:%s/api/v1/health", port)

	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		if status, ok := result["status"].(string); ok && status == "UP" {
			return true
		}
	}

	return false
}

func savePID(pid int) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	hubDir := filepath.Join(home, ".hubsaude")
	if err := os.MkdirAll(hubDir, 0755); err != nil {
		return err
	}

	pidFile := filepath.Join(hubDir, "assinador.pid")
	return os.WriteFile(pidFile, []byte(strconv.Itoa(pid)), 0644)
}
