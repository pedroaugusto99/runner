package assinador

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestResolveJarPathUsesExplicitPath(t *testing.T) {
	jarPath := createTempJar(t)

	got, err := ResolveJarPath(jarPath)
	if err != nil {
		t.Fatalf("ResolveJarPath retornou erro: %v", err)
	}
	if got != filepath.Clean(jarPath) {
		t.Fatalf("esperava %q, obteve %q", filepath.Clean(jarPath), got)
	}
}

func TestResolveJarPathReturnsUsageMessageWhenMissing(t *testing.T) {
	_, err := ResolveJarPath(filepath.Join(t.TempDir(), "assinador.jar"))
	if err == nil {
		t.Fatal("esperava erro para JAR ausente")
	}
	if !strings.Contains(err.Error(), "assinador.jar não encontrado") {
		t.Fatalf("mensagem inesperada: %v", err)
	}
}

func TestRunLocalSignPassesContractArguments(t *testing.T) {
	jarPath := createTempJar(t)
	javaPath := createFakeJava(t, 0)

	result, err := RunLocal(context.Background(), LocalRequest{
		Operation:  OperationSign,
		InputPath:  "entrada com espaço.json",
		OutputPath: "saida.json",
		JavaBin:    javaPath,
		JarPath:    jarPath,
		Timeout:    time.Second,
	})
	if err != nil {
		t.Fatalf("RunLocal retornou erro: %v", err)
	}

	expected := "-jar\n" + filepath.Clean(jarPath) + "\nsign\n--input\nentrada com espaço.json\n--output\nsaida.json\n"
	if result.Stdout != expected {
		t.Fatalf("argumentos inesperados:\n--- got ---\n%s--- want ---\n%s", result.Stdout, expected)
	}
}

func TestRunLocalValidatePassesContractArguments(t *testing.T) {
	jarPath := createTempJar(t)
	javaPath := createFakeJava(t, 0)

	result, err := RunLocal(context.Background(), LocalRequest{
		Operation:     OperationValidate,
		SignaturePath: "assinatura.json",
		JavaBin:       javaPath,
		JarPath:       jarPath,
		Timeout:       time.Second,
	})
	if err != nil {
		t.Fatalf("RunLocal retornou erro: %v", err)
	}

	expected := "-jar\n" + filepath.Clean(jarPath) + "\nvalidate\n--signature\nassinatura.json\n"
	if result.Stdout != expected {
		t.Fatalf("argumentos inesperados:\n--- got ---\n%s--- want ---\n%s", result.Stdout, expected)
	}
}

func TestRunLocalPropagatesProcessOutputOnFailure(t *testing.T) {
	jarPath := createTempJar(t)
	javaPath := createFakeJava(t, 7)

	result, err := RunLocal(context.Background(), LocalRequest{
		Operation:     OperationValidate,
		SignaturePath: "invalida.txt",
		JavaBin:       javaPath,
		JarPath:       jarPath,
		Timeout:       time.Second,
	})
	if err == nil {
		t.Fatal("esperava erro quando o processo retorna código diferente de zero")
	}
	if !strings.Contains(err.Error(), "código 7") {
		t.Fatalf("erro inesperado: %v", err)
	}
	if result.Code != 7 {
		t.Fatalf("código inesperado: %d", result.Code)
	}
	if result.Stderr == "" {
		t.Fatal("esperava stderr do processo preservado")
	}
}

func TestRunLocalReportsUsageErrorWhenJavaMissing(t *testing.T) {
	jarPath := createTempJar(t)
	missingJava := filepath.Join(t.TempDir(), "java-inexistente")

	result, err := RunLocal(context.Background(), LocalRequest{
		Operation:     OperationValidate,
		SignaturePath: "assinatura.json",
		JavaBin:       missingJava,
		JarPath:       jarPath,
		Timeout:       time.Second,
	})
	if err == nil {
		t.Fatal("esperava erro quando o executável Java não existe")
	}

	var usageErr UsageError
	if !errors.As(err, &usageErr) {
		t.Fatalf("esperava UsageError, obteve %T: %v", err, err)
	}
	if usageErr.Code != 1 {
		t.Fatalf("esperava código de uso 1, obteve %d", usageErr.Code)
	}
	if !strings.Contains(usageErr.Message, "não encontrado") {
		t.Fatalf("mensagem inesperada: %v", usageErr.Message)
	}
	if result.Code != 0 {
		t.Fatalf("esperava código 0 quando o processo nem inicia, obteve %d", result.Code)
	}
}

func createTempJar(t *testing.T) string {
	t.Helper()

	jarPath := filepath.Join(t.TempDir(), "assinador.jar")
	if err := os.WriteFile(jarPath, []byte("fake jar"), 0o600); err != nil {
		t.Fatalf("falha ao criar jar temporário: %v", err)
	}
	return jarPath
}

func createFakeJava(t *testing.T, exitCode int) string {
	t.Helper()

	dir := t.TempDir()
	name := "java-fake"
	if runtime.GOOS == "windows" {
		name = "java-fake.exe"
	}
	path := filepath.Join(dir, name)

	sourcePath := filepath.Join(dir, "main.go")
	source := fmt.Sprintf(`package main

import (
	"fmt"
	"os"
)

func main() {
	for _, arg := range os.Args[1:] {
		fmt.Println(arg)
	}
	fmt.Fprintln(os.Stderr, "erro simulado")
	os.Exit(%d)
}
`, exitCode)

	if err := os.WriteFile(sourcePath, []byte(source), 0o600); err != nil {
		t.Fatalf("falha ao criar fonte do java falso: %v", err)
	}
	if out, err := exec.Command("go", "build", "-o", path, sourcePath).CombinedOutput(); err != nil {
		t.Fatalf("falha ao compilar java falso: %v\n%s", err, string(out))
	}
	return path
}
