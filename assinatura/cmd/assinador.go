package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pedroaugusto99/runner/assinatura/internal/assinador"
	"github.com/pedroaugusto99/runner/assinatura/internal/output"
	"github.com/spf13/cobra"
)

var jvmQuietEnv = map[string]string{
	"SPRING_MAIN_WEB_APPLICATION_TYPE": "none",
	"SPRING_MAIN_BANNER_MODE":          "off",
	"LOGGING_LEVEL_ROOT":               "off",
}

func runAssinador(cmd *cobra.Command, req assinador.LocalRequest) error {
	if os.Getenv("SOFTHSM2_CONF") == "" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			configPath := filepath.Join(homeDir, ".config", "softhsm2", "softhsm2.conf")
			os.Setenv("SOFTHSM2_CONF", configPath)
		}
	}

	for key, value := range jvmQuietEnv {
		old, had := os.LookupEnv(key)
		os.Setenv(key, value)
		defer func(k, v string, restore bool) {
			if restore {
				os.Setenv(k, v)
			} else {
				os.Unsetenv(k)
			}
		}(key, old, had)
	}

	result, err := assinador.RunLocal(cmd.Context(), req)

	rendered := false
	if result.Stdout != "" {
		_, rendered = output.Render(cmd.OutOrStdout(), []byte(result.Stdout))
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
		if rendered {
			return exitError{Code: code, Silent: true}
		}
		fmt.Fprintf(cmd.ErrOrStderr(), "✖ %s\n", usageErr.Message)
		return exitError{Code: code, Silent: true}
	}

	fmt.Fprintf(cmd.ErrOrStderr(), "✖ falha inesperada ao executar o assinador.jar: %v\n", err)
	return exitError{Code: exitUnexpected, Silent: true}
}
