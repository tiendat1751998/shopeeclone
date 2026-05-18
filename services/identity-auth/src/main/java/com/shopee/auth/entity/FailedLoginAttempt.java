package com.shopee.auth.entity;

import jakarta.persistence.*;
import lombok.AllArgsConstructor;
import lombok.Builder;
import lombok.Data;
import lombok.NoArgsConstructor;

import java.time.LocalDateTime;
import java.util.UUID;

@Entity
@Table(name = "failed_login_attempts", indexes = {
    @Index(name = "idx_failed_login_email", columnList = "email,attemptedAt"),
    @Index(name = "idx_failed_login_ip", columnList = "ipAddress,attemptedAt")
})
@Data
@Builder
@NoArgsConstructor
@AllArgsConstructor
public class FailedLoginAttempt {

    @Id
    @Column(name = "id", nullable = false, updatable = false)
    private UUID id;

    @Column(name = "email", nullable = false, length = 255)
    private String email;

    @Column(name = "ip_address", nullable = false, length = 45)
    private String ipAddress;

    @Column(name = "attempted_at", nullable = false, updatable = false)
    private LocalDateTime attemptedAt;

    @PrePersist
    public void prePersist() {
        if (id == null) {
            id = UUID.randomUUID();
        }
        if (attemptedAt == null) {
            attemptedAt = LocalDateTime.now();
        }
    }
}
