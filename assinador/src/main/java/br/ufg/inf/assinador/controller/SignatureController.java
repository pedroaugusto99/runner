package br.ufg.inf.assinador.controller;

import br.ufg.inf.assinador.dto.OperationResult;
import br.ufg.inf.assinador.dto.SignHttpRequest;
import br.ufg.inf.assinador.dto.ValidateHttpRequest;
import br.ufg.inf.assinador.dto.OperationError;
import br.ufg.inf.assinador.dto.ValidationDetail;
import br.ufg.inf.assinador.exception.CryptoIntegrationException;
import br.ufg.inf.assinador.service.Pkcs11SignatureService;
import br.ufg.inf.assinador.service.SignatureService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.nio.charset.StandardCharsets;
import java.util.Collections;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

@RestController
@RequestMapping("/api/v1")
public class SignatureController {

    private static final Logger log = LoggerFactory.getLogger(SignatureController.class);

    private final SignatureService signatureService;
    private final Pkcs11SignatureService pkcs11SignatureService;

    public SignatureController(SignatureService signatureService, Pkcs11SignatureService pkcs11SignatureService) {
        this.signatureService = signatureService;
        this.pkcs11SignatureService = pkcs11SignatureService;
    }

    @GetMapping("/health")
    public ResponseEntity<Map<String, String>> healthCheck() {
        log.debug("Recebido pedido de health check (verificação de vitalidade).");
        Map<String, String> response = new HashMap<>();
        response.put("status", "UP");
        response.put("service", "assinador-jar");
        return ResponseEntity.ok(response);
    }

    @PostMapping("/sign")
    public ResponseEntity<OperationResult> signDocument(@RequestBody SignHttpRequest request) {
        log.info("Recebido pedido HTTP POST para /sign");

        try {
            if (request == null || request.inputPath() == null) {
                List<ValidationDetail> details = List.of(new ValidationDetail("input", "required"));
                OperationError errorDto = new OperationError("VALIDATION_ERROR", "Parâmetro input obrigatório.", details);
                return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(OperationResult.error("sign", errorDto));
            }

            if (request.outputPath() == null) {
                List<ValidationDetail> details = List.of(new ValidationDetail("output", "required"));
                OperationError errorDto = new OperationError("VALIDATION_ERROR", "Parâmetro output obrigatório.", details);
                return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(OperationResult.error("sign", errorDto));
            }

            if (!request.inputPath().startsWith("memory://")) {
                java.io.File file = new java.io.File(request.inputPath());
                if (!file.exists() || request.inputPath().isBlank()) {
                    List<ValidationDetail> details = List.of(new ValidationDetail("input", "not_found"));
                    OperationError errorDto = new OperationError("VALIDATION_ERROR", "O ficheiro de entrada não foi encontrado", details);
                    return ResponseEntity.status(HttpStatus.BAD_REQUEST).body(OperationResult.error("sign", errorDto));
                }

                return ResponseEntity.ok(signatureService.sign(request.inputPath(), request.outputPath()));
            }

            String contentToSign = "dados-padrao-hubsaude";
            String assinaturaBase64 = pkcs11SignatureService.assinarDados(contentToSign.getBytes(StandardCharsets.UTF_8));

            OperationResult result = OperationResult.signSuccess(assinaturaBase64, request.outputPath());
            return ResponseEntity.ok(result);

        } catch (CryptoIntegrationException e) {
            log.error("Dispositivo PKCS#11 indisponível no endpoint /sign: {}", e.getMessage());
            OperationError errorDto = new OperationError("HSM_UNAVAILABLE", e.getMessage(), Collections.emptyList());
            return ResponseEntity.status(HttpStatus.SERVICE_UNAVAILABLE).body(OperationResult.error("sign", errorDto));
        } catch (Exception e) {
            log.error("Erro inesperado no endpoint /sign", e);
            OperationError errorDto = new OperationError("INTERNAL_ERROR", e.getMessage(), Collections.emptyList());
            return ResponseEntity.status(HttpStatus.INTERNAL_SERVER_ERROR).body(OperationResult.error("sign", errorDto));
        }
    }

    @PostMapping("/validate")
    public ResponseEntity<OperationResult> validateSignature(@RequestBody ValidateHttpRequest request) {
        log.info("Recebido pedido HTTP POST para /validate");
        OperationResult result = signatureService.validate(request.signaturePath());
        return ResponseEntity.ok(result);
    }
}