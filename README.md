# Projeto Runner & Assinador Jar (FHIR/HubSaúde)

## Projeto da Disciplina de Implementação e Integração (2026/1)

### Grupo - Pedro e Claudio

<hr>



## Estado Atual do Projeto
O projeto evoluiu da Fase de Fundação para uma entrega funcional robusta. O Sistema Runner agora opera como um orquestrador completo, capaz de gerir o ciclo de vida do componente Java (assinador.jar), oferecendo modos de execução Local (via subprocesso) e Servidor (via HTTP/REST) com fallback automático, eliminando latências de inicialização da JVM (cold start).

### O que já foi feito (Sprint 1 a 3)
- [x] **Configuração de Ambiente**: Definição do módulo Go e ambiente de desenvolvimento.
- [x] **Estrutura de Pacotes (DT-06)**: Implementação do layout de diretórios separando binários (`cmd/`) de lógica interna (`internal/`).
- [x] **CLI de Assinatura (US-01.1)**: Criação do ponto de entrada do assinador com suporte ao comando `version`.
- [x] **Stub do Simulador (T-01.1.4)**: Implementação do binário inicial do simulador com framework Cobra.
- [x] **Base do assinador Java (US-02)**: Criação do projeto Maven com comandos locais iniciais `sign` e `validate`.
- [x] **Versionamento Semântico**: Preparação da variável global `version` para injeção dinâmica via pipeline.
- [x] **Automação (US-05.1)**: Configuração do GitHub Actions para Cross-Compilation (Windows, Linux, macOS).
- [x] **Segurança (US-05.3)**: Implementação da assinatura de artefatos com Cosign.
- [x] [US-01.5 — Iniciar assinador.jar no modo servidor](https://github.com/kyriosdata/runner/blob/802d241630ab3eac231834bc6c8afdd948c56856/docs/plano-revisitado-v2.md#us-015--iniciar-assinadorjar-no-modo-servidor)
- [x] [US-01.6 — Invocar assinador.jar via HTTP](https://github.com/kyriosdata/runner/blob/802d241630ab3eac231834bc6c8afdd948c56856/docs/plano-revisitado-v2.md#us-016--invocar-assinadorjar-via-http)
- [x] [US-01.7 — Detectar instância do assinador.jar em execução](https://github.com/kyriosdata/runner/blob/802d241630ab3eac231834bc6c8afdd948c56856/docs/plano-revisitado-v2.md#us-017--detectar-inst%C3%A2ncia-do-assinadorjar-em-execu%C3%A7%C3%A3o)
- [x] [US-01.8 — Interromper execução do assinador.jar](https://github.com/kyriosdata/runner/blob/802d241630ab3eac231834bc6c8afdd948c56856/docs/plano-revisitado-v2.md#us-018--interromper-execu%C3%A7%C3%A3o-do-assinadorjar)
- [x] [US-01.9 — Agendar interrupção do assinador.jar por inatividade](https://github.com/kyriosdata/runner/blob/802d241630ab3eac231834bc6c8afdd948c56856/docs/plano-revisitado-v2.md#us-019--agendar-interrup%C3%A7%C3%A3o-do-assinadorjar-por-inatividade)
- [x] [US-02.4 — Endpoints HTTP do assinador.jar](https://github.com/kyriosdata/runner/blob/802d241630ab3eac231834bc6c8afdd948c56856/docs/plano-revisitado-v2.md#us-024--endpoints-http-do-assinadorjar)
- [ ] [US-02.5 — Integração com dispositivo criptográfico via PKCS#11](https://github.com/kyriosdata/runner/blob/802d241630ab3eac231834bc6c8afdd948c56856/docs/plano-revisitado-v2.md#us-025--integra%C3%A7%C3%A3o-com-dispositivo-criptogr%C3%A1fico-via-pkcs11) - **Em andamento**

### O que falta fazer (Próximos Passos)
- [x] **Validação de parâmetros (US-02.2/02.3)**: validação rigorosa no `assinador.jar` (autoridade única), com erros estruturados (`campo`/`motivo`) e resultado determinístico de validação.
- [x] **Saída do contrato CLI/JAR (US-01.4)**: envelope JSON em `stdout`, exibição legível no CLI e códigos de saída coerentes — ver [ADR 0001](docs/adr/0001-contrato-cli-jar.md).
- [ ] **Nomenclatura `verify`/`validate`**: unificar o nome do comando (mantido `validate` por ora).
- [ ] **JDK e simulador**: download automático do JDK e ciclo de vida do `simulador.jar`.

## Estrutura do Repositório
O projeto separa os componentes principais em diretórios próprios:
- `assinatura/`: Código fonte do CLI de assinatura.
- `assinador/`: Código fonte do componente Java (Maven).
- `simulador/`: Código fonte do CLI de simulação.
- `docs/`: Documentação específica desta implementação.

## Como gerar os executáveis
Pré-requisitos: Go 1.25+, JDK 21 e Maven.

```bash
# 1. Empacotar o assinador.jar (artefato único, com Main-Class)
cd assinador
mvn package            # gera assinador/target/assinador.jar

# 2. Compilar o CLI de assinatura
cd ../assinatura
go build -o assinatura .
```

O CLI localiza o `assinador.jar` automaticamente (procurando em
`assinador/target/assinador.jar` a partir do diretório atual ou do binário);
também aceita `--jar <caminho>` ou a variável `ASSINADOR_JAR`.

## Como executar

```bash
cd assinatura

# Modo servidor (padrão): inicia o assinador.jar em background
./assinatura start

# Assinar — por padrão usa HTTP; faz fallback automático para subprocesso local
./assinatura sign --input entrada.json --output assinatura.json

# Forçar modo local (subprocesso java -jar), sem servidor
./assinatura sign --input entrada.json --output assinatura.json --local

# Validar assinatura
./assinatura validate --signature assinatura.json

# Parar servidor
./assinatura stop
```

A saída é estruturada e legível (✔/✖) e o código de saída reflete o resultado
(sucesso, assinatura inválida ou erro de parâmetro). A validação de parâmetros é
responsabilidade do `assinador.jar` (autoridade única) — ver
[ADR 0001](docs/adr/0001-contrato-cli-jar.md). Use `--help` em qualquer comando
para ver exemplos.

## Como executar os testes

```bash
# CLI e pacotes Go (use -short para pular os testes de integração com o JAR)
cd assinatura && go test ./...

# Componente Java (assinador.jar)
cd ../assinador && mvn test
```

Os testes de integração (`assinatura/cmd/integration_test.go`) exercitam o
contrato real CLI ↔ `assinador.jar` por subprocesso e por HTTP; são
automaticamente ignorados quando o JAR ou o `java` não estão disponíveis.

## Como contribuir
Trabalhe em branches curtas, abra PRs pequenos ligados a issues que referenciam
as histórias de usuário, e garanta `go test ./...` e `mvn test` verdes antes do
merge. Decisões não óbvias devem virar um ADR curto em `docs/adr/`.

Goiânia, 2026
