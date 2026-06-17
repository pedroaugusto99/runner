package cmd

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/spf13/cobra"
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

func newTestRoot(args ...string) (*cobra.Command, *bytes.Buffer, *bytes.Buffer) {
	var stdout, stderr bytes.Buffer
	root := newRootCommand("dev", "unknown", &stdout, &stderr)
	root.SetArgs(args)
	return root, &stdout, &stderr
}

func TestSignParsesFlags(t *testing.T) {
	cmd, _, err := newRootCommand("dev", "x", io.Discard, io.Discard).Find([]string{"sign"})
	if err != nil {
		t.Fatalf("comando sign não encontrado: %v", err)
	}
	for _, flag := range []string{"input", "output", "local", "port"} {
		if cmd.Flags().Lookup(flag) == nil {
			t.Fatalf("comando sign deveria expor a flag --%s", flag)
		}
	}
}

func TestHelpDocumentsCommandsAndFlags(t *testing.T) {
	root, stdout, _ := newTestRoot("--help")
	if err := root.Execute(); err != nil {
		t.Fatalf("Execute --help retornou erro: %v", err)
	}
	help := stdout.String()
	for _, want := range []string{"sign", "validate"} {
		if !strings.Contains(help, want) {
			t.Fatalf("ajuda deveria documentar o comando %q. Saída: %s", want, help)
		}
	}

	signHelp, signOut, _ := newTestRoot("sign", "--help")
	if err := signHelp.Execute(); err != nil {
		t.Fatalf("Execute sign --help retornou erro: %v", err)
	}
	for _, want := range []string{"--input", "--output"} {
		if !strings.Contains(signOut.String(), want) {
			t.Fatalf("ajuda do sign deveria documentar %q. Saída: %s", want, signOut.String())
		}
	}
}

func TestCommandsExposeAliases(t *testing.T) {
	root := newRootCommand("dev", "unknown", io.Discard, io.Discard)

	aliases := map[string][]string{
		"sign":     {"assinar", "criar"},
		"validate": {"validar", "verify"},
	}
	for name, wantAliases := range aliases {
		cmd, _, err := root.Find([]string{name})
		if err != nil {
			t.Fatalf("comando %q não encontrado: %v", name, err)
		}
		for _, alias := range wantAliases {
			if !contains(cmd.Aliases, alias) {
				t.Fatalf("comando %q deveria expor o alias %q (tem %v)", name, alias, cmd.Aliases)
			}
		}
	}
}

func contains(values []string, target string) bool {
	for _, v := range values {
		if v == target {
			return true
		}
	}
	return false
}
