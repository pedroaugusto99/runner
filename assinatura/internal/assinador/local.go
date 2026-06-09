package assinador

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os/exec"
	"time"

	"github.com/pedroaugusto99/runner/assinatura/internal/java"
)

func RunLocal(ctx context.Context, req LocalRequest) (LocalResult, error) {
	if req.JavaBin == "" {
		javaPath, provErr := java.EnsureJRE()
		if provErr != nil {
			fmt.Printf("Aviso do Provisionador: %v\n", provErr)
			req.JavaBin = "java"
		} else {
			req.JavaBin = javaPath
		}
	}

	if req.Timeout <= 0 {
		req.Timeout = 30 * time.Second
	}

	jarPath, err := ResolveJarPath(req.JarPath)
	if err != nil {
		return LocalResult{}, UsageError{Message: err.Error(), Code: 1}
	}

	args, err := req.args(jarPath)
	if err != nil {
		return LocalResult{}, UsageError{Message: err.Error(), Code: 1}
	}

	runCtx, cancel := context.WithTimeout(ctx, req.Timeout)
	defer cancel()

	command := exec.CommandContext(runCtx, req.JavaBin, args...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &stdout
	command.Stderr = &stderr

	err = command.Run()
	result := LocalResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}
	if command.ProcessState != nil {
		result.Code = command.ProcessState.ExitCode()
	}

	if errors.Is(runCtx.Err(), context.DeadlineExceeded) {
		return result, UsageError{Message: fmt.Sprintf("tempo limite excedido ao executar o assinador.jar após %s", req.Timeout), Code: 1}
	}
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return result, UsageError{Message: fmt.Sprintf("assinador.jar finalizou com código %d", exitErr.ExitCode()), Code: exitErr.ExitCode()}
		}
		if errors.Is(err, exec.ErrNotFound) || errors.Is(err, fs.ErrNotExist) {
			return result, UsageError{Message: fmt.Sprintf("executável Java %q não encontrado; instale o JDK ou informe o caminho com --java", req.JavaBin), Code: 1}
		}
		return result, err
	}

	return result, nil
}

func (r LocalRequest) args(jarPath string) ([]string, error) {
	base := []string{"-jar", jarPath}

	switch r.Operation {
	case OperationSign:
		if r.InputPath == "" {
			return nil, fmt.Errorf("informe o arquivo de entrada com --input")
		}
		if r.OutputPath == "" {
			return nil, fmt.Errorf("informe o arquivo de saída com --output")
		}
		return append(base, "sign", "--input", r.InputPath, "--output", r.OutputPath), nil
	case OperationValidate:
		if r.SignaturePath == "" {
			return nil, fmt.Errorf("informe o arquivo de assinatura com --signature")
		}
		return append(base, "validate", "--signature", r.SignaturePath), nil
	default:
		return nil, fmt.Errorf("operação não suportada: %s", r.Operation)
	}
}
