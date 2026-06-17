package cmd

import (
	"testing"
)

func TestValidateViaHTTP_US016(t *testing.T) {
	_, err := validateViaHTTP(nil, "9999", "assinatura.json")

	if err == nil {
		t.Error("Esperava erro de conexão ao tentar validar em porta sem servidor")
	}
}
