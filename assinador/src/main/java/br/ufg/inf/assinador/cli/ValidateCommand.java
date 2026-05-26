package br.ufg.inf.assinador.cli;

import br.ufg.inf.assinador.service.SignatureService;
import org.springframework.stereotype.Component;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.util.concurrent.Callable;

@Component
@Command(name = "validate", description = "Simula a validação de um pacote de assinatura FHIR.")
public class ValidateCommand implements Callable<Integer> {

    private final SignatureService signatureService;

    public ValidateCommand(SignatureService signatureService) {
        this.signatureService = signatureService;
    }

    @Option(names = {"-s", "--signature"}, description = "Caminho para o ficheiro de assinatura (.json)", required = true)
    private String signatureFile;

    @Override
    public Integer call() {
        try {
            System.out.println("A analisar o ficheiro de assinatura...");

            boolean isValid = signatureService.validate(signatureFile);

            if (isValid) {
                System.out.println("RESULTADO: [VÁLIDA] A assinatura foi verificada com sucesso.");
                return 0;
            } else {
                System.out.println("RESULTADO: [INVÁLIDA] A assinatura não confere ou o formato é incorreto.");
                return 1;
            }

        } catch (IllegalArgumentException e) {
            System.err.println("ERRO DE VALIDAÇÃO: " + e.getMessage());
            return 2;
        } catch (Exception e) {
            System.err.println("ERRO INTERNO: Falha ao validar a assinatura - " + e.getMessage());
            return 3;
        }
    }
}