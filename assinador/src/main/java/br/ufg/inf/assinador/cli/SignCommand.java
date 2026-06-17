package br.ufg.inf.assinador.cli;

import br.ufg.inf.assinador.dto.OperationError;
import br.ufg.inf.assinador.dto.OperationResult;
import br.ufg.inf.assinador.exception.ParameterValidationException;
import br.ufg.inf.assinador.service.SignatureService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.io.IOException;
import java.util.concurrent.Callable;

@Component
@Command(name = "sign",
        description = "Gera uma assinatura digital simulada para um pacote FHIR de entrada.",
        mixinStandardHelpOptions = true)
public class SignCommand implements Callable<Integer> {

    private static final Logger log = LoggerFactory.getLogger(SignCommand.class);

    private final SignatureService signatureService;
    private final OperationResultPrinter printer;

    public SignCommand(SignatureService signatureService, OperationResultPrinter printer) {
        this.signatureService = signatureService;
        this.printer = printer;
    }

    @Option(names = {"-i", "--input"}, description = "Caminho para o ficheiro FHIR (.json) a ser assinado")
    private String inputFile;

    @Option(names = {"-o", "--output"}, description = "Caminho onde a assinatura FHIR será guardada")
    private String outputFile;

    @Override
    public Integer call() throws IOException {
        log.debug("Executando comando sign (input={}, output={})", inputFile, outputFile);
        try {
            printer.print(signatureService.sign(inputFile, outputFile));
            return 0;
        } catch (ParameterValidationException e) {
            log.info("Parâmetros de sign rejeitados: {}", e.getMessage());
            printer.print(OperationResult.error("sign",
                    new OperationError("VALIDATION_ERROR", e.getMessage(), e.getDetails())));
            return 1;
        } catch (RuntimeException e) {
            log.error("Erro inesperado ao processar a assinatura", e);
            printer.print(OperationResult.error("sign",
                    new OperationError("INTERNAL_ERROR", e.getMessage(), null)));
            return 2;
        }
    }
}
