package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	exitUsage      = 1
	exitUnexpected = 2
)

type cliConfig struct {
	version string
	commit  string
	javaBin string
	jarPath string
	timeout time.Duration
}

func newRootCommand(version, commit string, stdout, stderr io.Writer) *cobra.Command {
	cfg := &cliConfig{version: version, commit: commit}

	rootCmd := &cobra.Command{
		Use:   "assinatura",
		Short: "CLI do Sistema Runner para operações de assinatura digital",
		Long: `assinatura facilita a execução do assinador.jar sem exigir que o usuário
conheça os detalhes de invocação da JVM.

Exemplos:
  assinatura sign --input bundle.json --output assinatura.json --local
  assinatura validate --signature assinatura.json --local
  assinatura version`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       formatVersion(cfg),
	}

	rootCmd.SetOut(stdout)
	rootCmd.SetErr(stderr)
	rootCmd.PersistentFlags().StringVar(&cfg.javaBin, "java", "java", "Executável Java usado para invocar o assinador.jar")
	rootCmd.PersistentFlags().StringVar(&cfg.jarPath, "jar", "", "Caminho para o assinador.jar; quando omitido, ASSINADOR_JAR ou caminhos conhecidos são usados")
	rootCmd.PersistentFlags().DurationVar(&cfg.timeout, "timeout", 30*time.Second, "Tempo máximo de espera por uma chamada ao assinador")

	rootCmd.AddCommand(newVersionCommand(cfg))
	rootCmd.AddCommand(newSignCommand(cfg))
	rootCmd.AddCommand(newValidateCommand(cfg))
	rootCmd.AddCommand(newStartCommand(cfg))
	rootCmd.AddCommand(newStopCommand(cfg))

	return rootCmd
}

func Execute(version, commit string) {
	cmd := newRootCommand(version, commit, os.Stdout, os.Stderr)
	if err := cmd.Execute(); err != nil {
		var exitErr exitError
		if errors.As(err, &exitErr) {
			if !exitErr.Silent && exitErr.Message != "" {
				fmt.Fprintln(os.Stderr, exitErr.Message)
			}
			os.Exit(exitErr.Code)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(exitUsage)
	}
}

type exitError struct {
	Message string
	Code    int
	Silent bool
}

func (e exitError) Error() string {
	return e.Message
}
