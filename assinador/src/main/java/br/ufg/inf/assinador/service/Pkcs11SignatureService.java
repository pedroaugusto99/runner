package br.ufg.inf.assinador.service;

import br.ufg.inf.assinador.exception.CryptoIntegrationException;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.core.io.ClassPathResource;
import org.springframework.stereotype.Service;

import java.io.File;
import java.io.IOException;
import java.security.*;
import java.util.Base64;

@Service
public class Pkcs11SignatureService {

    private static final Logger log = LoggerFactory.getLogger(Pkcs11SignatureService.class);
    private static final String PROVIDER_NAME = "SunPKCS11-SoftHSM2";

    @Value("${hubsaude.crypto.pkcs11.key-label}")
    private String keyLabel;

    @Value("${hubsaude.crypto.pkcs11.pin}")
    private String pin;

    public String assinarDados(byte[] dados) throws GeneralSecurityException {
        if (Security.getProvider(PROVIDER_NAME) != null) {
            Security.removeProvider(PROVIDER_NAME);
        }

        Provider baseProvider = Security.getProvider("SunPKCS11");
        if (baseProvider == null) {
            throw new NoSuchProviderException("ERRO DE INFRAESTRUTURA: O Provider SunPKCS11 nativo da JVM não está disponível.");
        }

        File tempCfg = null;
        try {
            try (java.io.InputStream is = new ClassPathResource("pkcs11.cfg").getInputStream()) {
                tempCfg = File.createTempFile("sunpkcs11-", ".cfg");
                tempCfg.deleteOnExit();
                java.nio.file.Files.copy(is, tempCfg.toPath(), java.nio.file.StandardCopyOption.REPLACE_EXISTING);
            }
        } catch (IOException e) {
            throw new CryptoIntegrationException("ERRO DE CONFIGURAÇÃO: Não foi possível ler o arquivo interno pkcs11.cfg para inicializar o subsistema.", e);
        }

        try {
            Provider providerConfigurado = baseProvider.configure(tempCfg.getAbsolutePath());
            Security.addProvider(providerConfigurado);

            KeyStore keyStore = KeyStore.getInstance("PKCS11", providerConfigurado);
            keyStore.load(null, pin.toCharArray());

            if (!keyStore.containsAlias(keyLabel)) {
                throw new CryptoIntegrationException("DISPOSITIVO INDISPONÍVEL: O token foi localizado, mas a chave de assinatura '" + keyLabel + "' não foi encontrada no SoftHSM2.");
            }

            PrivateKey privateKey = (PrivateKey) keyStore.getKey(keyLabel, null);
            Signature signature = Signature.getInstance("SHA256withRSA", providerConfigurado);
            signature.initSign(privateKey);
            signature.update(dados);

            byte[] assinaturaBruta = signature.sign();
            log.info("Assinatura digital PKCS#11 gerada com sucesso.");
            return Base64.getEncoder().encodeToString(assinaturaBruta);

        } catch (ProviderException e) {
            throw new CryptoIntegrationException("DISPOSITIVO INDISPONÍVEL: O dispositivo criptográfico (SoftHSM2) não foi detectado no sistema, está inacessível ou o PIN fornecido é inválido.", e);
        } catch (Exception e) {
            throw new CryptoIntegrationException("ERRO CRIPTOGRÁFICO: Falha operacional no barramento PKCS#11. Detalhes: " + e.getMessage(), e);
        }
    }
}