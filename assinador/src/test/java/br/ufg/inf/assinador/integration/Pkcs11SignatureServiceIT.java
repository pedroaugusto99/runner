package br.ufg.inf.assinador.integration;

import br.ufg.inf.assinador.exception.CryptoIntegrationException;
import br.ufg.inf.assinador.service.Pkcs11SignatureService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.test.context.ActiveProfiles;

import java.io.File;
import java.nio.charset.StandardCharsets;

import static org.junit.jupiter.api.Assertions.*;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

@SpringBootTest
@ActiveProfiles("test")
class Pkcs11SignatureServiceIT {

    @Autowired
    private Pkcs11SignatureService signatureService;

    @BeforeEach
    void setUp() {
        File softHsmLib = new File("/usr/lib/softhsm/libsofthsm2.so");
        assumeTrue(softHsmLib.exists(), "Abortando teste de integração: SoftHSM2 não instalado localmente.");
    }

    @Test
    void deveAssinarDadosComSucessoQuandoSoftHSM2EstiverDisponivel() {
        byte[] dadosParaAssinar = "Payload de teste do HubSaude".getBytes(StandardCharsets.UTF_8);

        try {
            String assinaturaBase64 = signatureService.assinarDados(dadosParaAssinar);

            assertNotNull(assinaturaBase64, "A assinatura retornada não deve ser nula.");
            assertFalse(assinaturaBase64.isBlank(), "A assinatura retornada não deve estar vazia.");

            assertDoesNotThrow(() -> java.util.Base64.getDecoder().decode(assinaturaBase64),
                    "O resultado retornado deve ser uma string Base64 válida.");

        } catch (Exception e) {
            if (e instanceof CryptoIntegrationException && e.getMessage().contains("DISPOSITIVO INDISPONÍVEL")) {
                logTestWarning("O driver do SoftHSM2 existe, mas o HubSaudeToken ou a chave não foram provisionados.");
                assumeTrue(false, "Ignorando teste: Dispositivo presente, mas massa criptográfica ausente.");
            } else {
                fail("Falha inesperada no teste de integração com o barramento PKCS11: " + e.getMessage(), e);
            }
        }
    }

    private void logTestWarning(String mensagem) {
        System.out.println("[PKCS11-IT-WARN] " + mensagem);
    }
}