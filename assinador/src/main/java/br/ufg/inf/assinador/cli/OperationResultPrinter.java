package br.ufg.inf.assinador.cli;

import br.ufg.inf.assinador.dto.OperationResult;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.stereotype.Component;

import java.io.IOException;

@Component
public class OperationResultPrinter {

    private final ObjectMapper objectMapper;

    public OperationResultPrinter(ObjectMapper objectMapper) {
        this.objectMapper = objectMapper;
    }

    public void print(OperationResult result) throws IOException {
        byte[] json = objectMapper.writeValueAsBytes(result);
        System.out.write(json);
        System.out.write('\n');
        System.out.flush();
    }
}
