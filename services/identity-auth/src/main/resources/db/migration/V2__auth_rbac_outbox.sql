-- ============================================
-- RBAC: Roles & Permissions
-- ============================================
CREATE TABLE IF NOT EXISTS roles (
    role_id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS permissions (
    permission_id VARCHAR(36) PRIMARY KEY,
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(resource, action)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id VARCHAR(36) NOT NULL,
    permission_id VARCHAR(36) NOT NULL,
    PRIMARY KEY (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(permission_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS user_roles (
    user_id VARCHAR(36) NOT NULL,
    role_id VARCHAR(36) NOT NULL,
    assigned_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_role_permissions_role ON role_permissions(role_id);

-- ============================================
-- Outbox Pattern
-- ============================================
CREATE TABLE IF NOT EXISTS outbox_events (
    event_id VARCHAR(36) PRIMARY KEY,
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSON NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMP NULL DEFAULT NULL,
    retry_count INT NOT NULL DEFAULT 0,
    last_error TEXT
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_outbox_processed ON outbox_events(processed, created_at);

-- ============================================
-- Failed Login Attempts (Account Lockout)
-- ============================================
CREATE TABLE IF NOT EXISTS failed_login_attempts (
    id VARCHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    attempted_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE INDEX idx_failed_login_email ON failed_login_attempts(email, attempted_at);
CREATE INDEX idx_failed_login_ip ON failed_login_attempts(ip_address, attempted_at);

-- ============================================
-- Seed default roles and permissions
-- ============================================
INSERT IGNORE INTO roles (role_id, name, description, is_system) VALUES
    (UUID(), 'SUPER_ADMIN', 'Full system access', TRUE),
    (UUID(), 'ADMIN', 'Administrative access', TRUE),
    (UUID(), 'SELLER', 'Seller/merchant account', TRUE),
    (UUID(), 'BUYER', 'Standard buyer account', TRUE);

INSERT IGNORE INTO permissions (permission_id, resource, action, description) VALUES
    (UUID(), 'users', 'read', 'View user profiles'),
    (UUID(), 'users', 'write', 'Create/update users'),
    (UUID(), 'users', 'delete', 'Delete users'),
    (UUID(), 'products', 'read', 'View products'),
    (UUID(), 'products', 'write', 'Create/update products'),
    (UUID(), 'products', 'delete', 'Delete products'),
    (UUID(), 'orders', 'read', 'View orders'),
    (UUID(), 'orders', 'write', 'Create/update orders'),
    (UUID(), 'orders', 'cancel', 'Cancel orders'),
    (UUID(), 'payments', 'read', 'View payments'),
    (UUID(), 'payments', 'refund', 'Process refunds'),
    (UUID(), 'inventory', 'read', 'View inventory'),
    (UUID(), 'inventory', 'write', 'Update inventory'),
    (UUID(), 'reports', 'read', 'View reports'),
    (UUID(), 'admin', 'access', 'Admin panel access');

-- Assign permissions to roles
INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r CROSS JOIN permissions p
WHERE r.name = 'SUPER_ADMIN';

INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.name = 'ADMIN'
  AND p.resource IN ('users', 'products', 'orders', 'payments');

INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.name = 'SELLER'
  AND ((p.resource = 'products' AND p.action IN ('read', 'write'))
    OR (p.resource = 'orders' AND p.action IN ('read'))
    OR (p.resource = 'inventory' AND p.action IN ('read', 'write')));

INSERT IGNORE INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.name = 'BUYER'
  AND ((p.resource = 'products' AND p.action = 'read')
    OR (p.resource = 'orders' AND p.action IN ('read', 'write', 'cancel')));
