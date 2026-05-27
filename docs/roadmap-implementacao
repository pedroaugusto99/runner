# Roadmap de Implementação

## Objetivo

Este documento organiza a evolução do Sistema Runner em incrementos pequenos, verificáveis e compatíveis com o estado atual do repositório. O foco é transformar a especificação em um plano executável, reduzindo risco e facilitando acompanhamento do time.

## Premissas

- O estado atual do projeto ainda é inicial, com CLI base em Go e workflows de build/release já criados.
- O componente `assinador.jar` ainda precisa ser implementado.
- O sistema deve evoluir de forma incremental, preservando a possibilidade de validação frequente.
- Sempre que possível, cada incremento deve gerar evidências objetivas: código, testes, documentação ou automação.

## Prioridades

| Prioridade | Tema | Motivo |
| --- | --- | --- |
| P0 | Contratos e comportamento esperado | Evita retrabalho entre Go e Java. |
| P1 | `assinador.jar` com simulação determinística | Destrava integração e testes. |
| P1 | CLI `assinatura` com fluxo local mínimo | Permite primeira entrega funcional fim a fim. |
| P2 | Modo servidor HTTP | Melhora desempenho e cobre parte central da US-01. |
| P2 | Gestão do `simulador.jar` | Expande escopo funcional do Runner. |
| P3 | Provisionamento automático de JDK | Melhora experiência do usuário final. |
| P3 | Empacotamento final multiplataforma | Fecha distribuição conforme especificação. |

## Estratégia de execução

A recomendação é trabalhar em ondas curtas, cada uma com um objetivo verificável e baixo acoplamento. O ideal é que cada etapa termine com:

- código funcional mínimo;
- teste automatizado ou evidência operacional;
- documentação atualizada;
- definição clara do próximo passo.

## Ondas de implementação

### Onda 1 — Fechar contratos e critérios verificáveis

**Objetivo**

Eliminar ambiguidade antes de implementar os componentes principais.

**Entradas**
- `docs/especificacao.md`
- `docs/rastreabilidade.md`
- `docs/integracao-assinador.md`

**Entregas esperadas**
- Lista fechada de comandos do CLI.
- Lista mínima de parâmetros de `sign` e `verify`.
- Definição oficial da porta padrão do modo servidor.
- Convenção de saída (`json`, `text`, `pretty`) documentada.
- Tabela de erros padronizados para Go e Java.

**Critério de pronto**
- O time consegue implementar Go e Java de forma desacoplada, sem decisões pendentes sobre o contrato principal.

### Onda 2 — Implementar núcleo do `assinador.jar`

**Objetivo**

Entregar o primeiro comportamento funcional do lado Java com simulação consistente.

**Entregas esperadas**
- Comando `sign` em modo local.
- Comando `verify` em modo local.
- Validação mínima de parâmetros obrigatórios.
- Respostas simuladas determinísticas.
- Códigos de saída consistentes para sucesso e erro.

**Sugestão de tarefas**
1. Criar estrutura do projeto Java em `assinador/`.
2. Implementar parser de argumentos.
3. Criar camada de validação de parâmetros.
4. Implementar respostas simuladas de assinatura e validação.
5. Padronizar saída estruturada em JSON.
6. Cobrir cenários de erro básicos com testes.

**Critério de pronto**
- É possível executar `sign` e `verify` via terminal e obter resposta previsível para cenários válidos e inválidos.

### Onda 3 — Integrar CLI `assinatura` ao modo local

**Objetivo**

Entregar a primeira integração fim a fim entre CLI Go e `assinador.jar`.

**Entregas esperadas**
- Comando `assinatura sign` funcional em modo local.
- Comando `assinatura verify` funcional em modo local.
- Flags mínimas para entrada e saída.
- Mensagens de erro amigáveis no CLI.
- Testes de integração básicos chamando o JAR localmente.

**Sugestão de tarefas**
1. Criar comandos Cobra para `sign` e `verify`.
2. Implementar pacote `internal/invoker` para execução de processo externo.
3. Traduzir flags do CLI para parâmetros do JAR.
4. Mapear erros técnicos para erros de uso.
5. Adicionar testes de fumaça do CLI.

**Critério de pronto**
- Um usuário consegue disparar assinatura e validação pelo binário Go sem invocar Java manualmente.

### Onda 4 — Adicionar modo servidor HTTP

**Objetivo**

Cobrir o fluxo preferencial da US-01 com menor custo de inicialização por requisição.

**Entregas esperadas**
- `assinador.jar` executando em modo servidor HTTP.
- Endpoint de saúde.
- Endpoint de `sign`.
- Endpoint de `verify`.
- Endpoint de desligamento controlado.
- Detecção de instância já ativa pelo CLI.

