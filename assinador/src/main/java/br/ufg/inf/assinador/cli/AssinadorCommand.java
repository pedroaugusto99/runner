package br.ufg.inf.assinador.cli;

import org.springframework.stereotype.Component;
import picocli.CommandLine.Command;
import java.util.concurrent.Callable;

@Component
@Command(name = "assinador",
        mixinStandardHelpOptions = true,
        version = "0.1.0",
        description = "Motor de Assinatura FHIR para HubSaude",
        subcommands = { SignCommand.class })
public class AssinadorCommand implements Callable<Integer> {

    @Override
    public Integer call() {
        System.out.println("Por favor, especifique um comando (ex: sign). Use --help para ver as opções disponíveis.");
        return 0;
    }
}