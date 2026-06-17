package br.ufg.inf.assinador.service;

import br.ufg.inf.assinador.dto.OperationResult;
import br.ufg.inf.assinador.dto.ValidationDetail;
import br.ufg.inf.assinador.exception.ParameterValidationException;
import br.ufg.inf.assinador.validation.ParameterValidator;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;

class FakeSignatureServiceTest {

    private FakeSignatureService service;

    @TempDir
    Path tempDir;

    @BeforeEach
    void setUp() {
        ObjectMapper objectMapper = new ObjectMapper();
        service = new FakeSignatureService(new ParameterValidator(objectMapper), objectMapper);
    }

    private Path writeFile(String name, String content) throws IOException {
        Path path = tempDir.resolve(name);
        Files.writeString(path, content);
        return path;
    }

    @Test
    void signComParametrosValidosGeraAssinaturaEArquivo() throws IOException {
        Path input = writeFile("bundle.json", "{\"resourceType\":\"Bundle\"}");
        Path output = tempDir.resolve("assinatura.json");

        OperationResult result = service.sign(input.toString(), output.toString());

        assertThat(result.success()).isTrue();
        assertThat(result.operation()).isEqualTo("sign");
        assertThat(result.signature()).contains("\"resourceType\": \"Signature\"");
        assertThat(result.output()).isEqualTo(output.toString());
        assertThat(Files.exists(output)).isTrue();
    }

    @Test
    void signComParametrosInvalidosLancaExcecaoAntesDeProcessar() {
        Path output = tempDir.resolve("assinatura.json");

        assertThatThrownBy(() -> service.sign(null, output.toString()))
                .isInstanceOf(ParameterValidationException.class)
                .satisfies(e -> {
                    ParameterValidationException pve = (ParameterValidationException) e;
                    assertThat(pve.getDetails()).contains(new ValidationDetail("input", "required"));
                });
        assertThat(Files.exists(output)).isFalse();
    }

    @Test
    void validateComAssinaturaReconhecidaRetornaValida() throws IOException {
        Path signature = writeFile("assinatura.json", "{\"resourceType\":\"Signature\"}");

        OperationResult result = service.validate(signature.toString());

        assertThat(result.success()).isTrue();
        assertThat(result.valid()).isTrue();
    }

    @Test
    void validateComConteudoNaoReconhecidoRetornaInvalida() throws IOException {
        Path signature = writeFile("assinatura.json", "{\"resourceType\":\"Bundle\"}");

        OperationResult result = service.validate(signature.toString());

        assertThat(result.success()).isTrue();
        assertThat(result.valid()).isFalse();
    }

    @Test
    void validateComParametroInvalidoLancaExcecao() {
        assertThatThrownBy(() -> service.validate("   "))
                .isInstanceOf(ParameterValidationException.class)
                .satisfies(e -> {
                    ParameterValidationException pve = (ParameterValidationException) e;
                    assertThat(pve.getDetails()).contains(new ValidationDetail("signature", "required"));
                });
    }
}
