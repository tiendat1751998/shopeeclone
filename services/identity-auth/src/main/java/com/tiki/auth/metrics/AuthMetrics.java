package com.tiki.auth.metrics;

import io.micrometer.core.instrument.Counter;
import io.micrometer.core.instrument.MeterRegistry;
import io.micrometer.core.instrument.Timer;
import org.springframework.stereotype.Component;

import java.util.concurrent.Callable;
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
        this.registrations = Counter.builder("tiki.auth.registrations")
                .description("Total successful registrations")
                .register(registry);

        this.logins = Counter.builder("tiki.auth.logins")
                .description("Total successful logins")
                .register(registry);

        this.loginFailures = Counter.builder("tiki.auth.login.failures")
                .description("Total failed login attempts")
                .tag("reason", "unknown")
                .register(registry);

        this.registrationFailures = Counter.builder("tiki.auth.registration.failures")
                .description("Total failed registrations")
                .tag("reason", "unknown")
                .register(registry);

        this.tokenRefreshes = Counter.builder("tiki.auth.token.refreshes")
                .description("Total successful token refreshes")
                .register(registry);

        this.tokenRefreshFailures = Counter.builder("tiki.auth.token.refresh.failures")
                .description("Total failed token refreshes")
                .tag("reason", "unknown")
                .register(registry);

        this.loginDuration = Timer.builder("tiki.auth.login.duration")
                .description("Login request duration")
                .publishPercentiles(0.5, 0.9, 0.95, 0.99)
                .register(registry);

        this.registrationDuration = Timer.builder("tiki.auth.registration.duration")
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

    public <T> T recordLoginDuration(Callable<T> callable) {
        try {
            return loginDuration.recordCallable(callable);
        } catch (Exception e) {
            if (e instanceof RuntimeException) {
                throw (RuntimeException) e;
            }
            throw new RuntimeException(e);
        }
    }

    public <T> T recordRegistrationDuration(Callable<T> callable) {
        try {
            return registrationDuration.recordCallable(callable);
        } catch (Exception e) {
            if (e instanceof RuntimeException) {
                throw (RuntimeException) e;
            }
            throw new RuntimeException(e);
        }
    }
}
