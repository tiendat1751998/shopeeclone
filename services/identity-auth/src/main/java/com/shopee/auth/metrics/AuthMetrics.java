package com.shopee.auth.metrics;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import org.springframework.stereotype.Component;

import java.util.concurrent.TimeUnit;

@Component
public class AuthMetrics {

    private final Counter registrations;
    private final Counter logins;
    private final Counter loginFailures;
    private final Counter registrationFailures;
    private final Counter tokenRefreshes;
    private final Counter tokenRefreshFailures;
    private final Timer loginDuration;
    private final Timer registrationDuration;

    public AuthMetrics(MeterRegistry registry) {
        this.registrations = Counter.builder("shopee.auth.registrations")
            .description("Total successful registrations")
            .register(registry);

        this.logins = Counter.builder("shopee.auth.logins")
            .description("Total successful logins")
            .register(registry);

        this.loginFailures = Counter.builder("shopee.auth.login.failures")
            .description("Total failed login attempts")
            .tag("reason", "unknown")
            .register(registry);

        this.registrationFailures = Counter.builder("shopee.auth.registration.failures")
            .description("Total failed registrations")
            .tag("reason", "unknown")
            .register(registry);

        this.tokenRefreshes = Counter.builder("shopee.auth.token.refreshes")
            .description("Total successful token refreshes")
            .register(registry);

        this.tokenRefreshFailures = Counter.builder("shopee.auth.token.refresh.failures")
            .description("Total failed token refreshes")
            .tag("reason", "unknown")
            .register(registry);

        this.loginDuration = Timer.builder("shopee.auth.login.duration")
            .description("Login request duration")
            .publishPercentiles(0.5, 0.9, 0.95, 0.99)
            .register(registry);

        this.registrationDuration = Timer.builder("shopee.auth.registration.duration")
            .description("Registration request duration")
            .publishPercentiles(0.5, 0.9, 0.95, 0.99)
            .register(registry);
    }

    public void incrementRegistrations() {
        registrations.increment();
    }

    public void incrementLogins() {
        logins.increment();
    }

    public void incrementLoginFailures(String reason) {
        loginFailures.increment();
    }

    public void incrementRegistrationFailures(String reason) {
        registrationFailures.increment();
    }

    public void incrementTokenRefreshes() {
        tokenRefreshes.increment();
    }

    public void incrementTokenRefreshFailures(String reason) {
        tokenRefreshFailures.increment();
    }

    public <T> T recordLoginDuration(java.util.concurrent.Callable<T> callable) throws Exception {
        return loginDuration.call(callable);
    }

    public <T> T recordRegistrationDuration(java.util.concurrent.Callable<T> callable) throws Exception {
        return registrationDuration.call(callable);
    }
}
