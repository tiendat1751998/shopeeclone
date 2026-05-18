package com.shopee.auth;

import com.shopee.auth.security.RateLimiterService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.data.redis.core.ValueOperations;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class RateLimiterServiceTest {

    @Mock
    private RedisTemplate<String, String> redisTemplate;

    @Mock
    private ValueOperations<String, String> valueOps;

    private RateLimiterService rateLimiterService;

    @BeforeEach
    void setUp() {
        rateLimiterService = new RateLimiterService(redisTemplate);
        when(redisTemplate.opsForValue()).thenReturn(valueOps);
    }

    @Test
    void isLoginAllowed_NoPreviousAttempts_ReturnsTrue() {
        when(valueOps.get(anyString())).thenReturn(null);
        boolean allowed = rateLimiterService.isLoginAllowed("test@example.com");
        assertThat(allowed).isTrue();
    }

    @Test
    void isLoginAllowed_UnderLimit_ReturnsTrue() {
        when(valueOps.get(anyString())).thenReturn("3");
        boolean allowed = rateLimiterService.isLoginAllowed("test@example.com");
        assertThat(allowed).isTrue();
    }

    @Test
    void isLoginAllowed_OverLimit_ReturnsFalse() {
        when(valueOps.get(anyString())).thenReturn("5");
        boolean allowed = rateLimiterService.isLoginAllowed("test@example.com");
        assertThat(allowed).isFalse();
    }

    @Test
    void isIpAllowed_NoPreviousAttempts_ReturnsTrue() {
        when(valueOps.get(anyString())).thenReturn(null);
        boolean allowed = rateLimiterService.isIpAllowed("192.168.1.1");
        assertThat(allowed).isTrue();
    }

    @Test
    void isIpAllowed_UnderLimit_ReturnsTrue() {
        when(valueOps.get(anyString())).thenReturn("10");
        boolean allowed = rateLimiterService.isIpAllowed("192.168.1.1");
        assertThat(allowed).isTrue();
    }

    @Test
    void isIpAllowed_OverLimit_ReturnsFalse() {
        when(valueOps.get(anyString())).thenReturn("25");
        boolean allowed = rateLimiterService.isIpAllowed("192.168.1.1");
        assertThat(allowed).isFalse();
    }

    @Test
    void isIpAllowed_NullIp_ReturnsTrue() {
        boolean allowed = rateLimiterService.isIpAllowed(null);
        assertThat(allowed).isTrue();
    }

    @Test
    void isIpAllowed_EmptyIp_ReturnsTrue() {
        boolean allowed = rateLimiterService.isIpAllowed("");
        assertThat(allowed).isTrue();
    }

    @Test
    void recordLoginAttempt_IncrementsCounter() {
        when(valueOps.increment(anyString())).thenReturn(1L);
        rateLimiterService.recordLoginAttempt("test@example.com");
        verify(valueOps).increment(contains("ratelimit:login:"));
    }

    @Test
    void recordIpAttempt_IncrementsCounter() {
        when(valueOps.increment(anyString())).thenReturn(1L);
        rateLimiterService.recordIpAttempt("192.168.1.1");
        verify(valueOps).increment(contains("ratelimit:ip:"));
    }

    @Test
    void recordIpAttempt_NullIp_DoesNothing() {
        rateLimiterService.recordIpAttempt(null);
        verify(valueOps, never()).increment(anyString());
    }

    @Test
    void resetLoginRate_DeletesKey() {
        rateLimiterService.resetLoginRate("test@example.com");
        verify(redisTemplate).delete(contains("ratelimit:login:"));
    }
}
