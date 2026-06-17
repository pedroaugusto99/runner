package br.ufg.inf.assinador.dto;

import com.fasterxml.jackson.annotation.JsonInclude;

import java.time.Instant;
import java.util.LinkedHashMap;
import java.util.Map;

@JsonInclude(JsonInclude.Include.NON_NULL)
public record OperationResult(
        boolean success,
        String operation,
        String signature,
        Boolean valid,
        String message,
        String output,
        OperationError error,
        Map<String, Object> metadata) {

    private static Map<String, Object> simulatedMetadata() {
        Map<String, Object> metadata = new LinkedHashMap<>();
        metadata.put("simulated", true);
        metadata.put("timestamp", Instant.now().toString());
        return metadata;
    }

    public static OperationResult signSuccess(String signature, String outputPath) {
        return new OperationResult(true, "sign", signature, null, null, outputPath, null, simulatedMetadata());
    }

    public static OperationResult validateResult(boolean valid, String message) {
        return new OperationResult(true, "validate", null, valid, message, null, null, simulatedMetadata());
    }

    public static OperationResult error(String operation, OperationError error) {
        return new OperationResult(false, operation, null, null, null, null, error, null);
    }
}
