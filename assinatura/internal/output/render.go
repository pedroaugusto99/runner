package output

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

var reasonMessages = map[string]string{
	"required":             "parâmetro obrigatório ausente",
	"not_found":            "arquivo não encontrado",
	"invalid_format":       "formato inválido (esperado um arquivo .json)",
	"invalid_json":         "conteúdo não é um JSON válido",
	"output_dir_not_found": "diretório de saída inexistente",
}

var reasonHints = map[string]string{
	"required":             "informe o parâmetro e tente novamente",
	"not_found":            "verifique se o caminho do arquivo está correto",
	"invalid_format":       "use um arquivo com extensão .json",
	"invalid_json":         "corrija o conteúdo JSON do arquivo",
	"output_dir_not_found": "crie o diretório ou ajuste o caminho de saída",
}

func Render(w io.Writer, raw []byte) (*Result, bool) {
	res, ok := extractEnvelope(raw)
	if !ok {
		if trimmed := strings.TrimSpace(string(raw)); trimmed != "" {
			fmt.Fprintln(w, trimmed)
		}
		return nil, false
	}

	switch {
	case res.Error != nil:
		renderError(w, res.Error)
	case res.Operation == "validate":
		renderValidate(w, res)
	default:
		renderSign(w, res)
	}
	return res, true
}

func extractEnvelope(raw []byte) (*Result, bool) {
	scanner := bufio.NewScanner(bytes.NewReader(raw))
	scanner.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if !strings.HasPrefix(line, "{") || !strings.HasSuffix(line, "}") {
			continue
		}

		var probe map[string]json.RawMessage
		if err := json.Unmarshal([]byte(line), &probe); err != nil {
			continue
		}
		_, hasOperation := probe["operation"]
		_, hasSuccess := probe["success"]
		if !hasOperation && !hasSuccess {
			continue
		}

		var res Result
		if err := json.Unmarshal([]byte(line), &res); err != nil {
			continue
		}
		return &res, true
	}
	return nil, false
}

func renderSign(w io.Writer, res *Result) {
	fmt.Fprintln(w, "✔ Assinatura criada com sucesso")
	if res.Output != "" {
		fmt.Fprintf(w, "  arquivo: %s\n", res.Output)
	}
	if simulated(res) {
		fmt.Fprintln(w, "  modo:    simulação")
	}
}

func renderValidate(w io.Writer, res *Result) {
	if res.Valid != nil && *res.Valid {
		fmt.Fprintln(w, "✔ Assinatura VÁLIDA")
	} else {
		fmt.Fprintln(w, "✖ Assinatura INVÁLIDA")
	}
	if res.Message != "" {
		fmt.Fprintf(w, "  %s\n", res.Message)
	}
}

func renderError(w io.Writer, e *ResultError) {
	fmt.Fprintln(w, "✖ Operação rejeitada")
	if e.Message != "" {
		fmt.Fprintf(w, "  %s\n", e.Message)
	}
	for _, d := range e.Details {
		fmt.Fprintf(w, "  • %s: %s\n", d.Field, reasonMessage(d.Reason))
		if hint := reasonHints[d.Reason]; hint != "" {
			fmt.Fprintf(w, "    → %s\n", hint)
		}
	}
}

func reasonMessage(reason string) string {
	if msg, ok := reasonMessages[reason]; ok {
		return msg
	}
	return reason
}

func simulated(res *Result) bool {
	if res.Metadata == nil {
		return false
	}
	v, ok := res.Metadata["simulated"].(bool)
	return ok && v
}
