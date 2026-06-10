package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/pedroaugusto99/runner/assinatura/internal/assinador"
	"github.com/spf13/cobra"
)

func runAssinador(cmd *cobra.Command, req assinador.LocalRequest) error {
	os.Setenv("SPRING_MAIN_WEB_APPLICATION_TYPE", "none")

	defer os.Unsetenv("SPRING_MAIN_WEB_APPLICATION_TYPE")

	result, err := assinador.RunLocal(cmd.Context(), req)
	if result.Stdout != "" {
		fmt.Fprint(cmd.OutOrStdout(), result.Stdout)
	}
	if result.Stderr != "" {
		fmt.Fprint(cmd.ErrOrStderr(), result.Stderr)
	}
	if err == nil {
		return nil
	}

	var usageErr assinador.UsageError
	if errors.As(err, &usageErr) {
		code := usageErr.Code
		if code == 0 {
			code = exitUsage
		}
		return exitError{Message: usageErr.Message, Code: code}
	}

	return exitError{Message: fmt.Sprintf("falha inesperada ao executar o assinador.jar: %v", err), Code: exitUnexpected}
}
