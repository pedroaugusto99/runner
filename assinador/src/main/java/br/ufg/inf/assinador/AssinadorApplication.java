package br.ufg.inf.assinador;

import br.ufg.inf.assinador.cli.AssinadorCommand;
import org.springframework.boot.CommandLineRunner;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import picocli.CommandLine;
import picocli.CommandLine.IFactory;

@SpringBootApplication
public class AssinadorApplication implements CommandLineRunner {

    private final AssinadorCommand assinadorCommand;
    private final IFactory factory;
    private int exitCode = 0;

    public AssinadorApplication(AssinadorCommand assinadorCommand, IFactory factory) {
        this.assinadorCommand = assinadorCommand;
        this.factory = factory;
    }

    public static void main(String[] args) {
        System.exit(SpringApplication.exit(SpringApplication.run(AssinadorApplication.class, args)));
    }

    @Override
    public void run(String... args) {
        this.exitCode = new CommandLine(assinadorCommand, factory).execute(args);
    }

    @org.springframework.context.annotation.Bean
    public org.springframework.boot.ExitCodeGenerator exitCodeGenerator() {
        return () -> exitCode;
    }
}