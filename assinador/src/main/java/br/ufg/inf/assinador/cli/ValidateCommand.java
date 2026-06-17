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
@Command(name = "validate",
        description = "Simula a validação de um pacote de assinatura FHIR.",
        mixinStandardHelpOptions = true)
public class ValidateCommand implements Callable<Integer> {

    private static final Logger log = LoggerFactory.getLogger(ValidateCommand.class);

    private final SignatureService signatureService;
    private final OperationResultPrinter printer;

    public ValidateCommand(SignatureService signatureService, OperationResultPrinter printer) {
        this.signatureService = signatureService;
        this.printer = printer;
    }

    @Option(names = {"-s", "--signature"}, description = "Caminho para o ficheiro de assinatura (.json)")
    private String signatureFile;

    @Override
    public Integer call() throws IOException {
        log.debug("Executando comando validate (signature={})", signatureFile);
        try {
            OperationResult result = signatureService.validate(signatureFile);
            printer.print(result);
            return Boolean.TRUE.equals(result.valid()) ? 0 : 1;
        } catch (ParameterValidationException e) {
            log.info("Parâmetros de validate rejeitados: {}", e.getMessage());
            printer.print(OperationResult.error("validate",
                    new OperationError("VALIDATION_ERROR", e.getMessage(), e.getDetails())));
            return 2;
        } catch (RuntimeException e) {
            log.error("Erro inesperado ao validar a assinatura", e);
            printer.print(OperationResult.error("validate",
                    new OperationError("INTERNAL_ERROR", e.getMessage(), null)));
            return 3;
        }
    }
}
