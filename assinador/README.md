# assinador.jar

Aplicação Java (Spring Boot) que **simula** operações de assinatura digital
do Sistema Runner (HubSaúde). Suporta dois modos de invocação:

- **CLI direto**: `java -jar assinador.jar <comando> [opções]`
- **Servidor HTTP**: `java -jar assinador.jar` (escuta em `:8080` por padrão)

## Pré-requisitos

- JDK 21+
- Maven 3.9+

## Como executar (desenvolvimento)

```bash
# A partir da raiz do módulo (esta pasta)
mvn spring-boot:run

# Rodar os testes
mvn test

# Construir o artefato final (gera target/assinador.jar)
mvn clean package
```

## Estrutura

```
src/main/java/br/ufg/inf/assinador/
├── AssinadorApplication.java   # entry point Spring Boot
├── cli/                        # runner para modo CLI direto
├── controller/                 # endpoints HTTP (modo servidor)
├── service/                    # regras de simulação (criação/validação)
├── dto/                        # requests/responses
└── exception/                  # tratamento global de erros
```
## FakeSignature

```

package br.ufg.inf.assinador.service;

import org.springframework.stereotype.Service;
import java.io.File;

@Service
public class FakeSignatureService implements SignatureService {

    @Override
    public String sign(String inputFilePath) {
        File file = new File(inputFilePath);
        if (!file.exists() || !file.isFile()) {
            throw new IllegalArgumentException("O ficheiro de entrada não foi encontrado ou é inválido: " + inputFilePath);
        }

        return """
               {
                 "resourceType": "Signature",
                 "status": "simulated",
                 "type": [
                   {
                     "system": "urn:iso-astm:E1762-95:2013",
                     "code": "1.2.840.10065.1.12.1.1"
                   }
                 ],
                 "data": "YXNzaW5hdHVyYSBzaW11bGFkYSBkZSBleGVtcGxv"
               }
               """;
    }

    @Override
    public boolean validate(String signatureFilePath) {
        File file = new File(signatureFilePath);
        if (!file.exists() || !file.isFile()) {
            throw new IllegalArgumentException("O ficheiro de assinatura não foi encontrado: " + signatureFilePath);
        }

        return signatureFilePath.toLowerCase().endsWith(".json");
    }
}

```
