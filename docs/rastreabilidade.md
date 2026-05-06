# Matriz de Rastreabilidade

## Objetivo

Este documento relaciona as histórias de usuário definidas em `docs/especificacao.md` com o estado atual do repositório, indicando evidências objetivas, lacunas e próximos passos. O objetivo é facilitar acompanhamento, planejamento e demonstração de aderência aos requisitos sem alterar a implementação atual.

## Legenda de status

- **Não iniciado**: não há evidência de implementação no repositório.
- **Em preparação**: existem artefatos de base, documentação ou automação parcial.
- **Parcial**: parte dos critérios de aceitação já é atendida.
- **Concluído**: todos os critérios de aceitação estão cobertos por implementação e evidências.

## Visão geral

| História | Título | Status atual | Observação resumida |
| --- | --- | --- | --- |
| US-01 | Invocar `assinador.jar` via CLI | Em preparação | CLI base existe, mas sem comandos de assinatura/validação e sem integração com Java. |
| US-02 | Simular assinatura digital com validação de parâmetros | Não iniciado | Projeto `assinador/` existe, porém sem implementação funcional observável do JAR. |
| US-03 | Gerenciar ciclo de vida do Simulador do HubSaúde | Em preparação | Binário `simulador` existe como stub, sem start/stop/status nem download dinâmico. |
| US-04 | Provisionar JDK automaticamente | Não iniciado | Há diretório `internal/jdk/`, mas ainda sem código. |
| US-05 | Disponibilizar binários multiplataforma | Parcial | Há workflows de build/release, mas ainda com cobertura e artefatos divergentes da especificação. |

## Detalhamento por história

### US-01 — Invocar `assinador.jar` via CLI

**Status:** Em preparação

**Evidências atuais**
- Existe o binário `cmd/assinatura/` com CLI baseada em Cobra.
- Existe o comando `version` em `cmd/assinatura/cmd/version.go`.
- Existe descrição de propósito no comando raiz em `cmd/assinatura/cmd/root.go`.

**Lacunas frente aos critérios de aceitação**
- Não há comandos para criar assinatura ou validar assinatura.
- Não há integração com `assinador.jar` por linha de comando.
- Não há integração por HTTP com uma instância em modo servidor.
- Não há detecção de instância já em execução.
- Não há suporte a porta padrão/configurável, parada remota ou parada programada.
- Não há formatação de saída de operações de assinatura/validação.

**Próximos passos sugeridos**
1. Definir a interface do CLI (`sign`, `verify`, `server start`, `server stop`, `server status`).
2. Formalizar contrato de invocação local e HTTP do `assinador.jar`.
3. Implementar um pacote `internal/invoker` com abstração por modo (`local` e `http`).
4. Criar testes de aceitação cobrindo os fluxos principais e de erro.

### US-02 — Simular assinatura digital com validação de parâmetros

**Status:** Não iniciado

**Evidências atuais**
- Existe a pasta `assinador/`, indicando reserva de espaço para o componente Java.
- A especificação já define claramente o comportamento esperado de simulação e validação rigorosa.

**Lacunas frente aos critérios de aceitação**
- Não foi localizada implementação funcional do `assinador.jar`.
- Não há validação de parâmetros conforme as referências FHIR.
- Não há simulação de criação/validação com respostas pré-construídas.
- Não há mensagens de erro estruturadas.
- Não há evidência de suporte a PKCS#11, ainda que simulado ou encapsulado.

**Próximos passos sugeridos**
1. Definir contrato de entrada/saída do `assinador.jar` antes da implementação.
2. Criar tabela de validações obrigatórias baseada nas referências da seção 10 da especificação.
3. Implementar respostas simuladas determinísticas para facilitar testes de integração.
4. Padronizar códigos de erro e mensagens legíveis para CLI e HTTP.

### US-03 — Gerenciar ciclo de vida do Simulador do HubSaúde

**Status:** Em preparação

**Evidências atuais**
- Existe o binário `cmd/simulador/`.
- O comando raiz do simulador já descreve a intenção de gerenciar o ciclo de vida.
- O arquivo `README.md` cita o simulador como stub usado para testes iniciais.

**Lacunas frente aos critérios de aceitação**
- Não há comandos de iniciar, parar e consultar status.
- Não há checagem de portas antes da inicialização.
- Não há download dinâmico da release mais recente do `simulador.jar`.
- Não há lógica para reutilizar artefato já existente localmente.
- Não há evidência de gerenciamento de processo em execução.

**Próximos passos sugeridos**
1. Definir comandos `start`, `stop` e `status` no CLI.
2. Criar pacote responsável por resolver e baixar artefatos do simulador.
3. Adicionar verificação de portas e de processo local.
4. Padronizar diretório de cache local para binários e metadados.

### US-04 — Provisionar JDK automaticamente

**Status:** Não iniciado

**Evidências atuais**
- Existe o diretório `internal/jdk/`, indicando a intenção de encapsular a funcionalidade.

**Lacunas frente aos critérios de aceitação**
- Não há detecção da presença do JDK nem validação de versão.
- Não há download automático por plataforma.
- Não há disponibilização do JDK baixado para uso pelo Runner.
- Não há evidência de compatibilidade entre Linux, Windows e macOS.

**Próximos passos sugeridos**
1. Definir a versão mínima de JDK exigida pelo sistema.
2. Escolher a origem oficial dos downloads por plataforma.
3. Implementar resolução de ambiente local antes de cair para download.
4. Padronizar cache e política de reaproveitamento do JDK já baixado.

### US-05 — Disponibilizar binários multiplataforma

**Status:** Parcial

**Evidências atuais**
- Existe `/.github/workflows/build.yml` com jobs de teste e cross-compilation.
- Existe `/.github/workflows/release.yml` com geração de artefatos, checksums e uso de Cosign.
- O projeto já usa versionamento via tag em release.

**Lacunas frente aos critérios de aceitação**
- Os workflows atuais geram apenas artefatos de `assinatura`, não de `simulador`.
- Os nomes dos artefatos atuais não seguem integralmente o formato exemplificado na especificação.
- A especificação exige publicação de `<artefato>`, `<artefato>.sig` e `<artefato>.pem` para cada binário; o workflow atual precisa evidenciar esse empacotamento final.
- A especificação menciona formatos como `.AppImage` e `.dmg`, enquanto o workflow atual gera binários crus.
- Não há testes automatizados no repositório para sustentar a etapa `go test ./...` com cobertura útil.

**Próximos passos sugeridos**
1. Ajustar a convenção de nomes dos artefatos publicados.
2. Expandir a release para incluir também o binário `simulador` quando aplicável.
3. Garantir publicação explícita dos arquivos `.sig` e `.pem` por artefato.
4. Registrar no repositório o procedimento de verificação dos artefatos assinados.

## Observações gerais

- O projeto está bem posicionado para a fase de implementação, mas ainda em estágio inicial de entrega funcional.
- A maior oportunidade imediata está em consolidar contratos, critérios verificáveis e automação aderente à especificação.
- Esta matriz deve ser atualizada a cada incremento relevante, preferencialmente junto com testes e evidências de uso.
