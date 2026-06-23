package br.ufg.inf.assinador.exception;

public class CryptoIntegrationException extends RuntimeException {
    public CryptoIntegrationException(String message) {
        super(message);
    }

    public CryptoIntegrationException(String message, Throwable cause) {
        super(message, cause);
    }
}