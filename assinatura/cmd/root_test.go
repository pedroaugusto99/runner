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

func TestSignWithoutLocalReportsHTTPNotAvailable(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	root := newRootCommand("dev", "unknown", &stdout, &stderr)
	root.SetArgs([]string{"sign", "--input", "entrada.json", "--output", "assinatura.json"})

	err := root.Execute()
	if err == nil {
		t.Fatal("esperava erro quando modo HTTP é solicitado")
	}
	if !strings.Contains(err.Error(), "modo HTTP ainda não está disponível") {
		t.Fatalf("erro inesperado: %v", err)
	}
}
