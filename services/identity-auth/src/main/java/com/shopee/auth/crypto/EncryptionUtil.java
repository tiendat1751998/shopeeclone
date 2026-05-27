package com.shopee.auth.crypto;

import jakarta.annotation.PostConstruct;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.Cipher;
import javax.crypto.spec.GCMParameterSpec;
import javax.crypto.spec.SecretKeySpec;
import java.security.SecureRandom;
import java.util.Base64;

@Component
public class EncryptionUtil {

    private static final Logger log = LoggerFactory.getLogger(EncryptionUtil.class);
    private static final String ALGORITHM = "AES/GCM/NoPadding";
    private static final int GCM_TAG_LENGTH = 128;
    private static final int GCM_IV_LENGTH = 12;

    @Value("${encryption.key:}")
    private String encryptionKey;

    private SecretKeySpec keySpec;
    private boolean enabled;

    @PostConstruct
    public void init() {
        if (encryptionKey == null || encryptionKey.isBlank()) {
            log.warn("Encryption key not configured - PII data will not be encrypted at rest");
            enabled = false;
            return;
        }
        byte[] keyBytes = hexStringToByteArray(encryptionKey);
        if (keyBytes == null) {
            keyBytes = encryptionKey.getBytes(java.nio.charset.StandardCharsets.UTF_8);
        }
        if (keyBytes.length != 16 && keyBytes.length != 24 && keyBytes.length != 32) {
            log.warn("Encryption key must be 16, 24, or 32 bytes - PII will not be encrypted");
            enabled = false;
            return;
        }
        this.keySpec = new SecretKeySpec(keyBytes, "AES");
        this.enabled = true;
        log.info("PII encryption initialized (AES-{} GCM)", keyBytes.length * 8);
    }

    private static byte[] hexStringToByteArray(String hex) {
        int len = hex.length();
        if (len % 2 != 0) return null;
        try {
            byte[] data = new byte[len / 2];
            for (int i = 0; i < len; i += 2) {
                data[i / 2] = (byte) ((Character.digit(hex.charAt(i), 16) << 4)
                    + Character.digit(hex.charAt(i + 1), 16));
            }
            return data;
        } catch (Exception e) {
            return null;
        }
    }

    public String encrypt(String plaintext) {
        if (!enabled || plaintext == null || plaintext.isEmpty()) {
            return plaintext;
        }
        try {
            Cipher cipher = Cipher.getInstance(ALGORITHM);
            byte[] iv = new byte[GCM_IV_LENGTH];
            SecureRandom.getInstanceStrong().nextBytes(iv);
            GCMParameterSpec spec = new GCMParameterSpec(GCM_TAG_LENGTH, iv);
            cipher.init(Cipher.ENCRYPT_MODE, keySpec, spec);
            byte[] ciphertext = cipher.doFinal(plaintext.getBytes(java.nio.charset.StandardCharsets.UTF_8));
            byte[] combined = new byte[GCM_IV_LENGTH + ciphertext.length];
            System.arraycopy(iv, 0, combined, 0, GCM_IV_LENGTH);
            System.arraycopy(ciphertext, 0, combined, GCM_IV_LENGTH, ciphertext.length);
            return Base64.getEncoder().encodeToString(combined);
        } catch (Exception e) {
            log.error("Encryption failed", e);
            throw new RuntimeException("Encryption failed", e);
        }
    }

    public String decrypt(String encrypted) {
        if (!enabled || encrypted == null || encrypted.isEmpty()) {
            return encrypted;
        }
        try {
            byte[] combined = Base64.getDecoder().decode(encrypted);
            Cipher cipher = Cipher.getInstance(ALGORITHM);
            GCMParameterSpec spec = new GCMParameterSpec(GCM_TAG_LENGTH, combined, 0, GCM_IV_LENGTH);
            cipher.init(Cipher.DECRYPT_MODE, keySpec, spec);
            byte[] plaintext = cipher.doFinal(combined, GCM_IV_LENGTH, combined.length - GCM_IV_LENGTH);
            return new String(plaintext, java.nio.charset.StandardCharsets.UTF_8);
        } catch (Exception e) {
            log.error("Decryption failed", e);
            return encrypted;
        }
    }

    public boolean isEnabled() {
        return enabled;
    }
}
