package br.ufg.inf.assinador.cli;

import br.ufg.inf.assinador.dto.OperationError;
import br.ufg.inf.assinador.dto.OperationResult;
import br.ufg.inf.assinador.dto.ValidationDetail;
import br.ufg.inf.assinador.exception.ParameterValidationException;
import br.ufg.inf.assinador.service.SignatureService;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Component;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.io.File;
import java.io.IOException;
import java.util.Collections;
import java.util.List;
import java.util.concurrent.Callable;

@Component
@Command(name = "validate",
        description = "Valida um pacote de assinatura FHIR gerado via PKCS#11.",
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
    @SuppressWarnings("java:S106")
    public Integer call() throws Exception {
        log.debug("Executando comando validate (signature={})", signatureFile);

        try {
            System.out.println("A analisar o ficheiro de assinatura...");
            File file = new File(signatureFile);
            if (!file.exists()) {
                List<ValidationDetail> detalhes = List.of(new ValidationDetail("signature", "not_found"));
                throw new ParameterValidationException("O ficheiro de assinatura não foi encontrado: " + signatureFile, detalhes);
            }

            OperationResult result = signatureService.validate(signatureFile);
            printer.print(result);

            System.out.println("RESULTADO: [VÁLIDA] A assinatura PKCS#11 foi verificada com sucesso.");
            return Boolean.TRUE.equals(result.valid()) ? 0 : 1;

        } catch (ParameterValidationException e) {
            log.info("Parâmetros de validate rejeitados: {}", e.getMessage());
            printer.print(OperationResult.error("validate",
                    new OperationError("VALIDATION_ERROR", e.getMessage(), e.getDetails())));
            return 2;
        } catch (IOException e) {
            log.error("Erro de I/O ao processar ou imprimir o resultado da validacao", e);
            OperationError errorDto = new OperationError("IO_ERROR", "Falha operacional na leitura ou impressao: " + e.getMessage(), Collections.emptyList());
            printer.print(OperationResult.error("validate", errorDto));
            return 3;
        } catch (RuntimeException e) {
            log.error("Erro inesperado ao validar a assinatura", e);
            printer.print(OperationResult.error("validate",
                    new OperationError("INTERNAL_ERROR", e.getMessage(), null)));
            return 4;
        }
    }
}
