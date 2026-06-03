package java

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func EnsureJRE() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("erro ao obter diretório home: %v", err)
	}

	hubsaudeDir := filepath.Join(home, ".hubsaude", "jdk")

	if _, err := os.Stat(hubsaudeDir); !os.IsNotExist(err) {
		if exePath, err := findJavaExe(hubsaudeDir); err == nil {
			fmt.Println("[Provisioner] JRE local encontrado em:", exePath)
			return exePath, nil
		}
	}

	fmt.Println("[Provisioner] JRE 21 não encontrado localmente.")
	fmt.Println("[Provisioner] Iniciando download do servidor Adoptium. Isto pode demorar um pouco...")

	if err := os.MkdirAll(hubsaudeDir, 0755); err != nil {
		return "", err
	}

	downloadURL, ext := getAdoptiumURL()
	if downloadURL == "" {
		return "", fmt.Errorf("plataforma não suportada: %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	archivePath := filepath.Join(home, ".hubsaude", "jre_download"+ext)

	if err := downloadFile(archivePath, downloadURL); err != nil {
		return "", fmt.Errorf("falha ao baixar o JRE: %v", err)
	}

	fmt.Println("[Provisioner] Download concluído! Extraindo arquivos...")

	if err := extractArchive(archivePath, hubsaudeDir); err != nil {
		return "", err
	}
	os.Remove(archivePath)

	return findJavaExe(hubsaudeDir)
}

func getAdoptiumURL() (string, string) {
	osName := runtime.GOOS
	if osName == "darwin" {
		osName = "mac"
	}
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "x64"
	} else if arch == "arm64" {
		arch = "aarch64"
	}

	url := fmt.Sprintf("https://api.adoptium.net/v3/binary/latest/21/ga/%s/%s/jre/hotspot/normal/eclipse", osName, arch)
	ext := ".tar.gz"
	if osName == "windows" {
		ext = ".zip"
	}
	return url, ext
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status HTTP inválido: %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func extractArchive(archivePath, targetDir string) error {
	cmd := exec.Command("tar", "-xf", archivePath, "-C", targetDir)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao extrair %s: %v. Saída: %s", archivePath, err, string(out))
	}
	return nil
}

func findJavaExe(dir string) (string, error) {
	var javaPath string
	targetName := "java"
	if runtime.GOOS == "windows" {
		targetName = "java.exe"
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == targetName {
			if filepath.Base(filepath.Dir(path)) == "bin" {
				javaPath = path
				return filepath.SkipDir
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}
	if javaPath != "" {
		return javaPath, nil
	}
	return "", fmt.Errorf("executável '%s' não encontrado dentro de %s", targetName, dir)
}
