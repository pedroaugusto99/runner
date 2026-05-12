package main

import (
	"os/exec"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("falha ao executar comando: %v", err)
	}

	expected := "dev"
	if !strings.Contains(string(out), expected) {
		t.Errorf("esperava conter %q, obteve %q", expected, string(out))
	}
}
