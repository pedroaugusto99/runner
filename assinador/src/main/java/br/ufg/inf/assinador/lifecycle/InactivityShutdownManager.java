package br.ufg.inf.assinador.lifecycle;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.boot.SpringApplication;
import org.springframework.context.ApplicationContext;
import org.springframework.stereotype.Component;

import jakarta.annotation.PostConstruct;
import java.util.Timer;
import java.util.TimerTask;

@Component
public class InactivityShutdownManager {

    private static final Logger log = LoggerFactory.getLogger(InactivityShutdownManager.class);

    private final ApplicationContext appContext;
    private final long timeoutMs;
    private Timer timer;

    public InactivityShutdownManager(ApplicationContext appContext,
                                     @Value("${assinador.inactivity-timeout-ms:300000}") long timeoutMs) {
        this.appContext = appContext;
        this.timeoutMs = timeoutMs;
    }

    @PostConstruct
    public void start() {
        log.info("Modo servidor ativo. Auto-shutdown programado para uma janela de inatividade de {} ms", timeoutMs);
        resetTimer();
    }

    public synchronized void resetTimer() {
        if (timer != null) {
            timer.cancel();
        }
        timer = new Timer(true);
        timer.schedule(new TimerTask() {
            @Override
            public void run() {
                log.info("Inatividade prolongada detetada ({} ms sem receber pedidos). A encerrar o servidor graciosamente...", timeoutMs);
                SpringApplication.exit(appContext, () -> 0);
            }
        }, timeoutMs);
        log.debug("Timer de inatividade reiniciado.");
    }
}