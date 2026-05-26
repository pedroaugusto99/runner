package br.ufg.inf.assinador.service;

import org.springframework.stereotype.Service;
import java.io.File;

@Service
public class FakeSignatureService implements SignatureService {

    @Override
    public String sign(String inputFilePath) {
        File file = new File(inputFilePath);
        if (!file.exists() || !file.isFile()) {
            throw new IllegalArgumentException("O ficheiro de entrada não foi encontrado ou é inválido: " + inputFilePath);
        }

        return """
               {
                 "resourceType": "Signature",
                 "status": "simulated",
                 "type": [
                   {
                     "system": "urn:iso-astm:E1762-95:2013",
                     "code": "1.2.840.10065.1.12.1.1"
                   }
                 ],
                 "data": "YXNzaW5hdHVyYSBzaW11bGFkYSBkZSBleGVtcGxv"
               }
               """;
    }

    @Override
    public boolean validate(String signatureFilePath) {
        File file = new File(signatureFilePath);
        if (!file.exists() || !file.isFile()) {
            throw new IllegalArgumentException("O ficheiro de assinatura não foi encontrado: " + signatureFilePath);
        }

        return signatureFilePath.toLowerCase().endsWith(".json");
    }
}