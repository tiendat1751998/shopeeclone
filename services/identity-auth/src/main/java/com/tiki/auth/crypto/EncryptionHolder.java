package com.tiki.auth.crypto;

import jakarta.annotation.PostConstruct;
import lombok.Getter;
import org.springframework.stereotype.Component;

@Component
public class EncryptionHolder {

    @Getter
    private static EncryptionUtil instance;

    private final EncryptionUtil encryptionUtil;

    public EncryptionHolder(EncryptionUtil encryptionUtil) {
        this.encryptionUtil = encryptionUtil;
    }

    @PostConstruct
    public void init() {
        instance = encryptionUtil;
    }
}
