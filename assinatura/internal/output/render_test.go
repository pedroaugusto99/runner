package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderSignSuccess(t *testing.T) {
	raw := []byte(`{"success":true,"operation":"sign","signature":"...","output":"assinatura.json","metadata":{"simulated":true}}`)

	var buf bytes.Buffer
	res, ok := Render(&buf, raw)
	if !ok {
		t.Fatal("esperava renderização bem-sucedida")
	}
	out := buf.String()
	if !strings.Contains(out, "✔") || !strings.Contains(out, "assinatura.json") {
		t.Fatalf("saída de sign inesperada: %q", out)
	}
	if code := res.ExitCode(); code != 0 {
		t.Fatalf("esperava exit 0 para sign de sucesso, obteve %d", code)
	}
}

func TestRenderValidateValidAndInvalid(t *testing.T) {
	valid := []byte(`{"success":true,"operation":"validate","valid":true,"message":"ok"}`)
	var buf bytes.Buffer
	res, ok := Render(&buf, valid)
	if !ok || !strings.Contains(buf.String(), "VÁLIDA") {
		t.Fatalf("esperava VÁLIDA, obteve %q", buf.String())
	}
	if res.ExitCode() != 0 {
		t.Fatalf("esperava exit 0 para válida, obteve %d", res.ExitCode())
	}

	invalid := []byte(`{"success":true,"operation":"validate","valid":false,"message":"nao confere"}`)
	buf.Reset()
	res, ok = Render(&buf, invalid)
	if !ok || !strings.Contains(buf.String(), "INVÁLIDA") {
		t.Fatalf("esperava INVÁLIDA, obteve %q", buf.String())
	}
	if res.ExitCode() != 1 {
		t.Fatalf("esperava exit 1 para inválida, obteve %d", res.ExitCode())
	}
}

func TestRenderValidationError(t *testing.T) {
	raw := []byte(`{"success":false,"operation":"sign","error":{"code":"VALIDATION_ERROR","message":"Parâmetro inválido: input (required).","details":[{"field":"input","reason":"required"}]}}`)

	var buf bytes.Buffer
	res, ok := Render(&buf, raw)
	if !ok {
		t.Fatal("esperava renderização do envelope de erro")
	}
	out := buf.String()
	if !strings.Contains(out, "✖") || !strings.Contains(out, "input") {
		t.Fatalf("saída de erro deveria citar o campo: %q", out)
	}
	if !strings.Contains(out, "obrigatório") || !strings.Contains(out, "→") {
		t.Fatalf("saída de erro deveria conter motivo legível e orientação: %q", out)
	}
	if res.ExitCode() != 1 {
		t.Fatalf("esperava exit 1 para erro de sign, obteve %d", res.ExitCode())
	}
}

func TestRenderExtractsEnvelopeFromNoisyOutput(t *testing.T) {
	raw := []byte(`  .   ____ Spring Boot banner
2026-06-16 INFO  Starting AssinadorApplication
{"success":true,"operation":"validate","valid":true,"message":"ok"}
`)

	var buf bytes.Buffer
	_, ok := Render(&buf, raw)
	if !ok {
		t.Fatal("esperava extrair o envelope de uma saída ruidosa")
	}
	if !strings.Contains(buf.String(), "VÁLIDA") {
		t.Fatalf("saída inesperada: %q", buf.String())
	}
}

func TestRenderFallsBackToRawWhenNotEnvelope(t *testing.T) {
	raw := []byte("erro inesperado da JVM\nstack trace...")

	var buf bytes.Buffer
	res, ok := Render(&buf, raw)
	if ok || res != nil {
		t.Fatal("não deveria reconhecer envelope em texto cru")
	}
	if !strings.Contains(buf.String(), "erro inesperado") {
		t.Fatalf("fallback deveria imprimir o texto cru: %q", buf.String())
	}
}

func TestValidateExitCodeWhenParameterError(t *testing.T) {
	res := &Result{Success: false, Operation: "validate"}
	if res.ExitCode() != 2 {
		t.Fatalf("esperava exit 2 para erro de parâmetro em validate, obteve %d", res.ExitCode())
	}
}
