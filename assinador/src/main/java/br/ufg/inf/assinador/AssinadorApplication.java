package br.ufg.inf.assinador;

import br.ufg.inf.assinador.cli.AssinadorCommand;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import picocli.CommandLine;
import picocli.CommandLine.IFactory;

@SpringBootApplication
public class AssinadorApplication implements CommandLineRunner {

    private static final Logger log = LoggerFactory.getLogger(AssinadorApplication.class);

    private final AssinadorCommand assinadorCommand;
    private final IFactory factory;
    private int exitCode = 0;

    public AssinadorApplication(AssinadorCommand assinadorCommand, IFactory factory) {
        this.assinadorCommand = assinadorCommand;
        this.factory = factory;
    }

    public static void main(String[] args) {
        if (args.length == 0) {
            SpringApplication.run(AssinadorApplication.class, args);
        } else {
            System.exit(SpringApplication.exit(SpringApplication.run(AssinadorApplication.class, args)));
        }
    }

    @Override
    public void run(String... args) {
        if (args.length == 0) {
            log.info("Nenhum argumento de CLI fornecido. Iniciando no Modo Servidor (HTTP)...");
        } else {
            log.info("Argumentos detetados. A iniciar processamento no Modo CLI (Local)...");
            this.exitCode = new CommandLine(assinadorCommand, factory).execute(args);
        }
    }

    @org.springframework.context.annotation.Bean
    public org.springframework.boot.ExitCodeGenerator exitCodeGenerator() {
        return () -> exitCode;
    }
}