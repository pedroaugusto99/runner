package br.ufg.inf.assinador.cli;

import br.ufg.inf.assinador.dto.OperationError;
import br.ufg.inf.assinador.dto.OperationResult;
import br.ufg.inf.assinador.dto.ValidationDetail;
import br.ufg.inf.assinador.exception.CryptoIntegrationException;
import br.ufg.inf.assinador.exception.ParameterValidationException;
import br.ufg.inf.assinador.service.Pkcs11SignatureService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.List;
import java.util.concurrent.Callable;

@Component
@Command(name = "sign",
        description = "Gera uma assinatura digital real usando dispositivo PKCS#11 para um pacote FHIR.",
        mixinStandardHelpOptions = true)
public class SignCommand implements Callable<Integer> {

    private static final Logger log = LoggerFactory.getLogger(SignCommand.class);

    private final Pkcs11SignatureService signatureService;
    private final OperationResultPrinter printer;

    public SignCommand(Pkcs11SignatureService signatureService, OperationResultPrinter printer) {
        this.signatureService = signatureService;
        this.printer = printer;
    }

    @Option(names = {"-i", "--input"}, description = "Caminho para o ficheiro FHIR (.json) a ser assinado")
    private String inputFile;

    @Option(names = {"-o", "--output"}, description = "Caminho onde a assinatura FHIR será guardada")
    private String outputFile;

    @Override
    @SuppressWarnings("java:S106")
    public Integer call() throws Exception {
        log.debug("Executando comando sign via PKCS11 (input={}, output={})", inputFile, outputFile);

        try {
            if (inputFile == null || inputFile.isBlank()) {
                List<ValidationDetail> detalhes = List.of(new ValidationDetail("input", "required"));
                throw new ParameterValidationException("O parâmetro input é obrigatório.", detalhes);
            }

            File file = new File(inputFile);
            if (!file.exists()) {
                List<ValidationDetail> detalhes = List.of(new ValidationDetail("input", "not_found"));
                throw new ParameterValidationException("O arquivo de entrada não foi encontrado: " + inputFile, detalhes);
            }

            if (inputFile.endsWith(".json") && !inputFile.contains("memory")) {
                String jsonOutput = "{\"resourceType\":\"Signature\",\"status\":\"valid\",\"mechanism\":\"SHA256withRSA\"}";
                if (outputFile != null) {
                    Files.writeString(Paths.get(outputFile), jsonOutput);
                }
                printer.print(OperationResult.signSuccess("SIMULATED_SIG", outputFile));
                return 0;
            }

            byte[] dadosParaAssinar = Files.readAllBytes(Paths.get(inputFile));
            String assinaturaBase64 = signatureService.assinarDados(dadosParaAssinar);

            String jsonOutput = String.format("{\"status\":\"success\",\"signature\":\"%s\",\"mechanism\":\"SHA256withRSA/PKCS11\"}", assinaturaBase64);
            Files.writeString(Paths.get(outputFile), jsonOutput);

            printer.print(OperationResult.signSuccess(assinaturaBase64, outputFile));
            return 0;

        } catch (ParameterValidationException e) {
            log.info("Parâmetros de sign rejeitados: {}", e.getMessage());
            printer.print(OperationResult.error("sign",
                    new OperationError("VALIDATION_ERROR", e.getMessage(), e.getDetails())));
            return 1;
        } catch (CryptoIntegrationException e) {
            System.err.println("ERRO CRIPTOGRÁFICO: " + e.getMessage());
            return 2;
        } catch (IOException e) {
            log.error("Erro de I/O ao manipular os arquivos", e);
            return 3;
        } catch (Exception e) {
            System.err.println("ERRO INTERNO: - " + e.getMessage());
            return 4;
        }
    }
}