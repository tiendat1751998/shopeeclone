package com.tiki.auth.service.session;

import com.tiki.auth.entity.RefreshToken;
import com.tiki.auth.repository.RefreshTokenRepository;
import com.tiki.auth.security.JwtTokenProvider;
import io.jsonwebtoken.Claims;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class SessionService {

    private static final Logger log = LoggerFactory.getLogger(SessionService.class);

    private final RefreshTokenRepository refreshTokenRepository;
    private final JwtTokenProvider jwtTokenProvider;

    @Transactional
    public RefreshToken createRefreshToken(UUID userId, String token) {
        RefreshToken refreshToken = RefreshToken.builder()
            .token(token)
            .userId(userId)
            .expiresAt(LocalDateTime.now().plusSeconds(jwtTokenProvider.getRefreshTtlSeconds()))
            .revoked(false)
            .build();
        return refreshTokenRepository.save(refreshToken);
    }

    @Transactional
    public void revokeToken(String refreshTokenValue) {
        if (refreshTokenValue == null || refreshTokenValue.isBlank()) {
            return;
        }
        refreshTokenRepository.findByToken(refreshTokenValue)
            .ifPresent(token -> {
                token.setRevoked(true);
                refreshTokenRepository.save(token);
            });
    }

    @Transactional
    public void revokeAllForUser(UUID userId) {
        refreshTokenRepository.deleteByUserId(userId);
    }

    @Transactional
    public Claims validateAndRotate(String refreshTokenValue) {
        Claims claims = jwtTokenProvider.validateRefreshToken(refreshTokenValue);
        if (claims == null) {
            throw new BadCredentialsException("Invalid or expired refresh token");
        }

        RefreshToken storedToken = refreshTokenRepository.findByToken(refreshTokenValue)
            .orElseThrow(() -> new BadCredentialsException("Refresh token not found"));

        if (storedToken.getRevoked()) {
            log.warn("Attempted to use revoked refresh token for user: {}", storedToken.getUserId());
            refreshTokenRepository.deleteByUserId(storedToken.getUserId());
            throw new BadCredentialsException("Refresh token has been revoked");
        }

        if (storedToken.getExpiresAt().isBefore(LocalDateTime.now())) {
            throw new BadCredentialsException("Refresh token has expired");
        }

        storedToken.setRevoked(true);
        refreshTokenRepository.save(storedToken);

        return claims;
    }
}
