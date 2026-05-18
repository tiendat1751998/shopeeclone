package com.shopee.auth.repository;

import com.shopee.auth.entity.FailedLoginAttempt;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.time.LocalDateTime;
import java.util.UUID;

@Repository
public interface FailedLoginRepository extends JpaRepository<FailedLoginAttempt, UUID> {

    int countByEmailAndAttemptedAtAfter(String email, LocalDateTime after);

    int countByIpAddressAndAttemptedAtAfter(String ipAddress, LocalDateTime after);

    void deleteByAttemptedAtBefore(LocalDateTime before);
}
