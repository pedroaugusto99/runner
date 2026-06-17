package br.ufg.inf.assinador.validation;

import br.ufg.inf.assinador.dto.ValidationDetail;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.stereotype.Component;

import java.io.File;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;

@Component
public class ParameterValidator {

    private static final String JSON_EXTENSION = ".json";

    private final ObjectMapper objectMapper;

    public ParameterValidator(ObjectMapper objectMapper) {
        this.objectMapper = objectMapper;
    }

    public List<ValidationDetail> validateSign(String input, String output) {
        List<ValidationDetail> details = new ArrayList<>();
        validateJsonFileParam("input", input, details);
        validateOutputParam("output", output, details);
        return details;
    }

    public List<ValidationDetail> validateValidate(String signature) {
        List<ValidationDetail> details = new ArrayList<>();
        validateJsonFileParam("signature", signature, details);
        return details;
    }

    private void validateJsonFileParam(String field, String value, List<ValidationDetail> details) {
        if (isBlank(value)) {
            details.add(new ValidationDetail(field, "required"));
            return;
        }

        File file = new File(value);
        if (!file.exists() || !file.isFile()) {
            details.add(new ValidationDetail(field, "not_found"));
            return;
        }

        if (!value.toLowerCase().endsWith(JSON_EXTENSION)) {
            details.add(new ValidationDetail(field, "invalid_format"));
            return;
        }

        try {
            objectMapper.readTree(file);
        } catch (Exception e) {
            details.add(new ValidationDetail(field, "invalid_json"));
        }
    }

    private void validateOutputParam(String field, String value, List<ValidationDetail> details) {
        if (isBlank(value)) {
            details.add(new ValidationDetail(field, "required"));
            return;
        }

        Path parent = Paths.get(value).getParent();
        if (parent != null && !parent.toFile().isDirectory()) {
            details.add(new ValidationDetail(field, "output_dir_not_found"));
        }
    }

    private static boolean isBlank(String value) {
        return value == null || value.isBlank();
    }
}
