package com.shopee.auth.service;

import com.shopee.auth.dto.AuthRequest;
import com.shopee.auth.dto.AuthResponse;
import com.shopee.auth.dto.UserResponse;
import com.shopee.auth.entity.RefreshToken;
import com.shopee.auth.entity.User;
import com.shopee.auth.exception.DuplicateResourceException;
import com.shopee.auth.metrics.AuthMetrics;
import com.shopee.auth.repository.RefreshTokenRepository;
import com.shopee.auth.repository.UserRepository;
import com.shopee.auth.security.AccountLockoutService;
import com.shopee.auth.security.JwtTokenProvider;
import com.shopee.auth.security.RateLimiterService;
import io.jsonwebtoken.Claims;
import jakarta.servlet.http.HttpServletRequest;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.slf4j.MDC;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.time.LocalDateTime;
import java.util.Map;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class AuthService {

    private static final Logger log = LoggerFactory.getLogger(AuthService.class);

    private final UserRepository userRepository;
    private final RefreshTokenRepository refreshTokenRepository;
    private final PasswordEncoder passwordEncoder;
    private final JwtTokenProvider jwtTokenProvider;
    private final RateLimiterService rateLimiterService;
    private final AccountLockoutService accountLockoutService;
    private final OutboxPublisher outboxPublisher;
    private final AuthMetrics authMetrics;
    private final HttpServletRequest httpServletRequest;

    @Transactional
    public AuthResponse register(AuthRequest.Register request) {
        String email = request.getEmail().toLowerCase().trim();

        if (userRepository.existsByEmail(email)) {
            authMetrics.incrementRegistrationFailures("duplicate_email");
            throw new DuplicateResourceException("Email already registered");
        }

        if (request.getPhone() != null && !request.getPhone().isBlank()
            && userRepository.existsByPhone(request.getPhone().trim())) {
            authMetrics.incrementRegistrationFailures("duplicate_phone");
            throw new DuplicateResourceException("Phone already registered");
        }

        validatePasswordStrength(request.getPassword());

        User user = User.builder()
            .email(email)
            .phone(request.getPhone() != null ? request.getPhone().trim() : null)
            .passwordHash(passwordEncoder.encode(request.getPassword()))
            .fullName(request.getFullName().trim())
            .role("BUYER")
            .isVerified(false)
            .isActive(true)
            .build();

        user = userRepository.save(user);

        String accessToken = jwtTokenProvider.generateAccessToken(user);
        String refreshToken = jwtTokenProvider.generateRefreshToken(user);

        saveRefreshToken(user.getUserId(), refreshToken);

        outboxPublisher.publish("user", user.getUserId().toString(), "user.registered", Map.of(
            "user_id", user.getUserId().toString(),
            "email", user.getEmail(),
            "role", user.getRole()
        ));

        authMetrics.incrementRegistrations();
        log.info("User registered successfully: {}", email);

        return buildAuthResponse(user, accessToken, refreshToken);
    }

    @Transactional
    public AuthResponse login(AuthRequest.Login request) {
        String email = request.getEmail().toLowerCase().trim();
        String ipAddress = resolveClientIp();

        if (accountLockoutService.isAccountLocked(email)) {
            authMetrics.incrementLoginFailures("account_locked");
            log.warn("Login blocked - account locked for email: {} from IP: {}", email, ipAddress);
            throw new BadCredentialsException("Account temporarily locked due to too many failed attempts");
        }

        if (!rateLimiterService.isLoginAllowed(email)) {
            authMetrics.incrementLoginFailures("rate_limited");
            log.warn("Login blocked - rate limit exceeded for email: {} from IP: {}", email, ipAddress);
            throw new BadCredentialsException("Too many login attempts. Please try again later.");
        }

        if (!rateLimiterService.isIpAllowed(ipAddress)) {
            authMetrics.incrementLoginFailures("ip_blocked");
            log.warn("Login blocked - IP rate limit exceeded: {}", ipAddress);
            throw new BadCredentialsException("Too many requests from this IP. Please try again later.");
        }

        User user = userRepository.findByEmail(email)
            .orElseThrow(() -> {
                rateLimiterService.recordLoginAttempt(email);
                rateLimiterService.recordIpAttempt(ipAddress);
                accountLockoutService.recordFailedAttempt(email, ipAddress);
                authMetrics.incrementLoginFailures("invalid_email");
                return new BadCredentialsException("Invalid email or password");
            });

        if (!passwordEncoder.matches(request.getPassword(), user.getPasswordHash())) {
            rateLimiterService.recordLoginAttempt(email);
            rateLimiterService.recordIpAttempt(ipAddress);
            accountLockoutService.recordFailedAttempt(email, ipAddress);
            authMetrics.incrementLoginFailures("invalid_password");
            log.warn("Failed login attempt for email: {} from IP: {}", email, ipAddress);
            throw new BadCredentialsException("Invalid email or password");
        }

        if (!user.getIsActive()) {
            authMetrics.incrementLoginFailures("account_disabled");
            throw new BadCredentialsException("Account is deactivated");
        }

        rateLimiterService.resetLoginRate(email);
        accountLockoutService.clearFailedAttempts(email);

        String accessToken = jwtTokenProvider.generateAccessToken(user);
        String refreshToken = jwtTokenProvider.generateRefreshToken(user);

        saveRefreshToken(user.getUserId(), refreshToken);

        outboxPublisher.publish("user", user.getUserId().toString(), "user.logged_in", Map.of(
            "user_id", user.getUserId().toString(),
            "email", user.getEmail(),
            "ip", ipAddress
        ));

        authMetrics.incrementLogins();
        log.info("User logged in: {} from IP: {}", email, ipAddress);

        return buildAuthResponse(user, accessToken, refreshToken);
    }

    @Transactional
    public AuthResponse refresh(String refreshTokenValue) {
        Claims claims = jwtTokenProvider.validateRefreshToken(refreshTokenValue);
        if (claims == null) {
            authMetrics.incrementTokenRefreshFailures("invalid_token");
            throw new BadCredentialsException("Invalid or expired refresh token");
        }

        RefreshToken storedToken = refreshTokenRepository.findByToken(refreshTokenValue)
            .orElseThrow(() -> {
                authMetrics.incrementTokenRefreshFailures("not_found");
                return new BadCredentialsException("Refresh token not found");
            });

        if (storedToken.getRevoked()) {
            authMetrics.incrementTokenRefreshFailures("revoked");
            log.warn("Attempted to use revoked refresh token for user: {}", storedToken.getUserId());
            refreshTokenRepository.deleteByUserId(storedToken.getUserId());
            throw new BadCredentialsException("Refresh token has been revoked");
        }

        if (storedToken.getExpiresAt().isBefore(LocalDateTime.now())) {
            authMetrics.incrementTokenRefreshFailures("expired");
            throw new BadCredentialsException("Refresh token has expired");
        }

        storedToken.setRevoked(true);
        refreshTokenRepository.save(storedToken);

        UUID userId = UUID.fromString(claims.getSubject());
        User user = userRepository.findById(userId)
            .orElseThrow(() -> new UsernameNotFoundException("User not found"));

        String newAccessToken = jwtTokenProvider.generateAccessToken(user);
        String newRefreshToken = jwtTokenProvider.generateRefreshToken(user);

        saveRefreshToken(user.getUserId(), newRefreshToken);

        authMetrics.incrementTokenRefreshes();
        log.info("Token refreshed for user: {}", userId);

        return buildAuthResponse(user, newAccessToken, newRefreshToken);
    }

    @Transactional
    public void logout(String userId, String refreshTokenValue) {
        UUID uuid = UUID.fromString(userId);

        if (refreshTokenValue != null && !refreshTokenValue.isBlank()) {
            refreshTokenRepository.findByToken(refreshTokenValue)
                .ifPresent(token -> {
                    token.setRevoked(true);
                    refreshTokenRepository.save(token);
                });
        }

        refreshTokenRepository.deleteByUserId(uuid);

        outboxPublisher.publish("user", userId, "user.logged_out", Map.of(
            "user_id", userId
        ));

        log.info("User logged out: {}", userId);
    }

    @Transactional(readOnly = true)
    public UserResponse getUserById(String userId) {
        User user = userRepository.findById(UUID.fromString(userId))
            .orElseThrow(() -> new UsernameNotFoundException("User not found: " + userId));

        return UserResponse.builder()
            .userId(user.getUserId())
            .email(user.getEmail())
            .phone(user.getPhone())
            .fullName(user.getFullName())
            .role(user.getRole())
            .verified(user.getIsVerified())
            .createdAt(user.getCreatedAt())
            .updatedAt(user.getUpdatedAt())
            .build();
    }

    @Transactional(readOnly = true)
    public boolean validateToken(String token) {
        return jwtTokenProvider.validateAccessToken(token) != null;
    }

    private void saveRefreshToken(UUID userId, String token) {
        RefreshToken refreshToken = RefreshToken.builder()
            .token(token)
            .userId(userId)
            .expiresAt(LocalDateTime.now().plusSeconds(jwtTokenProvider.getRefreshTtlSeconds()))
            .revoked(false)
            .build();
        refreshTokenRepository.save(refreshToken);
    }

    private AuthResponse buildAuthResponse(User user, String accessToken, String refreshToken) {
        return AuthResponse.builder()
            .userId(user.getUserId())
            .email(user.getEmail())
            .phone(user.getPhone())
            .fullName(user.getFullName())
            .role(user.getRole())
            .accessToken(accessToken)
            .refreshToken(refreshToken)
            .expiresIn(jwtTokenProvider.getAccessTtlSeconds())
            .build();
    }

    private void validatePasswordStrength(String password) {
        if (password == null || password.length() < 8) {
            throw new IllegalArgumentException("Password must be at least 8 characters");
        }
        if (!password.matches(".*[A-Z].*")) {
            throw new IllegalArgumentException("Password must contain at least one uppercase letter");
        }
        if (!password.matches(".*[a-z].*")) {
            throw new IllegalArgumentException("Password must contain at least one lowercase letter");
        }
        if (!password.matches(".*\\d.*")) {
            throw new IllegalArgumentException("Password must contain at least one digit");
        }
    }

    private String resolveClientIp() {
        String ip = httpServletRequest.getHeader("X-Forwarded-For");
        if (ip == null || ip.isBlank() || "unknown".equalsIgnoreCase(ip)) {
            ip = httpServletRequest.getHeader("X-Real-IP");
        }
        if (ip == null || ip.isBlank() || "unknown".equalsIgnoreCase(ip)) {
            ip = httpServletRequest.getRemoteAddr();
        }
        if (ip != null && ip.contains(",")) {
            ip = ip.split(",")[0].trim();
        }
        return ip;
    }
}
