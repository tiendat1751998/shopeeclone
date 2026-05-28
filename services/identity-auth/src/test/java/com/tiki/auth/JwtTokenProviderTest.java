package com.tiki.auth;

import com.tiki.auth.entity.User;
import com.tiki.auth.security.JwtTokenProvider;
import com.tiki.auth.security.JwksProvider;
import io.jsonwebtoken.Claims;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;

class JwtTokenProviderTest {

    private JwksProvider jwksProvider;
    private JwtTokenProvider jwtTokenProvider;

    @BeforeEach
    void setUp() {
        jwksProvider = new JwksProvider();
        jwtTokenProvider = new JwtTokenProvider(jwksProvider);

        try {
            var jwksField = JwksProvider.class.getDeclaredField("configuredPrivateKey");
            jwksField.setAccessible(true);
            jwksField.set(jwksProvider, "test");
            var jwksField2 = JwksProvider.class.getDeclaredField("configuredPublicKey");
            jwksField2.setAccessible(true);
            jwksField2.set(jwksProvider, "test");
            jwksProvider.init();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }

        try {
            var field = JwtTokenProvider.class.getDeclaredField("accessSecret");
            field.setAccessible(true);
            field.set(jwtTokenProvider, "test-access-secret-key-must-be-32-bytes-long!!");
            var field2 = JwtTokenProvider.class.getDeclaredField("refreshSecret");
            field2.setAccessible(true);
            field2.set(jwtTokenProvider, "test-refresh-secret-key-must-be-32-bytes-long!!");
            jwtTokenProvider.init();
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    @Test
    void generateAccessToken_ContainsExpectedClaims() {
        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("test@example.com")
            .fullName("Test User")
            .role("BUYER")
            .build();

        String token = jwtTokenProvider.generateAccessToken(user);
        assertThat(token).isNotBlank();

        Claims claims = jwtTokenProvider.validateAccessToken(token);
        assertThat(claims).isNotNull();
        assertThat(claims.getSubject()).isEqualTo(user.getUserId().toString());
        assertThat(claims.get("email", String.class)).isEqualTo("test@example.com");
        assertThat(claims.get("role", String.class)).isEqualTo("BUYER");
        assertThat(claims.get("scope", String.class)).isNotBlank();
    }

    @Test
    void generateRefreshToken_ContainsExpectedClaims() {
        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("test@example.com")
            .fullName("Test User")
            .role("BUYER")
            .build();

        String token = jwtTokenProvider.generateRefreshToken(user);
        assertThat(token).isNotBlank();

        Claims claims = jwtTokenProvider.validateRefreshToken(token);
        assertThat(claims).isNotNull();
        assertThat(claims.getSubject()).isEqualTo(user.getUserId().toString());
        assertThat(claims.get("type", String.class)).isEqualTo("refresh");
    }

    @Test
    void validateAccessToken_InvalidToken_ReturnsNull() {
        Claims claims = jwtTokenProvider.validateAccessToken("invalid-token");
        assertThat(claims).isNull();
    }

    @Test
    void validateRefreshToken_InvalidToken_ReturnsNull() {
        Claims claims = jwtTokenProvider.validateRefreshToken("invalid-token");
        assertThat(claims).isNull();
    }

    @Test
    void validateRefreshToken_AccessTokenUsedAsRefresh_ReturnsNull() {
        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("test@example.com")
            .fullName("Test User")
            .role("BUYER")
            .build();

        String accessToken = jwtTokenProvider.generateAccessToken(user);
        Claims claims = jwtTokenProvider.validateRefreshToken(accessToken);
        assertThat(claims).isNull();
    }

    @Test
    void getUserIdFromToken_ValidToken_ReturnsUserId() {
        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("test@example.com")
            .fullName("Test User")
            .role("BUYER")
            .build();

        String token = jwtTokenProvider.generateAccessToken(user);
        String userId = jwtTokenProvider.getUserIdFromToken(token);
        assertThat(userId).isEqualTo(user.getUserId().toString());
    }

    @Test
    void getUserIdFromToken_InvalidToken_ReturnsNull() {
        String userId = jwtTokenProvider.getUserIdFromToken("invalid");
        assertThat(userId).isNull();
    }

    @Test
    void adminRole_ScopeContainsAdminAccess() {
        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("admin@example.com")
            .fullName("Admin")
            .role("ADMIN")
            .build();

        String token = jwtTokenProvider.generateAccessToken(user);
        Claims claims = jwtTokenProvider.validateAccessToken(token);
        String scope = claims.get("scope", String.class);
        assertThat(scope).contains("admin:access");
    }

    @Test
    void sellerRole_ScopeContainsProductAndInventory() {
        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("seller@example.com")
            .fullName("Seller")
            .role("SELLER")
            .build();

        String token = jwtTokenProvider.generateAccessToken(user);
        Claims claims = jwtTokenProvider.validateAccessToken(token);
        String scope = claims.get("scope", String.class);
        assertThat(scope).contains("products:read");
        assertThat(scope).contains("inventory:read");
    }

    @Test
    void buyerRole_ScopeContainsProductAndOrderRead() {
        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("buyer@example.com")
            .fullName("Buyer")
            .role("BUYER")
            .build();

        String token = jwtTokenProvider.generateAccessToken(user);
        Claims claims = jwtTokenProvider.validateAccessToken(token);
        String scope = claims.get("scope", String.class);
        assertThat(scope).contains("products:read");
        assertThat(scope).contains("orders:read");
    }
}
