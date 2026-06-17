package br.ufg.inf.assinador.service;

import br.ufg.inf.assinador.dto.OperationResult;

public interface SignatureService {

    OperationResult sign(String inputPath, String outputPath);

    OperationResult validate(String signaturePath);
}
