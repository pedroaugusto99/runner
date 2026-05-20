package br.ufg.inf.assinador.service;

public interface SignatureService {
    String sign(String inputFilePath);
    boolean validate(String signatureFilePath);
}
