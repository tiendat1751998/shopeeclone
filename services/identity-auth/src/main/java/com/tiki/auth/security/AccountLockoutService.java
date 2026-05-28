package com.tiki.auth.security;

import com.tiki.auth.entity.FailedLoginAttempt;
import com.tiki.auth.repository.FailedLoginRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import java.time.LocalDateTime;

@Service
public class AccountLockoutService {

    private static final Logger log = LoggerFactory.getLogger(AccountLockoutService.class);

    private final FailedLoginRepository failedLoginRepository;

    @Value("${account-lockout.max-failed-attempts:5}")
    private int maxFailedAttempts;

    @Value("${account-lockout.lockout-duration-minutes:15}")
    private int lockoutDurationMinutes;

    public AccountLockoutService(FailedLoginRepository failedLoginRepository) {
        this.failedLoginRepository = failedLoginRepository;
    }

    public void recordFailedAttempt(String email, String ipAddress) {
        FailedLoginAttempt attempt = FailedLoginAttempt.builder()
            .email(email.toLowerCase())
            .ipAddress(ipAddress)
            .build();
        failedLoginRepository.save(attempt);

        log.warn("Failed login attempt recorded for email: {} from IP: {}", email, ipAddress);
    }

    public boolean isAccountLocked(String email) {
        LocalDateTime since = LocalDateTime.now().minusMinutes(lockoutDurationMinutes);
        int attempts = failedLoginRepository.countByEmailAndAttemptedAtAfter(email.toLowerCase(), since);
        return attempts >= maxFailedAttempts;
    }

    public boolean isIpBlocked(String ipAddress) {
        if (ipAddress == null) {
            return false;
        }
        int maxIpAttempts = 20;
        LocalDateTime since = LocalDateTime.now().minusMinutes(15);
        int attempts = failedLoginRepository.countByIpAddressAndAttemptedAtAfter(ipAddress, since);
        return attempts >= maxIpAttempts;
    }

    public void clearFailedAttempts(String email) {
        log.info("Clearing failed login attempts for email: {}", email);
        failedLoginRepository.deleteByEmail(email.toLowerCase());
    }
}
