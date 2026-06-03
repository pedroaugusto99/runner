# Projeto Runner & Assinador Jar (FHIR/HubSaúde)

## Projeto da Disciplina de Implementação e Integração (2026/1)

### Grupo - Pedro e Claudio

<hr>



## Estado Atual do Projeto
Atualmente, o projeto está entre a **Fase de Fundação** e a primeira entrega funcional integrada. A estrutura base dos executáveis Go existe, o componente Java já possui comandos locais iniciais e o CLI `assinatura` consegue invocar o `assinador.jar` em modo local.

### O que já foi feito (Sprint 1)
- [x] **Configuração de Ambiente**: Definição do módulo Go e ambiente de desenvolvimento.
- [x] **Estrutura de Pacotes (DT-06)**: Implementação do layout de diretórios separando binários (`cmd/`) de lógica interna (`internal/`).
- [x] **CLI de Assinatura (US-01.1)**: Criação do ponto de entrada do assinador com suporte ao comando `version`.
- [x] **Stub do Simulador (T-01.1.4)**: Implementação do binário inicial do simulador com framework Cobra.
- [x] **Base do assinador Java (US-02)**: Criação do projeto Maven com comandos locais iniciais `sign` e `validate`.
- [x] **Versionamento Semântico**: Preparação da variável global `version` para injeção dinâmica via pipeline.
- [x] **Automação (US-05.1)**: Configuração do GitHub Actions para Cross-Compilation (Windows, Linux, macOS).
- [x] **Segurança (US-05.3)**: Implementação da assinatura de artefatos com Cosign.

### O que falta fazer (Próximos Passos)
- [ ] **Contrato CLI/JAR**: Alinhar comandos, JSON de saída, códigos de erro e nomenclatura `verify`/`validate`.
- [ ] **Testes de contrato**: Validar `assinatura` chamando o JAR real por subprocesso e, depois, por HTTP.
- [x] **Integração local (Sprint 3)**: Implementação inicial de `assinatura sign` e `assinatura validate` invocando `assinador.jar`.
- [ ] **Modo servidor**: Implementar health check, start/status/stop e auto-shutdown no modo HTTP.
- [ ] **JDK e simulador**: Adicionar download automático do JDK e ciclo de vida do `simulador.jar`.

## Estrutura do Repositório
O projeto segue o layout padrão para aplicações Go com múltiplos binários:
- `cmd/assinatura`: Código fonte do binário principal de assinatura.
- `cmd/simulador`: Código fonte do binário de simulação.
- `internal/`: Pacotes privados compartilhados entre os binários.
- `assinador/`: Código fonte do componente Java (Maven).

## Como executar localmente
Certifique-se de ter o **Go 1.25+** instalado.

```bash
# Para ver a versão da assinatura
go run ./cmd/assinatura version

# Para assinar em modo local depois de gerar assinador/target/assinador.jar
go run ./cmd/assinatura sign --input entrada.json --output assinatura.json --local

# Para validar em modo local
go run ./cmd/assinatura validate --signature assinatura.json --local

# Para executar o stub do simulador
go run ./cmd/simulador

# Para testar os binários Go
go test ./...

# Para testar o assinador Java
cd assinador
mvn test
```
Goiânia, 2026
