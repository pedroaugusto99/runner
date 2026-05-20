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
