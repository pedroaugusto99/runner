package br.ufg.inf.assinador.controller;

import br.ufg.inf.assinador.dto.OperationResult;
import br.ufg.inf.assinador.dto.SignHttpRequest;
import br.ufg.inf.assinador.dto.ValidateHttpRequest;
import br.ufg.inf.assinador.service.SignatureService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/v1")
public class SignatureController {

    private static final Logger log = LoggerFactory.getLogger(SignatureController.class);

    private final SignatureService signatureService;

    public SignatureController(SignatureService signatureService) {
        this.signatureService = signatureService;
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
        OperationResult result = signatureService.sign(request.inputPath(), request.outputPath());
        return ResponseEntity.ok(result);
    }

    @PostMapping("/validate")
    public ResponseEntity<OperationResult> validateSignature(@RequestBody ValidateHttpRequest request) {
        log.info("Recebido pedido HTTP POST para /validate");
        OperationResult result = signatureService.validate(request.signaturePath());
        return ResponseEntity.ok(result);
    }
}
