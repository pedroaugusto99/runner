# Matriz de Rastreabilidade

## Objetivo

Este documento relaciona as histórias de usuário definidas na [especificacao upstream fixada](https://github.com/kyriosdata/runner/blob/main/especificacao.md) com o estado atual do repositório, indicando evidências objetivas, lacunas e próximos passos. O objetivo é facilitar acompanhamento, planejamento e demonstração de aderência aos requisitos sem alterar a implementação atual.

## Legenda de status

- **Não iniciado**: não há evidência de implementação no repositório.
- **Em preparação**: existem artefatos de base, documentação ou automação parcial.
- **Parcial**: parte dos critérios de aceitação já é atendida.
- **Concluído**: todos os critérios de aceitação estão cobertos por implementação e evidências.

## Estado observado em 2026-06-02

Validações executadas localmente:

- `env GOCACHE=/tmp/runner-go-build-cache GOMODCACHE=/tmp/runner-go-mod-cache go test ./...`, em `assinatura/`
- `env GOCACHE=/tmp/runner-go-build-cache GOMODCACHE=/tmp/runner-go-mod-cache go vet ./...`, em `assinatura/`
- `env GOCACHE=/tmp/runner-go-build-cache GOMODCACHE=/tmp/runner-go-mod-cache go test ./...`, em `simulador/`
- `env GOCACHE=/tmp/runner-go-build-cache GOMODCACHE=/tmp/runner-go-mod-cache go vet ./...`, em `simulador/`
- `mvn -Dmaven.repo.local=/tmp/runner-m2 test`, em `assinador/`

Resultado: os comandos passaram. A evidência Java ainda é limitada porque há apenas teste de carga de contexto; os testes de contrato com o JAR real ainda precisam ser adicionados.

## Visão geral

| História | Título | Status atual | Observação resumida |
| --- | --- | --- | --- |
| US-01 | Invocar `assinador.jar` via CLI | Parcial | CLI possui `sign` e `validate` em modo local, com execução do `assinador.jar`; modo HTTP e gestão de servidor ainda não existem. |
| US-02 | Simular assinatura digital com validação de parâmetros | Parcial | O Java já possui comandos locais `sign` e `validate` com simulação simples, mas sem contrato JSON, validação FHIR rigorosa, HTTP ou PKCS#11. |
| US-03 | Gerenciar ciclo de vida do Simulador do HubSaúde | Em preparação | Binário `simulador` existe como stub, sem start/stop/status nem download dinâmico. |
| US-04 | Provisionar JDK automaticamente | Não iniciado | Há diretório `assinatura/internal/jdk/`, mas ainda sem implementação. |
| US-05 | Disponibilizar binários multiplataforma | Parcial | Há workflows de build/release, mas ainda com cobertura e artefatos divergentes da especificação. |

## Detalhamento por história

### US-01 — Invocar `assinador.jar` via CLI

**Status:** Parcial

**Evidências atuais**
- Existe o componente `assinatura/` com CLI baseada em Cobra.
- Existe o comando `version` em `assinatura/cmd/version.go`, com suporte a versão e SHA curto por `-ldflags`.
- Existem comandos `sign` e `validate` em `assinatura/cmd/`.
- Existe integração local com o JAR em `assinatura/internal/assinador/`, preservando argumentos, `stdout`, `stderr`, timeout e código de saída.
- Existe descrição de propósito no comando raiz em `assinatura/cmd/root.go`.
- Existe pacote interno para integração com o `assinador.jar` (`assinatura/internal/assinador/`) e pacotes reservados para JDK e releases (`assinatura/internal/jdk/`, `assinatura/internal/release/`).

**Lacunas frente aos critérios de aceitação**
- Não há teste de contrato executando o JAR real a partir do CLI Go.
- Não há integração por HTTP com uma instância em modo servidor.
- Não há detecção de instância já em execução.
- Não há suporte a porta padrão/configurável, parada remota ou parada programada.
- A formatação de saída ainda depende do contrato atual do `assinador.jar`.

**Próximos passos sugeridos**
1. Fechar a interface pública do CLI (`sign`, `verify`, `server start`, `server stop`, `server status`) em issue ligada à US-01.
2. Ajustar o contrato de `docs/integracao-assinador.md` para coincidir com o Java atual ou alterar o Java para cumprir o contrato.
3. Criar testes de contrato que executem o JAR real via subprocesso.
4. Implementar modo HTTP após o `assinador.jar` expor endpoints de saúde, assinatura e validação.

### US-02 — Simular assinatura digital com validação de parâmetros

**Status:** Parcial

**Evidências atuais**
- Existe projeto Maven em `assinador/` com Spring Boot e picocli.
- Existem comandos locais `sign` e `validate`.
- `FakeSignatureService` simula criação de assinatura e validação baseada em arquivo `.json`.
- `mvn test` passa localmente com dependências resolvidas.

**Lacunas frente aos critérios de aceitação**
- A operação de validação no Java se chama `validate`, enquanto o contrato de integração usa `verify`.
- Não há saída JSON padronizada conforme `docs/integracao-assinador.md`.
- A validação de parâmetros ainda é mínima e baseada em existência/extensão de arquivo, não nas referências FHIR.
- Não há mensagens de erro estruturadas com código, campo e motivo.
- Não há separação clara entre resultado em `stdout` e diagnóstico em `stderr` em todos os casos.
- Não há testes funcionais para `sign`, `validate` e cenários negativos.
- Não há evidência de suporte a PKCS#11, ainda que simulado ou encapsulado.

**Próximos passos sugeridos**
1. Decidir e registrar se o comando final será `verify` ou `validate`; depois alinhar código, docs e testes.
2. Implementar saída JSON padronizada para sucesso e erro.
3. Extrair uma tabela mínima de validações FHIR e transformar cada regra em teste.
4. Adicionar testes unitários para `FakeSignatureService` e testes de CLI para `sign`/`verify`.
5. Criar simulação explícita de PKCS#11 no contrato antes de implementar integração real ou fake.

### US-03 — Gerenciar ciclo de vida do Simulador do HubSaúde

**Status:** Em preparação

**Evidências atuais**
- Existe o componente `simulador/`.
- O comando raiz do simulador já descreve a intenção de gerenciar o ciclo de vida.
- O arquivo `README.md` cita o simulador como stub usado para testes iniciais.

**Lacunas frente aos critérios de aceitação**
- Não há comandos de iniciar, parar e consultar status.
- Não há checagem de portas antes da inicialização.
- Não há download dinâmico da release mais recente do `simulador.jar`.
- Não há lógica para reutilizar artefato já existente localmente.
- Não há evidência de gerenciamento de processo em execução.

**Próximos passos sugeridos**
1. Definir comandos `start`, `stop` e `status` em issue ligada à US-03.
2. Criar pacote responsável por resolver e baixar artefatos do simulador via GitHub Releases.
3. Adicionar verificação de portas, health check e readiness antes de declarar sucesso.
4. Padronizar diretório de cache local para binários e metadados.
5. Testar porta ocupada, simulador ausente e reuso do artefato já baixado.

### US-04 — Provisionar JDK automaticamente

**Status:** Não iniciado

**Evidências atuais**
- Existe o diretório `assinatura/internal/jdk/`, indicando a intenção de encapsular a funcionalidade.

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
5. Adicionar erro amigável quando o JDK não puder ser obtido.

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
- Não há testes automatizados suficientes no repositório para sustentar as etapas `go test ./...` dos módulos Go com cobertura útil.

**Próximos passos sugeridos**
1. Incluir build e testes Java no CI, além de `go test`.
2. Ajustar a convenção de nomes dos artefatos publicados.
3. Expandir a release para incluir também o binário `simulador` quando aplicável.
4. Garantir publicação explícita dos arquivos `.sig` e `.pem` por artefato.
5. Registrar no repositório o procedimento de verificação dos artefatos assinados.

## Próximos passos priorizados

1. **P0: alinhar contrato CLI/JAR.** Resolver a divergência `verify` vs. `validate`, definir formato JSON final, códigos de erro, exit codes e porta padrão. Registrar decisões não óbvias como ADR curto.
2. **P0: criar testes de contrato.** Executar o JAR real via subprocesso para sucesso e erro de `sign`/`verify`, preservando argumentos com espaços e acentos.
3. **P1: consolidar o `assinador.jar` local.** Trocar mensagens livres por respostas estruturadas, fortalecer validações FHIR mínimas e separar `stdout` de `stderr`.
4. **P1: integrar `assinatura` ao modo local.** Implementar `assinatura sign` e `assinatura verify` chamando o JAR, com falha clara para JAR/JDK ausente.
5. **P2: adicionar modo HTTP.** Implementar health check real, start idempotente, status, stop e timeout/conexão recusada/resposta malformada.
6. **P2: colocar Java no CI.** O pipeline deve rodar Go e Maven em Linux e Windows para comprovar portabilidade real.
7. **P3: evoluir simulador e JDK.** Implementar ciclo de vida do `simulador.jar`, download condicional e provisionamento automático do JDK.

## Observações gerais

- O projeto está bem posicionado para a fase de implementação, mas ainda em estágio inicial de entrega funcional.
- A maior oportunidade imediata está em consolidar contratos, critérios verificáveis e automação aderente à especificação.
- Esta matriz deve ser atualizada a cada incremento relevante, preferencialmente junto com testes e evidências de uso.
