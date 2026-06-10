package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommandUsesInjectedVersion(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := newRootCommand("v1.2.3-test", "abc1234", &stdout, &stderr)
	root.SetArgs([]string{"version"})

	if err := root.Execute(); err != nil {
		t.Fatalf("Execute retornou erro: %v", err)
	}
	if got := strings.TrimSpace(stdout.String()); got != "assinatura v1.2.3-test (commit abc1234)" {
		t.Fatalf("saída inesperada: %q", got)
	}
	if stderr.Len() != 0 {
		t.Fatalf("stderr inesperado: %q", stderr.String())
	}
}

func TestSignWithoutLocalReportsFallback(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := newRootCommand("dev", "unknown", &stdout, &stderr)
	root.SetArgs([]string{"sign", "--input", "entrada.json", "--output", "assinatura.json", "--port", "9999"})

	_ = root.Execute()

	got := stdout.String()
	if !strings.Contains(got, "Servidor inativo") && !strings.Contains(got, "A processar operação no Modo Local") {
		t.Errorf("Esperava que o sistema tentasse fallback para o modo local. Saída: %q", got)
	}
}
