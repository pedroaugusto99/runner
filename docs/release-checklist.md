# Checklist de Release

## Objetivo

Este checklist operacionaliza os requisitos de distribuição descritos em `docs/especificacao.md`, com foco em consistência, rastreabilidade e segurança dos artefatos publicados.

## Quando usar

Use este checklist sempre que uma versão do projeto for preparada para publicação via GitHub Releases.

## 1. Pré-release

- [ ] Confirmar a versão SemVer da release (ex.: `v1.0.0`).
- [ ] Garantir que a branch e a tag correspondem ao estado aprovado para publicação.
- [ ] Revisar `README.md` e documentação de uso para a versão corrente.
- [ ] Confirmar que mudanças relevantes estão refletidas em documentação técnica e de usuário.
- [ ] Verificar se os workflows de CI/CD estão verdes.
- [ ] Garantir que não há arquivos temporários ou artefatos locais indevidos incluídos no repositório.

## 2. Testes e validação

- [ ] Executar `go vet ./...`.
- [ ] Executar `go test ./...`.
- [ ] Validar execução básica dos binários CLI previstos para a release.
- [ ] Confirmar que a versão exibida pelos binários corresponde à tag da release.
- [ ] Registrar limitações conhecidas que não bloqueiam a publicação.

## 3. Artefatos esperados

Conforme a especificação, a release deve contemplar artefatos multiplataforma e os arquivos associados de assinatura.

### `assinatura`

- [ ] Publicar artefato para Windows `amd64`.
- [ ] Publicar artefato para Linux `amd64`.
- [ ] Publicar artefato para macOS `amd64`.
- [ ] Publicar `checksums` da release.

### `simulador`

- [ ] Publicar artefato para Windows `amd64`.
- [ ] Publicar artefato para Linux `amd64`.
- [ ] Publicar artefato para macOS `amd64`.
- [ ] Publicar `checksums` da release, caso separados por componente.

## 4. Convenção de nomes

Verificar aderência ao padrão definido ou acordado pelo time. A especificação exemplifica nomes como:

- `assinatura-1.0.0-windows-amd64.exe`
- `assinatura-1.0.0-linux-amd64.AppImage`
- `assinatura-1.0.0-macos-amd64.dmg`
- `simulador-1.0.0-windows-amd64.exe`
- `simulador-1.0.0-linux-amd64.AppImage`
- `simulador-1.0.0-macos-amd64.dmg`

Checklist:

- [ ] Os nomes dos arquivos incluem produto, versão, sistema operacional e arquitetura.
- [ ] O sufixo/extensão final bate com o formato que será realmente distribuído.
- [ ] Não existem nomes ambíguos ou inconsistentes entre workflow e documentação.

## 5. Integridade e assinatura com Cosign

Para cada artefato distribuído, a release deve conter:

- [ ] `<artefato>`
- [ ] `<artefato>.sig`
- [ ] `<artefato>.pem`

Checklist operacional:

- [ ] A assinatura foi gerada automaticamente pelo pipeline.
- [ ] O pipeline usou identidade baseada em OIDC.
- [ ] A assinatura foi registrada no transparency log do Sigstore.
- [ ] Todos os artefatos possuem seus arquivos `.sig` correspondentes.
- [ ] Todos os artefatos possuem seus arquivos `.pem` correspondentes.
- [ ] O arquivo de `checksums` também está assinado, se essa for a política adotada pelo time.

## 6. Verificação dos artefatos

Após a publicação, validar ao menos um artefato de ponta a ponta usando `cosign verify-blob`.

Exemplo:

```bash
cosign verify-blob \
  --certificate assinatura-1.0.0-linux-amd64.AppImage.pem \
  --signature assinatura-1.0.0-linux-amd64.AppImage.sig \
  assinatura-1.0.0-linux-amd64.AppImage
```

Checklist:

- [ ] O comando de verificação retorna sucesso.
- [ ] O certificado corresponde ao artefato verificado.
- [ ] A evidência de verificação foi registrada no processo de release.

## 7. Publicação no GitHub Releases

- [ ] Criar ou confirmar a tag da versão.
- [ ] Publicar a release com notas geradas ou revisadas manualmente.
- [ ] Anexar todos os artefatos previstos.
- [ ] Confirmar anexos de `.sig` e `.pem` para cada binário.
- [ ] Confirmar presença do arquivo de `checksums`.
- [ ] Confirmar que o texto da release comunica o escopo entregue e as limitações conhecidas.

## 8. Pós-release

- [ ] Baixar pelo menos um artefato publicado para validar disponibilidade.
- [ ] Conferir se os links da release estão funcionando.
- [ ] Atualizar quadro de acompanhamento do projeto, se houver.
- [ ] Registrar incidentes ou ajustes necessários para a próxima release.

## 9. Lacunas observadas no estado atual do repositório

Com base no estado atual do projeto, vale revisar antes da primeira release funcional:

- [ ] Os workflows atuais precisam refletir integralmente a convenção final de nomes exigida pela especificação.
- [ ] A pipeline atual gera artefatos de `assinatura`, mas a especificação também prevê artefatos de `simulador`.
- [ ] É preciso confirmar se os formatos finais serão binários crus ou empacotados como `.AppImage` e `.dmg`.
- [ ] É importante garantir evidência explícita de publicação dos arquivos `.sig` e `.pem` associados a cada artefato.

## Resultado esperado

Ao concluir este checklist, a release deve estar alinhada com os requisitos de multiplataforma, integridade, assinatura criptográfica e rastreabilidade definidos para o Sistema Runner.
