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
- [ ] **Contrato CLI/JAR**: Alinhar comandos, JSON de saída, códigos de erro e nomenclatura `verify`/`validate`.
- [ ] **JDK e simulador**: Adicionar download automático do JDK e ciclo de vida do `simulador.jar`.

## Estrutura do Repositório
O projeto separa os componentes principais em diretórios próprios:
- `assinatura/`: Código fonte do CLI de assinatura.
- `assinador/`: Código fonte do componente Java (Maven).
- `simulador/`: Código fonte do CLI de simulação.
- `docs/`: Documentação específica desta implementação.

## Como executar localmente
Certifique-se de ter o Go 1.25+ e Maven instalados.

```bash
# 1. Iniciar servidor em background
cd assinatura
./assinatura start

# 2. Assinar (o CLI decide automaticamente entre HTTP ou Local)
./assinatura sign --input entrada.json --output assinatura.json

# 3. Validar assinatura
./assinatura validate --signature assinatura.json

# 4. Parar servidor
./assinatura stop

# Para testar o CLI assinatura
go test -v ./...

# Para executar o stub do simulador
cd ../simulador
go run .

# Para testar o assinador Java
cd ../assinador
mvn test
```
Goiânia, 2026
