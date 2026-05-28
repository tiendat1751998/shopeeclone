package com.tiki.auth;

import com.tiki.auth.entity.FailedLoginAttempt;
import com.tiki.auth.repository.FailedLoginRepository;
import com.tiki.auth.security.AccountLockoutService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.ArgumentCaptor;
import org.mockito.Captor;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import static org.assertj.core.api.Assertions.assertThat;
import static org.mockito.ArgumentMatchers.*;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class AccountLockoutServiceTest {

    @Mock
    private FailedLoginRepository failedLoginRepository;

    @Captor
    private ArgumentCaptor<FailedLoginAttempt> attemptCaptor;

    private AccountLockoutService accountLockoutService;

    @BeforeEach
    void setUp() {
        accountLockoutService = new AccountLockoutService(failedLoginRepository);
    }

    @Test
    void isAccountLocked_UnderLimit_ReturnsFalse() {
        when(failedLoginRepository.countByEmailAndAttemptedAtAfter(anyString(), any())).thenReturn(3);
        boolean locked = accountLockoutService.isAccountLocked("test@example.com");
        assertThat(locked).isFalse();
    }

    @Test
    void isAccountLocked_AtLimit_ReturnsTrue() {
        when(failedLoginRepository.countByEmailAndAttemptedAtAfter(anyString(), any())).thenReturn(5);
        boolean locked = accountLockoutService.isAccountLocked("test@example.com");
        assertThat(locked).isTrue();
    }

    @Test
    void isAccountLocked_OverLimit_ReturnsTrue() {
        when(failedLoginRepository.countByEmailAndAttemptedAtAfter(anyString(), any())).thenReturn(10);
        boolean locked = accountLockoutService.isAccountLocked("test@example.com");
        assertThat(locked).isTrue();
    }

    @Test
    void isIpBlocked_NullIp_ReturnsFalse() {
        boolean blocked = accountLockoutService.isIpBlocked(null);
        assertThat(blocked).isFalse();
    }

    @Test
    void recordFailedAttempt_SavesAttempt() {
        accountLockoutService.recordFailedAttempt("test@example.com", "192.168.1.1");
        verify(failedLoginRepository).save(attemptCaptor.capture());
        FailedLoginAttempt saved = attemptCaptor.getValue();
        assertThat(saved.getEmail()).isEqualTo("test@example.com");
        assertThat(saved.getIpAddress()).isEqualTo("192.168.1.1");
        assertThat(saved.getAttemptedAt()).isNotNull();
    }

    @Test
    void recordFailedAttempt_NormalizesEmailToLowercase() {
        accountLockoutService.recordFailedAttempt("Test@Example.COM", "10.0.0.1");
        verify(failedLoginRepository).save(attemptCaptor.capture());
        assertThat(attemptCaptor.getValue().getEmail()).isEqualTo("test@example.com");
    }
}
