package assinador

import (
	"fmt"
	"os"
	"path/filepath"
)

func ResolveJarPath(explicitPath string) (string, error) {
	if explicitPath != "" {
		return existingFile(explicitPath)
	}
	if envPath := os.Getenv("ASSINADOR_JAR"); envPath != "" {
		return existingFile(envPath)
	}

	for _, candidate := range candidateJarPaths() {
		if path, err := existingFile(candidate); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("assinador.jar não encontrado; gere com 'cd assinador && mvn package' ou informe --jar/ASSINADOR_JAR")
}

func candidateJarPaths() []string {
	const jarRelPath = "assinador/target/assinador.jar"

	var candidates []string
	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, walkUpCandidates(cwd, jarRelPath)...)
	}
	if exe, err := os.Executable(); err == nil {
		candidates = append(candidates, walkUpCandidates(filepath.Dir(exe), jarRelPath)...)
	}

	return candidates
}

func walkUpCandidates(start, relPath string) []string {
	var candidates []string
	current := filepath.Clean(start)
	for {
		candidates = append(candidates, filepath.Join(current, relPath))
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return candidates
}

func existingFile(path string) (string, error) {
	cleanPath := filepath.Clean(path)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return "", fmt.Errorf("assinador.jar não encontrado em %q", path)
	}
	if info.IsDir() {
		return "", fmt.Errorf("caminho do assinador.jar aponta para um diretório: %q", path)
	}
	return cleanPath, nil
}
