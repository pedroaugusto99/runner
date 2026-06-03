package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCommand(cfg *cliConfig) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Exibe a versão atual da aplicação",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "assinatura %s\n", formatVersion(cfg))
		},
	}
}

func formatVersion(cfg *cliConfig) string {
	if cfg.commit == "" {
		return cfg.version
	}
	return fmt.Sprintf("%s (commit %s)", cfg.version, cfg.commit)
}
