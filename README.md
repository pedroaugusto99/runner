# Projeto Runner & Assinador Jar (FHIR/HubSaúde)

## Projeto da Disciplina de Implementação e Integração (2026/1)

### Grupo - Pedro e Claudio

<hr>



Goiânia, 2026

## Estado Atual do Projeto
Atualmente, o projeto encontra-se no final da **Fase de Fundação (Sprint 1)**. A estrutura base do código e os esqueletos dos executáveis foram estabelecidos e validados localmente.

### O que já foi feito (Sprint 1)
- [x] **Configuração de Ambiente**: Definição do módulo Go e ambiente de desenvolvimento.
- [x] **Estrutura de Pacotes (DT-06)**: Implementação do layout de diretórios separando binários (`cmd/`) de lógica interna (`internal/`).
- [x] **CLI de Assinatura (US-01.1)**: Criação do ponto de entrada do assinador com suporte ao comando `version`.
- [x] **Stub do Simulador (T-01.1.4)**: Implementação do binário inicial do simulador para testes de integração.
- [x] **Versionamento Semântico**: Preparação da variável global `version` para injeção via pipeline de CI/CD.

### O que falta fazer (Próximos Passos)
- [ ] **Automação (US-05.1)**: Configuração do GitHub Actions para Cross-Compilation (Windows, Linux, macOS).
- [ ] **Segurança (US-05.3)**: Implementação da assinatura de artefatos com Cosign.
- [ ] **Projeto Java (Sprint 2)**: Desenvolvimento do `assinador.jar` dentro da pasta `assinador/`.
- [ ] **Integração (Sprint 3)**: Lógica de download automático do JDK e do assinador.

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

# Para executar o stub do simulador
go run ./cmd/simulador
