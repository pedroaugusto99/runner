package br.ufg.inf.assinador.controller;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;

import java.util.HashMap;
import java.util.Map;

@RestController
@RequestMapping("/api/v1")
public class SignatureController {

    private static final Logger log = LoggerFactory.getLogger(SignatureController.class);
    private static final String STATUS = "status";

    @GetMapping("/health")
    public ResponseEntity<Map<String, String>> healthCheck() {
        log.debug("Recebido pedido de health check (verificação de vitalidade).");
        Map<String, String> response = new HashMap<>();
        response.put(STATUS, "UP");
        response.put("service", "assinador-jar");
        return ResponseEntity.ok(response);
    }

    @PostMapping("/sign")
    public ResponseEntity<Map<String, Object>> signDocument(@RequestBody Map<String, Object> request) {
        log.info("Recebido pedido HTTP POST para /sign");

        Map<String, Object> response = new HashMap<>();
        response.put(STATUS, "success");
        response.put("message", "Assinatura digital simulada com sucesso via HTTP");

        return ResponseEntity.ok(response);
    }

    @PostMapping("/validate")
    public ResponseEntity<Map<String, Object>> validateSignature(@RequestBody Map<String, Object> request) {
        log.info("Recebido pedido HTTP POST para /validate");

        Map<String, Object> response = new HashMap<>();
        response.put(STATUS, "success");
        response.put("message", "Validação concluída: A assinatura fornecida é válida.");

        return ResponseEntity.ok(response);
    }
}