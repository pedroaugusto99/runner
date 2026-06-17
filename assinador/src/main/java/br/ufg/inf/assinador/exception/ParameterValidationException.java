package br.ufg.inf.assinador.exception;

import br.ufg.inf.assinador.dto.ValidationDetail;

import java.util.List;

public class ParameterValidationException extends RuntimeException {

    private final transient List<ValidationDetail> details;

    public ParameterValidationException(String message, List<ValidationDetail> details) {
        super(message);
        this.details = List.copyOf(details);
    }

    public List<ValidationDetail> getDetails() {
        return details;
    }
}
