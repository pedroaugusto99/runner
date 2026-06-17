package br.ufg.inf.assinador.exception;

import br.ufg.inf.assinador.dto.OperationError;
import br.ufg.inf.assinador.dto.OperationResult;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.ExceptionHandler;
import org.springframework.web.bind.annotation.RestControllerAdvice;

@RestControllerAdvice
public class GlobalExceptionHandler {

    private static final Logger log = LoggerFactory.getLogger(GlobalExceptionHandler.class);

    @ExceptionHandler(ParameterValidationException.class)
    public ResponseEntity<OperationResult> handleValidation(ParameterValidationException ex) {
        log.info("Parâmetros inválidos rejeitados: {}", ex.getMessage());
        OperationError error = new OperationError("VALIDATION_ERROR", ex.getMessage(), ex.getDetails());
        return ResponseEntity.status(HttpStatus.BAD_REQUEST)
                .body(OperationResult.error("validation", error));
    }
}
