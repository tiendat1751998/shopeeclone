package com.shopee.auth.repository;

import com.shopee.auth.entity.Role;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import java.util.Optional;
import java.util.Set;
import java.util.UUID;

@Repository
public interface RoleRepository extends JpaRepository<Role, UUID> {

    Optional<Role> findByName(String name);

    @Query("SELECT DISTINCT r FROM Role r JOIN r.permissions p WHERE p.resource = :resource AND p.action = :action")
    Set<Role> findRolesByPermission(@Param("resource") String resource, @Param("action") String action);
}
