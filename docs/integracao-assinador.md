# Contrato de Integração — `assinatura` ↔ `assinador.jar`

## Objetivo

Este documento define o contrato esperado entre o CLI `assinatura` e o componente Java `assinador.jar`, conforme descrito em `docs/especificacao.md`. O foco é estabelecer uma base estável para implementação posterior, sem depender ainda da lógica final do assinador.

## Princípios

- O fluxo funcional deve ser o mesmo nos modos **local** e **HTTP**.
- O CLI deve apresentar mensagens legíveis ao usuário final.
- O `assinador.jar` deve ser a autoridade de validação dos parâmetros de negócio.
- Erros devem ser estruturados e previsíveis.
- As respostas de simulação devem ser determinísticas para facilitar testes.

## Operações suportadas

| Operação | Objetivo | Comando lógico no CLI |
| --- | --- | --- |
| Criar assinatura | Simular geração de assinatura digital | `assinatura sign` |
| Validar assinatura | Simular validação de assinatura digital | `assinatura verify` |
| Iniciar servidor | Subir `assinador.jar` em modo HTTP | `assinatura server start` |
| Parar servidor | Solicitar parada do `assinador.jar` | `assinatura server stop` |
| Consultar servidor | Verificar disponibilidade/status | `assinatura server status` |

## Modos de integração

### 1. Modo local (invocação direta)

O CLI executa o JAR diretamente, por exemplo:

```bash
java -jar assinador.jar sign ...
java -jar assinador.jar verify ...
```

**Características esperadas**
- Cada chamada executa um processo independente.
- Adequado para automação pontual e scripts.
- O CLI traduz flags e argumentos para o formato esperado pelo JAR.
- O resultado é retornado via `stdout` em formato estruturado.
- Erros são retornados via `stderr` e código de saída diferente de zero.

### 2. Modo servidor (HTTP)

O `assinador.jar` permanece em execução localmente e recebe requisições HTTP do CLI.

**Características esperadas**
- O CLI deve preferir este modo por padrão, salvo orientação contrária.
- O CLI deve conseguir detectar uma instância já existente na porta configurada.
- O protocolo deve ser local (`localhost`) por padrão.
- As respostas devem usar JSON.

## Contrato lógico de entrada

Os nomes de campos abaixo representam o contrato de integração entre componentes. A lista exata de parâmetros de negócio deve ser refinada a partir das referências FHIR citadas na seção 10 de `docs/especificacao.md`.

### Operação `sign`

**Entrada lógica mínima**

```json
{
  "payload": "<conteudo-ou-caminho>",
  "contentType": "application/fhir+json",
  "keystore": {
    "type": "pkcs11",
    "library": "/caminho/ou/id/do-driver",
    "slot": "0",
    "pin": "******"
  },
  "signatureFormat": "XAdES|JAdES|CAdES|PAdES",
  "output": "json"
}
```

**Regras gerais esperadas**
- Campos obrigatórios devem ser rejeitados quando ausentes ou vazios.
- Tipos inválidos devem gerar erro de validação.
- Combinações incompatíveis de parâmetros devem gerar erro claro.
- O retorno de sucesso deve ser estável para uma mesma entrada válida dentro do modo de simulação.

### Operação `verify`

**Entrada lógica mínima**

```json
{
  "signedPayload": "<conteudo-assinado-ou-caminho>",
  "contentType": "application/fhir+json",
  "signatureFormat": "XAdES|JAdES|CAdES|PAdES",
  "output": "json"
}
```

**Regras gerais esperadas**
- Deve validar presença e formato mínimo do conteúdo assinado.
- Deve retornar resultado determinístico de validação no modo simulado.
- Deve informar claramente o motivo de erro de validação de entrada.

## Contrato lógico de saída

### Resposta de sucesso — `sign`

