package br.ufg.inf.assinador.cli;

import br.ufg.inf.assinador.service.SignatureService;
import org.springframework.stereotype.Component;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.concurrent.Callable;

@Component
@Command(name = "sign", description = "Gera uma assinatura digital simulada para um pacote FHIR de entrada.")
public class SignCommand implements Callable<Integer> {

    private final SignatureService signatureService;

    public SignCommand(SignatureService signatureService) {
        this.signatureService = signatureService;
    }

    @Option(names = {"-i", "--input"}, description = "Caminho para o ficheiro FHIR (.json) a ser assinado", required = true)
    private String inputFile;

    @Option(names = {"-o", "--output"}, description = "Caminho onde a assinatura FHIR será guardada", required = true)
    private String outputFile;

    @Override
    public Integer call() {
        try {
            System.out.println("A validar os parâmetros de entrada...");

            String mockSignature = signatureService.sign(inputFile);

            Files.writeString(Paths.get(outputFile), mockSignature);

            System.out.println("SUCESSO: Assinatura simulada guardada com êxito em: " + outputFile);
            return 0;

        } catch (IllegalArgumentException e) {
            System.err.println("ERRO DE VALIDAÇÃO: " + e.getMessage());
            return 1;

        } catch (Exception e) {
            System.err.println("ERRO INTERNO: Falha ao processar a assinatura - " + e.getMessage());
            return 2;
        }
    }
}
