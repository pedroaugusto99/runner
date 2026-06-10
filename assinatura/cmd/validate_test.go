package cmd

import (
	"testing"
)

// Rastreabilidade: US-01.6 - Invocar assinador.jar via HTTP
// Issue Relacionada: #34
func TestValidateViaHTTP_US016(t *testing.T) {
	err := validateViaHTTP(nil, "9999", "assinatura.json")

	if err == nil {
		t.Error("Esperava erro de conexão ao tentar validar em porta sem servidor")
	}
}