**Sugestão de tarefas**
1. Adicionar modo `server` ao `assinador.jar`.
2. Implementar endpoints definidos no contrato.
3. Criar cliente HTTP em `internal/invoker`.
4. Implementar `assinatura server start|status|stop`.
5. Adicionar testes de integração do modo HTTP.

**Critério de pronto**
- O CLI usa uma instância em execução e consegue operar sem `cold start` a cada chamada.

### Onda 5 — Estruturar gerenciamento do `simulador.jar`

**Objetivo**

Entregar a base da US-03 com lifecycle management controlado pelo Runner.

**Entregas esperadas**
- Comandos `simulador start`, `simulador stop` e `simulador status`.
- Resolução da release mais recente do simulador.
- Download condicional do artefato quando ausente.
- Verificação de portas antes de iniciar.
- Armazenamento local padronizado do binário e metadados.

**Sugestão de tarefas**
1. Criar pacote de download/resolução em `internal/release`.
2. Criar camada de processo local para start/stop/status.
3. Implementar cache local do artefato.
4. Adicionar validação de portas e mensagens de erro adequadas.

**Critério de pronto**
- O usuário consegue subir e parar o simulador sem conhecer comandos Java nem localizar manualmente o JAR.

### Onda 6 — Provisionar JDK automaticamente

**Objetivo**

Reduzir dependência de instalação manual de Java no ambiente do usuário.

**Entregas esperadas**
- Detecção de JDK local compatível.
- Estratégia de fallback para download por plataforma.
- Cache reaproveitável do JDK.
- Uso transparente do JDK pelo Runner.

**Sugestão de tarefas**
1. Definir versão mínima suportada.
2. Implementar detector de JDK local em `internal/jdk`.
3. Implementar resolução de origem/download por sistema operacional.
4. Integrar o caminho do JDK ao invocador local.
5. Cobrir cenários de ausência e reaproveitamento.

**Critério de pronto**
- O usuário consegue usar o sistema suportado mesmo sem Java previamente instalado, desde que o download seja possível.

### Onda 7 — Fechar distribuição, qualidade e release

**Objetivo**

Aproximar o repositório do modelo de entrega exigido na especificação.

**Entregas esperadas**
- Cobertura mínima de testes para os fluxos principais.
- Versionamento refletido nos binários.
- Ajuste dos nomes de artefato conforme convenção final.
- Inclusão de `simulador` na pipeline de release, quando aplicável.
- Evidência explícita de publicação de arquivos `.sig` e `.pem`.
- Guia de uso e validação de artefatos atualizado.

**Critério de pronto**
- Uma release publicada demonstra aderência funcional mínima e conformidade operacional com a especificação.

## Dependências principais

| Item | Depende de |
| --- | --- |
| Integração do CLI Go | Contrato estável + `assinador.jar` local funcional |
| Modo HTTP | Contrato estável + base local do `assinador.jar` |
| Gestão do simulador | Resolução de release + gestão de processos |
| Provisionamento de JDK | Definição de versão mínima + estratégia por plataforma |
| Release final | Fluxos principais funcionando + convenção de artefatos decidida |

## Ordem mínima recomendada

1. Fechar decisões de contrato.
2. Entregar `assinador.jar` local com simulação.
3. Integrar `assinatura` ao modo local.
4. Adicionar modo servidor HTTP.
5. Implementar lifecycle do simulador.
6. Implementar provisionamento de JDK.
7. Refinar qualidade, empacotamento e release.

## Riscos e mitigação

| Risco | Impacto | Mitigação |
| --- | --- | --- |
| Contrato indefinido entre Go e Java | Alto retrabalho | Validar cedo o documento de integração com o professor/time. |
| Parâmetros de assinatura pouco claros | Implementação errada | Extrair tabela explícita de parâmetros das referências FHIR. |
| Provisionamento de JDK complexo demais cedo | Atraso nas entregas | Postergar para depois da integração local mínima. |
| Pipeline divergir da especificação | Release inconsistente | Revisar a checklist de release a cada incremento. |
| Falta de testes de regressão | Quebras silenciosas | Adicionar smoke tests desde as primeiras ondas. |

## Checklist de avanço por onda

Use a lista abaixo para controlar o progresso do roadmap:

- [ ] Onda 1 concluída
- [ ] Onda 2 concluída
- [ ] Onda 3 concluída
- [ ] Onda 4 concluída
- [ ] Onda 5 concluída
- [ ] Onda 6 concluída
- [ ] Onda 7 concluída

## Resultado esperado

Ao seguir este roadmap, o projeto evolui de um repositório com fundação pronta para um sistema funcional, testável e distribuível, mantendo coerência com `docs/especificacao.md` e com os artefatos de documentação já adicionados ao repositório.
