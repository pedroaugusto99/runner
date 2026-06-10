package cmd

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

func newStopCommand(cfg *cliConfig) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stop",
		Short: "Encerra a execução do assinador.jar em background",
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.Info("Iniciando processo de encerramento do servidor")

			pidFile, err := getPIDFilePath()
			if err != nil {
				return err
			}

			if _, err := os.Stat(pidFile); os.IsNotExist(err) {
				slog.Info("Ficheiro PID não encontrado", "caminho", pidFile)
				cmd.Println("ℹ️ Nenhum servidor parece estar em execução (ficheiro PID não encontrado).")
				return nil
			}

			pid, err := readPIDFromFile(pidFile)
			if err != nil {
				return err
			}

			cmd.Printf("A encerrar o servidor (PID: %d)...\n", pid)

			if err := terminateProcess(pid); err != nil {
				if errors.Is(err, os.ErrProcessDone) {
					cmd.Printf("⚠️ Processo %d não encontrado. A limpar ficheiro PID fantasma...\n", pid)
				} else {
					return err
				}
			} else {
				slog.Info("Servidor encerrado com sucesso", "pid", pid)
				cmd.Println("✅ Servidor encerrado com sucesso!")
			}

			removePIDFile(pidFile)

			return nil
		},
	}

	cmd.Flags().StringP("port", "p", "8080", "Porta do processo a encerrar")

	return cmd
}

func getPIDFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Falha ao obter diretório home", "erro", err)
		return "", fmt.Errorf("não foi possível aceder ao diretório de utilizador: %w", err)
	}
	return filepath.Join(home, ".hubsaude", "assinador.pid"), nil
}

func readPIDFromFile(pidFile string) (int, error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		slog.Error("Falha ao ler ficheiro PID", "erro", err)
		return 0, fmt.Errorf("falha ao ler ficheiro PID: %w", err)
	}

	pidStr := strings.TrimSpace(string(data))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		slog.Error("PID inválido no ficheiro", "conteudo", pidStr, "erro", err)
		return 0, fmt.Errorf("ficheiro PID corrompido: %w", err)
	}
	return pid, nil
}

func terminateProcess(pid int) error {
	slog.Debug("Tentando localizar processo do SO", "pid", pid)
	process, err := os.FindProcess(pid)
	if err != nil {
		slog.Warn("Processo não encontrado pelo SO", "pid", pid, "erro", err)
		return os.ErrProcessDone
	}

	err = process.Signal(os.Interrupt)
	if err != nil {
		slog.Warn("Falha no sinal de interrupção, tentando terminar forçadamente", "erro", err)
		killErr := process.Kill()
		if killErr != nil && !errors.Is(killErr, os.ErrProcessDone) {
			slog.Error("Falha ao encerrar processo forçadamente", "pid", pid, "erro", killErr)
			return fmt.Errorf("não foi possível encerrar o processo %d: %w", pid, killErr)
		}
	} else {
		slog.Debug("Sinal de interrupção enviado com sucesso")
	}

	return nil
}

func removePIDFile(pidFile string) {
	if err := os.Remove(pidFile); err != nil {
		slog.Warn("Falha ao remover ficheiro PID", "caminho", pidFile, "erro", err)
	} else {
		slog.Info("Ficheiro PID removido com sucesso")
	}
}
