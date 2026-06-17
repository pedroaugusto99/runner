package br.ufg.inf.assinador.dto;

import com.fasterxml.jackson.annotation.JsonInclude;

import java.util.List;

/**
 * Bloco de erro padronizado do envelope de saída.
 *
 * @param code    código estável do erro (ex.: "VALIDATION_ERROR", "INTERNAL_ERROR")
 * @param message mensagem legível para o usuário
 * @param details lista de campos inválidos e respectivos motivos
 */
@JsonInclude(JsonInclude.Include.NON_NULL)
public record OperationError(String code, String message, List<ValidationDetail> details) {
}
