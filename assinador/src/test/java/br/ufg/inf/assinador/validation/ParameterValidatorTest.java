package br.ufg.inf.assinador.validation;

import br.ufg.inf.assinador.dto.ValidationDetail;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;

import static org.assertj.core.api.Assertions.assertThat;

class ParameterValidatorTest {

    private ParameterValidator validator;

    @TempDir
    Path tempDir;

    @BeforeEach
    void setUp() {
        validator = new ParameterValidator(new ObjectMapper());
    }

    private Path writeFile(String name, String content) throws IOException {
        Path path = tempDir.resolve(name);
        Files.writeString(path, content);
        return path;
    }

    // ---- sign ----

    @Test
    void signComParametrosValidosNaoRetornaErros() throws IOException {
        Path input = writeFile("entrada.json", "{\"resourceType\":\"Bundle\"}");

        List<ValidationDetail> details =
                validator.validateSign(input.toString(), tempDir.resolve("saida.json").toString());

        assertThat(details).isEmpty();
    }

    @Test
    void signSemInputReportaRequired() {
        List<ValidationDetail> details =
                validator.validateSign(null, tempDir.resolve("saida.json").toString());

        assertThat(details).contains(new ValidationDetail("input", "required"));
    }

    @Test
    void signComInputInexistenteReportaNotFound() {
        String missing = tempDir.resolve("nao_existe.json").toString();

        List<ValidationDetail> details =
                validator.validateSign(missing, tempDir.resolve("saida.json").toString());

        assertThat(details).contains(new ValidationDetail("input", "not_found"));
    }

    @Test
    void signComExtensaoInvalidaReportaInvalidFormat() throws IOException {
        Path input = writeFile("entrada.txt", "{}");

        List<ValidationDetail> details =
                validator.validateSign(input.toString(), tempDir.resolve("saida.json").toString());

        assertThat(details).contains(new ValidationDetail("input", "invalid_format"));
    }

    @Test
    void signComJsonInvalidoReportaInvalidJson() throws IOException {
        Path input = writeFile("entrada.json", "isto não é json");

        List<ValidationDetail> details =
                validator.validateSign(input.toString(), tempDir.resolve("saida.json").toString());

        assertThat(details).contains(new ValidationDetail("input", "invalid_json"));
    }

    @Test
    void signSemOutputReportaRequired() throws IOException {
        Path input = writeFile("entrada.json", "{}");

        List<ValidationDetail> details = validator.validateSign(input.toString(), "  ");

        assertThat(details).contains(new ValidationDetail("output", "required"));
    }

    @Test
    void signComDiretorioDeSaidaInexistenteReportaOutputDirNotFound() throws IOException {
        Path input = writeFile("entrada.json", "{}");
        String output = tempDir.resolve("inexistente").resolve("saida.json").toString();

        List<ValidationDetail> details = validator.validateSign(input.toString(), output);

        assertThat(details).contains(new ValidationDetail("output", "output_dir_not_found"));
    }

    // ---- validate ----

    @Test
    void validateComAssinaturaValidaNaoRetornaErros() throws IOException {
        Path signature = writeFile("assinatura.json", "{\"resourceType\":\"Signature\"}");

        List<ValidationDetail> details = validator.validateValidate(signature.toString());

        assertThat(details).isEmpty();
    }

    @Test
    void validateSemSignatureReportaRequired() {
        List<ValidationDetail> details = validator.validateValidate("");

        assertThat(details).containsExactly(new ValidationDetail("signature", "required"));
    }

    @Test
    void validateComArquivoInexistenteReportaNotFound() {
        List<ValidationDetail> details =
                validator.validateValidate(tempDir.resolve("nao_existe.json").toString());

        assertThat(details).containsExactly(new ValidationDetail("signature", "not_found"));
    }

    @Test
    void validateComExtensaoInvalidaReportaInvalidFormat() throws IOException {
        Path signature = writeFile("assinatura.xml", "{}");

        List<ValidationDetail> details = validator.validateValidate(signature.toString());

        assertThat(details).containsExactly(new ValidationDetail("signature", "invalid_format"));
    }

    @Test
    void validateComJsonInvalidoReportaInvalidJson() throws IOException {
        Path signature = writeFile("assinatura.json", "{ quebrado");

        List<ValidationDetail> details = validator.validateValidate(signature.toString());

        assertThat(details).containsExactly(new ValidationDetail("signature", "invalid_json"));
    }
}
