package com.shopee.auth.security;

import com.shopee.auth.entity.User;
import io.jsonwebtoken.*;
import io.jsonwebtoken.security.Keys;
import jakarta.annotation.PostConstruct;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.SecretKey;
import java.nio.charset.StandardCharsets;
import java.security.PrivateKey;
import java.security.PublicKey;
import java.util.*;

@Component
public class JwtTokenProvider {

    private static final Logger log = LoggerFactory.getLogger(JwtTokenProvider.class);

    private static final String ISSUER = "shopee-clone";

    @Value("${jwt.access-secret:}")
    private String accessSecret;

    @Value("${jwt.refresh-secret:}")
    private String refreshSecret;

    @Value("${jwt.access-ttl-seconds:900}")
    private long accessTtlSeconds;

    @Value("${jwt.refresh-ttl-seconds:604800}")
    private long refreshTtlSeconds;

    private final JwksProvider jwksProvider;

    private SecretKey hmacAccessKey;
    private SecretKey hmacRefreshKey;
    private boolean useRsa;

    public JwtTokenProvider(JwksProvider jwksProvider) {
        this.jwksProvider = jwksProvider;
    }

    @PostConstruct
    public void init() {
        // Always initialize HMAC keys as fallback
        if (accessSecret != null && !accessSecret.isBlank()) {
            this.hmacAccessKey = Keys.hmacShaKeyFor(accessSecret.getBytes(StandardCharsets.UTF_8));
        }
        if (refreshSecret != null && !refreshSecret.isBlank()) {
            this.hmacRefreshKey = Keys.hmacShaKeyFor(refreshSecret.getBytes(StandardCharsets.UTF_8));
        }
        // Use RSA only when JWKS provider has valid keys loaded
        this.useRsa = jwksProvider.getPublicKey() != null && jwksProvider.getPrivateKey() != null;
        // Fail fast if no signing key is configured
        if (!this.useRsa && this.hmacAccessKey == null) {
            throw new IllegalStateException(
                "No signing key configured: provide either JWKS RSA keys or HMAC secrets"
            );
        }
        log.info("JWT signing mode: {}", useRsa ? "RSA (JWKS)" : "HMAC");
    }

    public String generateAccessToken(User user) {
        Date now = new Date();
        Date expiry = new Date(now.getTime() + accessTtlSeconds * 1000);

        JwtBuilder builder = Jwts.builder()
            .issuer(ISSUER)
            .subject(user.getUserId().toString())
            .claim("email", user.getEmail())
            .claim("role", user.getRole())
            .claim("scope", buildScope(user))
            .issuedAt(now)
            .expiration(expiry);

        if (useRsa) {
            builder = builder
                .header().keyId(jwksProvider.getKeyId()).and()
                .signWith(jwksProvider.getPrivateKey(), SignatureAlgorithm.RS256);
        } else {
            builder = builder.signWith(getEffectiveAccessKey());
        }

        return builder.compact();
    }

    public String generateRefreshToken(User user) {
        Date now = new Date();
        Date expiry = new Date(now.getTime() + refreshTtlSeconds * 1000);

        JwtBuilder builder = Jwts.builder()
            .issuer(ISSUER)
            .subject(user.getUserId().toString())
            .claim("type", "refresh")
            .issuedAt(now)
            .expiration(expiry);

        if (useRsa) {
            builder = builder.signWith(jwksProvider.getPrivateKey(), SignatureAlgorithm.RS256);
        } else {
            builder = builder.signWith(getEffectiveRefreshKey());
        }

        return builder.compact();
    }

    public Claims validateAccessToken(String token) {
        try {
            JwtParserBuilder parserBuilder = Jwts.parser();
            if (useRsa) {
                parserBuilder.verifyWith(jwksProvider.getPublicKey());
            } else {
                parserBuilder.verifyWith(getEffectiveAccessKey());
            }
            return parserBuilder.build().parseSignedClaims(token).getPayload();
        } catch (JwtException | IllegalArgumentException e) {
            log.warn("Invalid access token: {}", e.getMessage());
            return null;
        }
    }

    public Claims validateRefreshToken(String token) {
        try {
            JwtParserBuilder parserBuilder = Jwts.parser();
            if (useRsa) {
                parserBuilder.verifyWith(jwksProvider.getPublicKey());
            } else {
                parserBuilder.verifyWith(getEffectiveRefreshKey());
            }
            Claims claims = parserBuilder.build().parseSignedClaims(token).getPayload();

            if (!"refresh".equals(claims.get("type"))) {
                return null;
            }
            return claims;
        } catch (JwtException | IllegalArgumentException e) {
            log.warn("Invalid refresh token: {}", e.getMessage());
            return null;
        }
    }

    public String getUserIdFromToken(String token) {
        Claims claims = validateAccessToken(token);
        return claims != null ? claims.getSubject() : null;
    }

    public long getAccessTtlSeconds() {
        return accessTtlSeconds;
    }

    public long getRefreshTtlSeconds() {
        return refreshTtlSeconds;
    }

    private SecretKey getEffectiveAccessKey() {
        if (useRsa) {
            // This should not be called when using RSA, but return HMAC key as fallback
            return hmacAccessKey;
        }
        return hmacAccessKey;
    }

    private SecretKey getEffectiveRefreshKey() {
        if (useRsa) {
            return hmacRefreshKey;
        }
        return hmacRefreshKey;
    }

    private String buildScope(User user) {
        StringBuilder scope = new StringBuilder();
        if ("SELLER".equals(user.getRole())) {
            scope.append("products:read products:write orders:read inventory:read inventory:write");
        } else if ("ADMIN".equals(user.getRole()) || "SUPER_ADMIN".equals(user.getRole())) {
            scope.append("admin:access users:read users:write products:read products:write orders:read payments:read");
        } else {
            scope.append("products:read orders:read orders:write");
        }
        return scope.toString().trim();
    }
}
