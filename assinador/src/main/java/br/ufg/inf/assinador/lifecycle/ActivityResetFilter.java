package br.ufg.inf.assinador.lifecycle;

import jakarta.servlet.*;
import jakarta.servlet.http.HttpServletRequest;
import org.springframework.stereotype.Component;
import java.io.IOException;

@Component
public class ActivityResetFilter implements Filter {

    private final InactivityShutdownManager shutdownManager;

    public ActivityResetFilter(InactivityShutdownManager shutdownManager) {
        this.shutdownManager = shutdownManager;
    }

    @Override
    public void doFilter(ServletRequest request, ServletResponse response, FilterChain chain)
            throws IOException, ServletException {

        if (request instanceof HttpServletRequest) {
            shutdownManager.resetTimer();
        }

        chain.doFilter(request, response);
    }
}