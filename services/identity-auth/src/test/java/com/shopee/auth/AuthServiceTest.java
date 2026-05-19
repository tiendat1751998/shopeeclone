package com.shopee.auth;

import com.shopee.auth.dto.AuthRequest;
import com.shopee.auth.dto.AuthResponse;
import com.shopee.auth.dto.UserResponse;
import com.shopee.auth.entity.RefreshToken;
import com.shopee.auth.entity.User;
import com.shopee.auth.exception.DuplicateResourceException;
import com.shopee.auth.metrics.AuthMetrics;
import com.shopee.auth.repository.OutboxEventRepository;
import com.shopee.auth.repository.RefreshTokenRepository;
import com.shopee.auth.repository.RoleRepository;
import com.shopee.auth.repository.UserRepository;
import com.shopee.auth.repository.UserRoleRepository;
import com.shopee.auth.security.AccountLockoutService;
import com.shopee.auth.security.JwksProvider;
import com.shopee.auth.security.JwtTokenProvider;
import com.shopee.auth.security.RateLimiterService;
import com.shopee.auth.service.AuthService;
import com.shopee.auth.service.OutboxPublisher;
import com.shopee.auth.service.rbac.RoleService;
import com.shopee.auth.service.session.SessionService;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.servlet.http.HttpServletRequest;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.security.authentication.BadCredentialsException;
import org.springframework.security.core.userdetails.UsernameNotFoundException;
import org.springframework.security.crypto.bcrypt.BCryptPasswordEncoder;
import org.springframework.security.crypto.password.PasswordEncoder;

import java.time.LocalDateTime;
import java.util.Optional;
import java.util.UUID;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class AuthServiceTest {

    @Mock
    private UserRepository userRepository;

    @Mock
    private RefreshTokenRepository refreshTokenRepository;

    @Mock
    private RoleRepository roleRepository;

    @Mock
    private UserRoleRepository userRoleRepository;

    @Mock
    private HttpServletRequest httpServletRequest;

    @Mock
    private AuthMetrics authMetrics;

    @Mock
    private OutboxEventRepository outboxEventRepository;

    @Mock
    private JwksProvider jwksProvider;

    private PasswordEncoder passwordEncoder;
    private JwtTokenProvider jwtTokenProvider;
    private SessionService sessionService;
    private RoleService roleService;
    private AuthService authService;

    @BeforeEach
    void setUp() {
        passwordEncoder = new BCryptPasswordEncoder();
        jwtTokenProvider = new JwtTokenProvider(jwksProvider);
        jwtTokenProvider.init();

        setField(jwtTokenProvider, "accessSecret", "test-access-secret-key-must-be-at-least-256-bits");
        setField(jwtTokenProvider, "refreshSecret", "test-refresh-secret-key-must-be-at-least-256-bits");
        setField(jwtTokenProvider, "accessTtlSeconds", 900L);
        setField(jwtTokenProvider, "refreshTtlSeconds", 604800L);

        sessionService = new SessionService(refreshTokenRepository, jwtTokenProvider);
        roleService = new RoleService(roleRepository, userRoleRepository);

        ObjectMapper objectMapper = new ObjectMapper();
        OutboxPublisher outboxPublisher = new OutboxPublisher(outboxEventRepository, objectMapper);

        var rateLimiterService = mock(RateLimiterService.class);
        var accountLockoutService = mock(AccountLockoutService.class);

        authService = new AuthService(userRepository, passwordEncoder,
            jwtTokenProvider, rateLimiterService, accountLockoutService,
            sessionService, roleService, outboxPublisher, authMetrics, httpServletRequest);
    }

    private void setField(Object target, String fieldName, Object value) {
        try {
            var field = target.getClass().getDeclaredField(fieldName);
            field.setAccessible(true);
            field.set(target, value);
        } catch (Exception e) {
            throw new RuntimeException(e);
        }
    }

    @Test
    void register_Success() {
        var request = new AuthRequest.Register();
        request.setEmail("test@example.com");
        request.setPhone("0123456789");
        request.setPassword("password123");
        request.setFullName("Test User");

        when(userRepository.existsByEmail(any())).thenReturn(false);
        when(userRepository.existsByPhone(any())).thenReturn(false);
        when(userRepository.save(any())).thenAnswer(inv -> {
            User u = inv.getArgument(0);
            u.setUserId(UUID.randomUUID());
            return u;
        });
        when(refreshTokenRepository.save(any())).thenReturn(new RefreshToken());

        AuthResponse response = authService.register(request);

        assertThat(response).isNotNull();
        assertThat(response.getAccessToken()).isNotBlank();
        assertThat(response.getRefreshToken()).isNotBlank();
        assertThat(response.getEmail()).isEqualTo("test@example.com");
    }

    @Test
    void register_DuplicateEmail() {
        var request = new AuthRequest.Register();
        request.setEmail("existing@example.com");
        request.setPassword("password123");
        request.setFullName("Test");

        when(userRepository.existsByEmail(any())).thenReturn(true);

        assertThatThrownBy(() -> authService.register(request))
            .isInstanceOf(DuplicateResourceException.class)
            .hasMessageContaining("Email already registered");
    }

    @Test
    void login_Success() {
        var request = new AuthRequest.Login();
        request.setEmail("test@example.com");
        request.setPassword("password123");

        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("test@example.com")
            .passwordHash(passwordEncoder.encode("password123"))
            .fullName("Test User")
            .role("BUYER")
            .isActive(true)
            .build();

        when(userRepository.findByEmail(any())).thenReturn(Optional.of(user));
        when(refreshTokenRepository.save(any())).thenReturn(new RefreshToken());

        AuthResponse response = authService.login(request);

        assertThat(response).isNotNull();
        assertThat(response.getAccessToken()).isNotBlank();
    }

    @Test
    void login_WrongPassword() {
        var request = new AuthRequest.Login();
        request.setEmail("test@example.com");
        request.setPassword("wrongpassword");

        User user = User.builder()
            .userId(UUID.randomUUID())
            .email("test@example.com")
            .passwordHash(passwordEncoder.encode("correctpassword"))
            .fullName("Test User")
            .build();

        when(userRepository.findByEmail(any())).thenReturn(Optional.of(user));

        assertThatThrownBy(() -> authService.login(request))
            .isInstanceOf(BadCredentialsException.class);
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

        UserResponse response = authService.getUserById(userId.toString());

        assertThat(response).isNotNull();
        assertThat(response.getUserId()).isEqualTo(userId);
        assertThat(response.getEmail()).isEqualTo("test@example.com");
    }

    @Test
    void getUserById_NotFound() {
        when(userRepository.findById(any())).thenReturn(Optional.empty());

        assertThatThrownBy(() -> authService.getUserById(UUID.randomUUID().toString()))
            .isInstanceOf(UsernameNotFoundException.class);
    }
}
