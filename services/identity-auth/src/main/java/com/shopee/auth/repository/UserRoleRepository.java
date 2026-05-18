package com.shopee.auth.repository;

import com.shopee.auth.entity.Role;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.stereotype.Repository;

import java.util.Set;
import java.util.UUID;

@Repository
public interface UserRoleRepository extends JpaRepository<Object, UUID> {

    @Query(value = "SELECT r.* FROM roles r JOIN user_roles ur ON r.role_id = ur.role_id WHERE ur.user_id = ?1", nativeQuery = true)
    Set<Role> findRolesByUserId(UUID userId);

    @Query(value = "SELECT p.* FROM permissions p " +
           "JOIN role_permissions rp ON p.permission_id = rp.permission_id " +
           "JOIN user_roles ur ON rp.role_id = ur.role_id " +
           "WHERE ur.user_id = ?1", nativeQuery = true)
    Set<Object> findPermissionsByUserId(UUID userId);
}
