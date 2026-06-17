package br.ufg.inf.assinador.service;

import br.ufg.inf.assinador.dto.OperationResult;
import br.ufg.inf.assinador.dto.ValidationDetail;
import br.ufg.inf.assinador.exception.ParameterValidationException;
import br.ufg.inf.assinador.validation.ParameterValidator;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.stereotype.Service;

import java.io.File;
import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.List;

@Service
public class FakeSignatureService implements SignatureService {

    private static final String SIMULATED_SIGNATURE = """
            {
              "resourceType": "Signature",
              "status": "simulated",
              "type": [
                {
                  "system": "urn:iso-astm:E1762-95:2013",
                  "code": "1.2.840.10065.1.12.1.1"
                }
              ],
              "data": "YXNzaW5hdHVyYSBzaW11bGFkYSBkZSBleGVtcGxv"
            }
            """;

    private final ParameterValidator validator;
    private final ObjectMapper objectMapper;

    public FakeSignatureService(ParameterValidator validator, ObjectMapper objectMapper) {
        this.validator = validator;
        this.objectMapper = objectMapper;
    }

    @Override
    public OperationResult sign(String inputPath, String outputPath) {
        rejectIfInvalid(validator.validateSign(inputPath, outputPath));

        try {
            Files.writeString(Paths.get(outputPath), SIMULATED_SIGNATURE);
        } catch (IOException e) {
            throw new UncheckedIOException("Falha ao gravar a assinatura simulada em: " + outputPath, e);
        }

        return OperationResult.signSuccess(SIMULATED_SIGNATURE.strip(), outputPath);
    }

    @Override
    public OperationResult validate(String signaturePath) {
        rejectIfInvalid(validator.validateValidate(signaturePath));

        if (isSignatureResource(signaturePath)) {
            return OperationResult.validateResult(true, "Assinatura válida no modo de simulação.");
        }
        return OperationResult.validateResult(false,
                "Assinatura inválida: 'resourceType' diferente de 'Signature'.");
    }

    private boolean isSignatureResource(String signaturePath) {
        try {
            String resourceType = objectMapper.readTree(new File(signaturePath))
                    .path("resourceType")
                    .asText("");
            return "Signature".equals(resourceType);
        } catch (IOException e) {
            return false;
        }
    }

    private void rejectIfInvalid(List<ValidationDetail> details) {
        if (details.isEmpty()) {
            return;
        }
        ValidationDetail first = details.get(0);
        String message = String.format("Parâmetro inválido: %s (%s).", first.field(), first.reason());
        throw new ParameterValidationException(message, details);
    }
}
