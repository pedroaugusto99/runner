package cmd

import (
	"testing"
)

// Rastreabilidade: US-01.6 - Invocar assinador.jar via HTTP
// Issue Relacionada: #34
func TestSignFallbackLocal_US016(t *testing.T) {
	t.Run("Quando servidor esta offline, deve sinalizar modo local", func(t *testing.T) {
		port := "9999"

		isUp := isServerUp(port)
		if isUp {
			t.Errorf("Esperava que servidor na porta %s estivesse offline", port)
		}
	})
}
