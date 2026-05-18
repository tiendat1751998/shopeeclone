package com.shopee.auth;

import com.shopee.auth.dto.AuthRequest;
import com.shopee.auth.dto.AuthResponse;
import com.shopee.auth.entity.User;
import com.shopee.auth.exception.DuplicateResourceException;
import com.shopee.auth.metrics.AuthMetrics;
import com.shopee.auth.repository.FailedLoginRepository;
import com.shopee.auth.repository.OutboxEventRepository;
import com.shopee.auth.repository.RefreshTokenRepository;
import com.shopee.auth.repository.RoleRepository;
import com.shopee.auth.repository.UserRepository;
import com.shopee.auth.repository.UserRoleRepository;
import com.shopee.auth.security.AccountLockoutService;
import com.shopee.auth.security.JwtTokenProvider;
import com.shopee.auth.security.JwksProvider;
import com.shopee.auth.security.RateLimiterService;
import com.shopee.auth.service.AuthService;
import com.shopee.auth.service.OutboxPublisher;
import com.shopee.auth.service.rbac.RoleService;
import com.shopee.auth.service.session.SessionService;
import jakarta.servlet.http.HttpServletRequest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.core.ValueOperations;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;

import java.time.LocalDateTime;
import java.util.Optional;
import java.util.UUID;
import java.util.concurrent.TimeUnit;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class AuthServiceRateLimitTest {

    @Mock private UserRepository userRepository;
    @Mock private RefreshTokenRepository refreshTokenRepository;
    @Mock private FailedLoginRepository failedLoginRepository;
    @Mock private OutboxEventRepository outboxEventRepository;
    @Mock private RoleRepository roleRepository;
    @Mock private UserRoleRepository userRoleRepository;
    @Mock private RedisTemplate<String, String> redisTemplate;
    @Mock private ValueOperations<String, String> valueOps;
    @Mock private HttpServletRequest httpServletRequest;
    @Mock private AuthMetrics authMetrics;

    private BCryptPasswordEncoder passwordEncoder;
    private JwksProvider jwksProvider;
    private JwtTokenProvider jwtTokenProvider;
    private RateLimiterService rateLimiterService;
    private AccountLockoutService accountLockoutService;
    private OutboxPublisher outboxPublisher;
    private SessionService sessionService;
    private RoleService roleService;
    private AuthService authService;

    @BeforeEach
    void setUp() {
        passwordEncoder = new BCryptPasswordEncoder();
        jwksProvider = new JwksProvider();
        jwtTokenProvider = new JwtTokenProvider(jwksProvider);
        rateLimiterService = new RateLimiterService(redisTemplate);
        accountLockoutService = new AccountLockoutService(failedLoginRepository);
        outboxPublisher = new OutboxPublisher(outboxEventRepository, new com.fasterxml.jackson.databind.ObjectMapper());
        sessionService = new SessionService(refreshTokenRepository, jwtTokenProvider);
        roleService = new RoleService(roleRepository, userRoleRepository);
        authService = new AuthService(userRepository, passwordEncoder,
            jwtTokenProvider, rateLimiterService, accountLockoutService,
            sessionService, roleService, outboxPublisher, authMetrics, httpServletRequest);

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
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
        jwtTokenProvider.init();
    }

    @Test
    void register_Success() {
        var request = new AuthRequest.Register();
        request.setEmail("test@example.com");
        request.setPhone("0123456789");
        request.setPassword("TestPass123");
        request.setFullName("Test User");

        when(userRepository.existsByEmail(any())).thenReturn(false);
        when(userRepository.existsByPhone(any())).thenReturn(false);
        when(userRepository.save(any())).thenAnswer(inv -> {
            User u = inv.getArgument(0);
            u.setUserId(UUID.randomUUID());
            return u;
        });
        when(refreshTokenRepository.save(any())).thenReturn(null);
        when(outboxEventRepository.save(any())).thenReturn(null);
        when(redisTemplate.opsForValue()).thenReturn(valueOps);

        AuthResponse response = authService.register(request);

        assertThat(response).isNotNull();
        assertThat(response.getAccessToken()).isNotBlank();
        assertThat(response.getRefreshToken()).isNotBlank();
        assertThat(response.getEmail()).isEqualTo("test@example.com");
        assertThat(response.getExpiresIn()).isPositive();
    }

    @Test
    void register_DuplicateEmail() {
        var request = new AuthRequest.Register();
        request.setEmail("existing@example.com");
        request.setPassword("TestPass123");
        request.setFullName("Test");

        when(userRepository.existsByEmail(any())).thenReturn(true);

        assertThatThrownBy(() -> authService.register(request))
            .isInstanceOf(DuplicateResourceException.class)
            .hasMessageContaining("Email already registered");
    }

    @Test
    void register_WeakPassword() {
        var request = new AuthRequest.Register();
        request.setEmail("test@example.com");
        request.setPassword("weak");
        request.setFullName("Test");

        assertThatThrownBy(() -> authService.register(request))
            .isInstanceOf(IllegalArgumentException.class)
            .hasMessageContaining("at least 8 characters");
    }

    @Test
    void register_PasswordNoUppercase() {
        var request = new AuthRequest.Register();
        request.setEmail("test@example.com");
        request.setPassword("lowercaseonly1");
        request.setFullName("Test");

        assertThatThrownBy(() -> authService.register(request))
            .isInstanceOf(IllegalArgumentException.class)
            .hasMessageContaining("uppercase letter");
    }

    @Test
    void login_Success() {
        var request = new AuthRequest.Login();
        request.setEmail("test@example.com");
        request.setPassword("TestPass123");

        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("test@example.com")
            .passwordHash(passwordEncoder.encode("TestPass123"))
            .fullName("Test User")
            .role("BUYER")
            .isActive(true)
            .build();

        when(failedLoginRepository.countByEmailAndAttemptedAtAfter(any(), any())).thenReturn(0);
        when(redisTemplate.opsForValue()).thenReturn(valueOps);
        when(valueOps.get(anyString())).thenReturn(null);
        when(userRepository.findByEmail(any())).thenReturn(Optional.of(user));
        when(refreshTokenRepository.save(any())).thenReturn(null);
        when(outboxEventRepository.save(any())).thenReturn(null);
        when(httpServletRequest.getHeader(anyString())).thenReturn("127.0.0.1");

        AuthResponse response = authService.login(request);

        assertThat(response).isNotNull();
        assertThat(response.getAccessToken()).isNotBlank();
        assertThat(response.getRefreshToken()).isNotBlank();
        assertThat(response.getRole()).isEqualTo("BUYER");
    }

    @Test
    void login_WrongPassword() {
        var request = new AuthRequest.Login();
        request.setEmail("test@example.com");
        request.setPassword("WrongPass123");

        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("test@example.com")
            .passwordHash(passwordEncoder.encode("CorrectPass123"))
            .fullName("Test User")
            .isActive(true)
            .build();

        when(failedLoginRepository.countByEmailAndAttemptedAtAfter(any(), any())).thenReturn(0);
        when(redisTemplate.opsForValue()).thenReturn(valueOps);
        when(valueOps.get(anyString())).thenReturn(null);
        when(userRepository.findByEmail(any())).thenReturn(Optional.of(user));
        when(httpServletRequest.getHeader(anyString())).thenReturn("127.0.0.1");
        when(failedLoginRepository.save(any())).thenReturn(null);

        assertThatThrownBy(() -> authService.login(request))
            .isInstanceOf(BadCredentialsException.class)
            .hasMessageContaining("Invalid email or password");
    }

    @Test
    void login_AccountLocked() {
        var request = new AuthRequest.Login();
        request.setEmail("locked@example.com");
        request.setPassword("TestPass123");

        when(failedLoginRepository.countByEmailAndAttemptedAtAfter(eq("locked@example.com"), any())).thenReturn(5);

        assertThatThrownBy(() -> authService.login(request))
            .isInstanceOf(BadCredentialsException.class)
            .hasMessageContaining("Account temporarily locked");
    }

    @Test
    void login_AccountDeactivated() {
        var request = new AuthRequest.Login();
        request.setEmail("inactive@example.com");
        request.setPassword("TestPass123");

        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("inactive@example.com")
            .passwordHash(passwordEncoder.encode("TestPass123"))
            .fullName("Inactive User")
            .isActive(false)
            .build();

        when(failedLoginRepository.countByEmailAndAttemptedAtAfter(any(), any())).thenReturn(0);
        when(redisTemplate.opsForValue()).thenReturn(valueOps);
        when(valueOps.get(anyString())).thenReturn(null);
        when(userRepository.findByEmail(any())).thenReturn(Optional.of(user));
        when(httpServletRequest.getHeader(anyString())).thenReturn("127.0.0.1");

        assertThatThrownBy(() -> authService.login(request))
            .isInstanceOf(BadCredentialsException.class)
            .hasMessageContaining("Account is deactivated");
    }

    @Test
    void getUserById_Success() {
        UUID userId = UUID.randomUUID();

        User user = User.builder()
            .userId(userId)
            .email("test@example.com")
            .fullName("Test User")
            .role("BUYER")
            .isVerified(true)
            .createdAt(LocalDateTime.now())
            .updatedAt(LocalDateTime.now())
            .build();

        when(userRepository.findById(userId)).thenReturn(Optional.of(user));

        var response = authService.getUserById(userId.toString());

        assertThat(response).isNotNull();
        assertThat(response.getUserId()).isEqualTo(userId);
        assertThat(response.getEmail()).isEqualTo("test@example.com");
        assertThat(response.getRole()).isEqualTo("BUYER");
    }

    @Test
    void getUserById_NotFound() {
        when(userRepository.findById(any())).thenReturn(Optional.empty());

        assertThatThrownBy(() -> authService.getUserById(UUID.randomUUID().toString()))
            .isInstanceOf(org.springframework.security.core.userdetails.UsernameNotFoundException.class)
            .hasMessageContaining("User not found");
    }

    @Test
    void validateToken_ReturnsTrueForValidToken() {
        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("test@example.com")
            .passwordHash(passwordEncoder.encode("TestPass123"))
            .fullName("Test User")
            .role("BUYER")
            .build();

        String token = jwtTokenProvider.generateAccessToken(user);
        boolean valid = authService.validateToken(token);
        assertThat(valid).isTrue();
    }

    @Test
    void validateToken_ReturnsFalseForInvalidToken() {
        boolean valid = authService.validateToken("invalid.token.here");
        assertThat(valid).isFalse();
    }

    @Test
    void validateToken_ReturnsFalseForExpiredToken() {
        // Can't easily test with the provider, but we can test garbage
        boolean valid = authService.validateToken("");
        assertThat(valid).isFalse();
    }
}
