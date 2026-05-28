package com.tiki.auth.repository;

import com.tiki.auth.entity.Permission;
import com.tiki.auth.entity.Role;
import jakarta.persistence.EntityManager;
import jakarta.persistence.PersistenceContext;
import jakarta.persistence.Query;
import org.springframework.stereotype.Repository;

import java.util.HashSet;
import java.util.List;
import java.util.Set;
import java.util.UUID;

@Repository
public class UserRoleRepository {

    @PersistenceContext
    private EntityManager em;

    @SuppressWarnings("unchecked")
    public Set<Role> findRolesByUserId(UUID userId) {
        Query query = em.createNativeQuery(
                "SELECT r.* FROM roles r JOIN user_roles ur ON r.role_id = ur.role_id WHERE ur.user_id = ?1",
                Role.class);
        query.setParameter(1, userId.toString());
        List<Role> results = query.getResultList();
        return new HashSet<>(results);
    }

    @SuppressWarnings("unchecked")
    public Set<Permission> findPermissionsByUserId(UUID userId) {
        Query query = em.createNativeQuery(
                "SELECT DISTINCT p.* FROM permissions p " +
                "JOIN role_permissions rp ON p.permission_id = rp.permission_id " +
                "JOIN user_roles ur ON rp.role_id = ur.role_id " +
                "WHERE ur.user_id = ?1",
                Permission.class);
        query.setParameter(1, userId.toString());
        List<Permission> results = query.getResultList();
        return new HashSet<>(results);
    }
}