```json
{
  "success": true,
  "operation": "sign",
  "mode": "local|http",
  "signature": "<assinatura-simulada>",
  "algorithm": "<algoritmo-ou-placeholder>",
  "format": "<formato>",
  "metadata": {
    "timestamp": "2026-05-05T10:00:00Z",
    "simulated": true
  }
}
```

### Resposta de sucesso — `verify`

```json
{
  "success": true,
  "operation": "verify",
  "mode": "local|http",
  "valid": true,
  "message": "Assinatura válida no modo de simulação.",
  "metadata": {
    "timestamp": "2026-05-05T10:00:00Z",
    "simulated": true
  }
}
```

### Resposta de erro padronizada

```json
{
  "success": false,
  "operation": "sign|verify|server",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Parâmetro obrigatório ausente: payload.",
    "details": [
      {
        "field": "payload",
        "reason": "required"
      }
    ]
  }
}
```

## Códigos de erro sugeridos

| Código | Quando usar |
| --- | --- |
| `VALIDATION_ERROR` | Parâmetro ausente, vazio, inválido ou combinação inconsistente |
| `JDK_NOT_AVAILABLE` | Java/JDK não encontrado ou incompatível |
| `SIGNER_NOT_FOUND` | `assinador.jar` não encontrado |
| `SIGNER_START_FAILED` | Falha ao iniciar em modo servidor |
| `SIGNER_UNREACHABLE` | Instância HTTP não responde |
| `INTERNAL_ERROR` | Falha inesperada não classificada |

## Contrato HTTP proposto

### Base URL padrão

```text
http://127.0.0.1:8090
```

### Endpoints sugeridos

| Método | Rota | Objetivo |
| --- | --- | --- |
| `POST` | `/api/v1/sign` | Solicitar criação de assinatura simulada |
| `POST` | `/api/v1/verify` | Solicitar validação simulada |
| `GET` | `/api/v1/health` | Verificar disponibilidade do servidor |
| `POST` | `/api/v1/shutdown` | Solicitar parada controlada |

### Exemplo de requisição `POST /api/v1/sign`

```json
{
  "payload": "{\"resourceType\":\"Bundle\"}",
  "contentType": "application/fhir+json",
  "signatureFormat": "JAdES",
  "output": "json"
}
```

### Exemplo de resposta `GET /api/v1/health`

```json
{
  "status": "ok",
  "service": "assinador.jar",
  "mode": "http",
  "version": "dev"
}
```

## Mapeamento esperado no CLI

| Ação no CLI | Modo local | Modo HTTP |
| --- | --- | --- |
| `assinatura sign` | Executa `java -jar ... sign ...` | Faz `POST /api/v1/sign` |
| `assinatura verify` | Executa `java -jar ... verify ...` | Faz `POST /api/v1/verify` |
| `assinatura server status` | Não se aplica | Faz `GET /api/v1/health` |
| `assinatura server stop` | Não se aplica | Faz `POST /api/v1/shutdown` |

## Comportamentos operacionais esperados

- O CLI deve normalizar parâmetros e validar o mínimo necessário antes de invocar o assinador.
- O `assinador.jar` deve concentrar a validação de negócio e dos parâmetros específicos da assinatura.
- O CLI deve transformar erros técnicos em mensagens amigáveis, preservando detalhes quando solicitado.
- O mesmo erro lógico deve produzir o mesmo código em ambos os modos.
- O modo HTTP deve ser preferencial quando houver instância acessível e o usuário não solicitar modo local.

## Itens a decidir antes da implementação

- Lista final de parâmetros obrigatórios das operações `sign` e `verify`.
- Convenção para diferenciar conteúdo inline de caminho de arquivo.
- Porta padrão oficial do modo servidor.
- Estratégia de autenticação local para comando de desligamento, se necessária.
- Formato final das saídas em terminal (`json`, `text`, `pretty`).

## Resultado esperado deste contrato

Com este documento, o time pode implementar o cliente Go e o servidor/CLI Java de forma desacoplada, usando um acordo comum de comandos, estruturas de dados, erros e comportamento esperado.
