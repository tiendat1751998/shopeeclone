package com.tiki.auth.security;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.redis.core.RedisTemplate;
import org.springframework.stereotype.Service;

import java.time.Duration;
import java.util.concurrent.TimeUnit;

@Service
public class RateLimiterService {

    private static final Logger log = LoggerFactory.getLogger(RateLimiterService.class);

    private final RedisTemplate<String, String> redisTemplate;

    @Value("${rate-limit.login.max-attempts:5}")
    private int loginMaxAttempts;

    @Value("${rate-limit.login.window-seconds:300}")
    private int loginWindowSeconds;

    @Value("${rate-limit.ip.max-attempts:20}")
    private int ipMaxAttempts;

    @Value("${rate-limit.ip.window-seconds:60}")
    private int ipWindowSeconds;

    private static final String LOGIN_RATE_KEY = "ratelimit:login:%s";
    private static final String IP_RATE_KEY = "ratelimit:ip:%s";

    public RateLimiterService(RedisTemplate<String, String> redisTemplate) {
        this.redisTemplate = redisTemplate;
    }

    public boolean isLoginAllowed(String email) {
        String key = String.format(LOGIN_RATE_KEY, email.toLowerCase());
        return checkRate(key, loginMaxAttempts, loginWindowSeconds);
    }

    public boolean isIpAllowed(String ipAddress) {
        if (ipAddress == null || ipAddress.isBlank()) {
            return true;
        }
        String key = String.format(IP_RATE_KEY, ipAddress);
        return checkRate(key, ipMaxAttempts, ipWindowSeconds);
    }

    public void recordLoginAttempt(String email) {
        String key = String.format(LOGIN_RATE_KEY, email.toLowerCase());
        incrementRate(key, loginWindowSeconds);
    }

    public void recordIpAttempt(String ipAddress) {
        if (ipAddress == null || ipAddress.isBlank()) {
            return;
        }
        String key = String.format(IP_RATE_KEY, ipAddress);
        incrementRate(key, ipWindowSeconds);
    }

    public void resetLoginRate(String email) {
        String key = String.format(LOGIN_RATE_KEY, email.toLowerCase());
        redisTemplate.delete(key);
    }

    private boolean checkRate(String key, int maxAttempts, int windowSeconds) {
        String val = redisTemplate.opsForValue().get(key);
        if (val == null) {
            return true;
        }
        int attempts = Integer.parseInt(val);
        return attempts < maxAttempts;
    }

    private void incrementRate(String key, int windowSeconds) {
        Long count = redisTemplate.opsForValue().increment(key);
        if (count != null && count == 1) {
            redisTemplate.expire(key, Duration.ofSeconds(windowSeconds));
        }
    }
}
