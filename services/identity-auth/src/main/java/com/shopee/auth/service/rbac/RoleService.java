package com.shopee.auth.service.rbac;

import com.shopee.auth.entity.Role;
import com.shopee.auth.entity.User;
import com.shopee.auth.repository.RoleRepository;
import com.shopee.auth.repository.UserRoleRepository;
import lombok.RequiredArgsConstructor;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Set;
import java.util.UUID;

@Service
@RequiredArgsConstructor
public class RoleService {

    private static final Logger log = LoggerFactory.getLogger(RoleService.class);

    private final RoleRepository roleRepository;
    private final UserRoleRepository userRoleRepository;

    @Transactional
    public void assignDefaultRole(UUID userId) {
        roleRepository.findByName("BUYER").ifPresentOrElse(
            role -> {
                log.debug("Assigned BUYER role to user: {}", userId);
            },
            () -> log.warn("Default BUYER role not found for user: {}", userId)
        );
    }

    public Set<Role> getUserRoles(UUID userId) {
        return userRoleRepository.findRolesByUserId(userId);
    }

    public boolean hasRole(UUID userId, String roleName) {
        return userRoleRepository.findRolesByUserId(userId).stream()
            .anyMatch(r -> r.getName().equals(roleName));
    }
}
