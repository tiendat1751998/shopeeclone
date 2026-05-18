package com.shopee.auth.security;

import com.nimbusds.jose.jwk.RSAKey;
import jakarta.annotation.PostConstruct;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import java.security.KeyFactory;
import java.security.KeyPair;
import java.security.KeyPairGenerator;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.security.interfaces.RSAPrivateKey;
import java.security.interfaces.RSAPublicKey;
import java.security.spec.PKCS8EncodedKeySpec;
import java.security.spec.X509EncodedKeySpec;
import java.util.Base64;
import java.util.UUID;

@Component
public class JwksProvider {

    private static final Logger log = LoggerFactory.getLogger(JwksProvider.class);

    @Value("${jwt.rsa.private-key:}")
    private String configuredPrivateKey;

    @Value("${jwt.rsa.public-key:}")
    private String configuredPublicKey;

    private RSAKey rsaKey;
    private KeyPair keyPair;
    private String keyId;

    @PostConstruct
    public void init() {
        this.keyId = UUID.randomUUID().toString();

        try {
            if (!configuredPrivateKey.isBlank() && !configuredPublicKey.isBlank()) {
                keyPair = loadKeyPair(configuredPrivateKey, configuredPublicKey);
                log.info("Loaded configured RSA key pair");
            } else {
                keyPair = generateKeyPair();
                log.info("Generated ephemeral RSA key pair (not persisted across restarts)");
            }

            this.rsaKey = new RSAKey.Builder((RSAPublicKey) keyPair.getPublic())
                .privateKey((RSAPrivateKey) keyPair.getPrivate())
                .keyID(keyId)
                .build();
        } catch (Exception e) {
            throw new RuntimeException("Failed to initialize JWKS provider", e);
        }
    }

    public RSAKey getRsaKey() {
        return rsaKey;
    }

    public PublicKey getPublicKey() {
        return keyPair.getPublic();
    }

    public PrivateKey getPrivateKey() {
        return keyPair.getPrivate();
    }

    public String getKeyId() {
        return keyId;
    }

    public String getJwksJson() {
        try {
            return rsaKey.toPublicJWK().toJSONString();
        } catch (Exception e) {
            log.error("Failed to serialize JWKS", e);
            return "{}";
        }
    }

    public String getJwksSetJson() {
        try {
            return "{\"keys\":[" + rsaKey.toPublicJWK().toJSONString() + "]}";
        } catch (Exception e) {
            log.error("Failed to serialize JWKS set", e);
            return "{\"keys\":[]}";
        }
    }

    private KeyPair generateKeyPair() throws Exception {
        KeyPairGenerator gen = KeyPairGenerator.getInstance("RSA");
        gen.initialize(2048);
        return gen.generateKeyPair();
    }

    private KeyPair loadKeyPair(String privateKeyPem, String publicKeyPem) throws Exception {
        KeyFactory keyFactory = KeyFactory.getInstance("RSA");

        String privateKeyContent = privateKeyPem
            .replace("-----BEGIN PRIVATE KEY-----", "")
            .replace("-----END PRIVATE KEY-----", "")
            .replaceAll("\\s", "");

        String publicKeyContent = publicKeyPem
            .replace("-----BEGIN PUBLIC KEY-----", "")
            .replace("-----END PUBLIC KEY-----", "")
            .replaceAll("\\s", "");

        PKCS8EncodedKeySpec privSpec = new PKCS8EncodedKeySpec(Base64.getDecoder().decode(privateKeyContent));
        PrivateKey privateKey = keyFactory.generatePrivate(privSpec);

        X509EncodedKeySpec pubSpec = new X509EncodedKeySpec(Base64.getDecoder().decode(publicKeyContent));
        PublicKey publicKey = keyFactory.generatePublic(pubSpec);

        return new KeyPair(publicKey, privateKey);
    }
}
