# ADR 0001 — Contrato CLI ↔ assinador.jar (sign/validate)

- **Status:** Aceito
- **Data:** 2026-06-16
- **Histórias relacionadas:** US-01.2, US-01.3, US-01.4, US-02.2, US-02.3
- **Especificação (commit fixo):** https://github.com/kyriosdata/runner/blob/802d241630ab3eac231834bc6c8afdd948c56856/especificacao.md

## Contexto

Ao implementar a validação de parâmetros (US-02.2/02.3) e a integração do CLI
(US-01.2/01.3/01.4), foi preciso fechar o contrato entre o CLI `assinatura` (Go)
e o `assinador.jar` (Java). O documento `docs/integracao-assinador.md` descreve um
contrato-alvo mais rico (payload inline, `contentType`, `signatureFormat`,
keystore PKCS#11) que ainda não está implementado, e havia ambiguidades reais:
interface de arquivos vs. payload; `verify` vs. `validate`; formato da saída.

## Decisão

1. **Interface baseada em arquivos.** `sign --input <f> --output <f>` e
   `validate --signature <f>`. O contrato rico de `docs/integracao-assinador.md`
   permanece como evolução futura.
2. **O `assinador.jar` é a autoridade única de validação.** Toda validação de
   parâmetros (presença, existência, extensão `.json`, JSON parseável) vive em
   `ParameterValidator` e é compartilhada pelos comandos picocli e pelos
   endpoints HTTP. **O CLI não replica** essa validação (sem `MarkFlagRequired`,
   sem checagem de presença no Go) — ele apenas repassa os argumentos e renderiza
   o resultado.
3. **Saída estruturada em `stdout`, diagnóstico em `stderr`.** O jar emite um
   envelope JSON (uma linha, UTF-8) em `stdout`; logs/diagnósticos usam slf4j
   (`stderr`). O CLI Go faz parse do envelope e renderiza apresentação legível,
   idêntica nos modos local e HTTP.

   ```json
   { "success": true, "operation": "sign", "signature": "...", "output": "...", "metadata": {...} }
   { "success": true, "operation": "validate", "valid": true|false, "message": "...", "metadata": {...} }
   { "success": false, "operation": "sign|validate", "error": { "code": "VALIDATION_ERROR", "message": "...", "details": [ { "field": "input", "reason": "required" } ] } }
   ```

4. **Códigos de saída.** `sign`: `0` sucesso · `1` erro de validação · `2` erro
   interno. `validate`: `0` válida · `1` inválida · `2` erro de validação de
   parâmetro · `3` erro interno. O CLI propaga o código.
5. **Nome do comando: `validate`** (não `verify`), alinhado às histórias e ao
   código existente. A unificação `verify`/`validate` segue como item de contrato
   em aberto.
6. **Resultado determinístico de `validate`:** válida quando o JSON tem
   `resourceType == "Signature"` (formato produzido por `sign`); senão inválida.

## Consequências

- A validação roda no caminho padrão (HTTP) e no fallback local com a mesma
  regra e a mesma mensagem — sem divergência entre Go e Java.
- Erros do usuário (`VALIDATION_ERROR`, exit 1/2) ficam distintos de erros de
  sistema (`INTERNAL_ERROR`/infraestrutura).
- Este ADR **refina** `docs/integracao-assinador.md` no que tange ao formato de
  entrada e saída efetivamente implementado; o contrato rico daquele documento
  permanece como meta futura.
- A divergência `verify` vs. `validate` continua registrada como pendência.
