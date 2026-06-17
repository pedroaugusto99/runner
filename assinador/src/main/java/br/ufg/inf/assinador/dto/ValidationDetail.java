package br.ufg.inf.assinador.dto;

/**
 * Detalhe estruturado de uma falha de validação de parâmetro.
 *
 * @param field  nome do parâmetro inválido (ex.: "input", "signature")
 * @param reason motivo legível por máquina (ex.: "required", "not_found",
 *               "invalid_format", "invalid_json")
 */
public record ValidationDetail(String field, String reason) {
}
