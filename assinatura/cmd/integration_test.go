package cmd

import (
	"bytes"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pedroaugusto99/runner/assinatura/internal/assinador"
)

func resolveEnvOrSkip(t *testing.T) string {
	t.Helper()
	if testing.Short() {
		t.Skip("teste de integração ignorado em modo -short")
	}
	jarPath, err := assinador.ResolveJarPath("")
	if err != nil {
		t.Skipf("assinador.jar indisponível: %v", err)
	}
	if _, err := exec.LookPath("java"); err != nil {
		t.Skipf("java indisponível no PATH: %v", err)
	}
	return jarPath
}

func writeBundle(t *testing.T) (dir, input string) {
	t.Helper()
	dir = t.TempDir()
	input = filepath.Join(dir, "bundle.json")
	if err := os.WriteFile(input, []byte(`{"resourceType":"Bundle"}`), 0o600); err != nil {
		t.Fatalf("falha ao criar FHIR de entrada: %v", err)
	}
	return dir, input
}

func TestIntegrationLocalSignThenValidate(t *testing.T) {
	jarPath := resolveEnvOrSkip(t)
	dir, inputPath := writeBundle(t)
	outputPath := filepath.Join(dir, "assinatura.json")

	var stdout, stderr bytes.Buffer
	root := newRootCommand("dev", "test", &stdout, &stderr)
	root.SetArgs([]string{"sign", "--local", "--jar", jarPath, "--input", inputPath, "--output", outputPath})
	if err := root.Execute(); err != nil {
		t.Fatalf("sign retornou erro: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Assinatura criada") {
		t.Fatalf("saída de sign não foi renderizada como esperado: %q", stdout.String())
	}
	if _, err := os.Stat(outputPath); err != nil {
		t.Fatalf("assinatura não foi gravada em %s: %v", outputPath, err)
	}

	stdout.Reset()
	stderr.Reset()
	root = newRootCommand("dev", "test", &stdout, &stderr)
	root.SetArgs([]string{"validate", "--local", "--jar", jarPath, "--signature", outputPath})
	if err := root.Execute(); err != nil {
		t.Fatalf("validate retornou erro: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "VÁLIDA") {
		t.Fatalf("saída de validate não indicou validade: %q", stdout.String())
	}
}

func TestIntegrationLocalRejectsMissingParameter(t *testing.T) {
	jarPath := resolveEnvOrSkip(t)
	dir, _ := writeBundle(t)

	var stdout, stderr bytes.Buffer
	root := newRootCommand("dev", "test", &stdout, &stderr)
	root.SetArgs([]string{"sign", "--local", "--jar", jarPath, "--output", filepath.Join(dir, "x.json")})
	err := root.Execute()

	if err == nil {
		t.Fatal("esperava código de saída diferente de zero para parâmetro ausente")
	}
	out := stdout.String()
	if !strings.Contains(out, "input") || !strings.Contains(out, "obrigatório") {
		t.Fatalf("erro renderizado deveria citar o parâmetro e o motivo: %q", out)
	}
}

func TestIntegrationHTTPSignAndValidate(t *testing.T) {
	jarPath := resolveEnvOrSkip(t)
	port := freePort(t)
	startServer(t, jarPath, port)

	dir, inputPath := writeBundle(t)
	outputPath := filepath.Join(dir, "assinatura.json")

	var stdout, stderr bytes.Buffer
	root := newRootCommand("dev", "test", &stdout, &stderr)
	root.SetArgs([]string{"sign", "--port", port, "--input", inputPath, "--output", outputPath})
	if err := root.Execute(); err != nil {
		t.Fatalf("sign HTTP retornou erro: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Assinatura criada") {
		t.Fatalf("sign HTTP não renderizou sucesso: %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	root = newRootCommand("dev", "test", &stdout, &stderr)
	root.SetArgs([]string{"validate", "--port", port, "--signature", outputPath})
	if err := root.Execute(); err != nil {
		t.Fatalf("validate HTTP retornou erro: %v\nstderr: %s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "VÁLIDA") {
		t.Fatalf("validate HTTP não indicou validade: %q", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	root = newRootCommand("dev", "test", &stdout, &stderr)
	root.SetArgs([]string{"sign", "--port", port, "--input", filepath.Join(dir, "ausente.json"), "--output", outputPath})
	if err := root.Execute(); err == nil {
		t.Fatal("esperava erro para input inexistente via HTTP")
	}
	if !strings.Contains(stdout.String(), "não encontrado") {
		t.Fatalf("erro HTTP deveria ser renderizado com motivo: %q", stdout.String())
	}
}

func freePort(t *testing.T) string {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("falha ao alocar porta livre: %v", err)
	}
	defer l.Close()
	return strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
}

func startServer(t *testing.T, jarPath, port string) {
	t.Helper()

	var logBuf bytes.Buffer
	server := exec.Command("java", "-jar", jarPath)
	server.Env = append(os.Environ(), "SERVER_PORT="+port)
	server.Stdout = &logBuf
	server.Stderr = &logBuf
	if err := server.Start(); err != nil {
		t.Fatalf("falha ao iniciar o servidor: %v", err)
	}
	t.Cleanup(func() {
		_ = server.Process.Kill()
		_ = server.Wait()
	})

	deadline := time.Now().Add(40 * time.Second)
	for time.Now().Before(deadline) {
		if isServerUp(port) {
			return
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("servidor não respondeu ao health check na porta %s.\nlog:\n%s", port, logBuf.String())
}
