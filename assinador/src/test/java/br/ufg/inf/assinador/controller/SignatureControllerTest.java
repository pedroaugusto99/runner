package br.ufg.inf.assinador.controller;

import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.http.MediaType;
import org.springframework.test.web.servlet.MockMvc;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Map;

import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@SpringBootTest
@AutoConfigureMockMvc
class SignatureControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @TempDir
    Path tempDir;

    private Path writeFile(String name, String content) throws IOException {
        Path path = tempDir.resolve(name);
        Files.writeString(path, content);
        return path;
    }

    private String json(Map<String, String> body) throws IOException {
        return objectMapper.writeValueAsString(body);
    }

    @Test
    void healthRetornaUp() throws Exception {
        mockMvc.perform(get("/api/v1/health"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.status").value("UP"));
    }

    @Test
    void signComParametrosValidosRetorna200() throws Exception {
        Path input = writeFile("bundle.json", "{\"resourceType\":\"Bundle\"}");
        String body = json(Map.of(
                "inputPath", input.toString(),
                "outputPath", tempDir.resolve("assinatura.json").toString()));

        mockMvc.perform(post("/api/v1/sign").contentType(MediaType.APPLICATION_JSON).content(body))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.success").value(true))
                .andExpect(jsonPath("$.operation").value("sign"))
                .andExpect(jsonPath("$.signature").exists());
    }

    @Test
    void signSemInputRetorna400ComErroEstruturado() throws Exception {
        String body = json(Map.of("outputPath", tempDir.resolve("assinatura.json").toString()));

        mockMvc.perform(post("/api/v1/sign").contentType(MediaType.APPLICATION_JSON).content(body))
                .andExpect(status().isBadRequest())
                .andExpect(jsonPath("$.success").value(false))
                .andExpect(jsonPath("$.error.code").value("VALIDATION_ERROR"))
                .andExpect(jsonPath("$.error.details[0].field").value("input"))
                .andExpect(jsonPath("$.error.details[0].reason").value("required"));
    }

    @Test
    void validateComAssinaturaValidaRetornaValid() throws Exception {
        Path signature = writeFile("assinatura.json", "{\"resourceType\":\"Signature\"}");
        String body = json(Map.of("signaturePath", signature.toString()));

        mockMvc.perform(post("/api/v1/validate").contentType(MediaType.APPLICATION_JSON).content(body))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.success").value(true))
                .andExpect(jsonPath("$.valid").value(true));
    }

    @Test
    void validateComConteudoNaoReconhecidoRetornaInvalid() throws Exception {
        Path signature = writeFile("assinatura.json", "{\"resourceType\":\"Bundle\"}");
        String body = json(Map.of("signaturePath", signature.toString()));

        mockMvc.perform(post("/api/v1/validate").contentType(MediaType.APPLICATION_JSON).content(body))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.valid").value(false));
    }
}
