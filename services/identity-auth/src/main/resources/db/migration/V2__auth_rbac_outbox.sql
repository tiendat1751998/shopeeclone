-- ============================================
-- RBAC: Roles & Permissions
-- ============================================
CREATE TABLE IF NOT EXISTS roles (
    role_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(255),
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS permissions (
    permission_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    resource VARCHAR(100) NOT NULL,
    action VARCHAR(50) NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(resource, action)
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    permission_id UUID NOT NULL REFERENCES permissions(permission_id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(role_id) ON DELETE CASCADE,
    assigned_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role_id)
);

CREATE INDEX idx_user_roles_user ON user_roles(user_id);
CREATE INDEX idx_role_permissions_role ON role_permissions(role_id);

-- ============================================
-- Outbox Pattern
-- ============================================
CREATE TABLE IF NOT EXISTS outbox_events (
    event_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type VARCHAR(100) NOT NULL,
    aggregate_id VARCHAR(100) NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN NOT NULL DEFAULT FALSE,
    processed_at TIMESTAMP WITH TIME ZONE,
    retry_count INT NOT NULL DEFAULT 0,
    last_error TEXT
);

CREATE INDEX idx_outbox_processed ON outbox_events(processed, created_at) WHERE processed = FALSE;

-- ============================================
-- Failed Login Attempts (Account Lockout)
-- ============================================
CREATE TABLE IF NOT EXISTS failed_login_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,
    attempted_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_failed_login_email ON failed_login_attempts(email, attempted_at);
CREATE INDEX idx_failed_login_ip ON failed_login_attempts(ip_address, attempted_at);

-- ============================================
-- Seed default roles and permissions
-- ============================================
INSERT INTO roles (name, description, is_system) VALUES
    ('SUPER_ADMIN', 'Full system access', TRUE),
    ('ADMIN', 'Administrative access', TRUE),
    ('SELLER', 'Seller/merchant account', TRUE),
    ('BUYER', 'Standard buyer account', TRUE)
ON CONFLICT (name) DO NOTHING;

INSERT INTO permissions (resource, action, description) VALUES
    ('users', 'read', 'View user profiles'),
    ('users', 'write', 'Create/update users'),
    ('users', 'delete', 'Delete users'),
    ('products', 'read', 'View products'),
    ('products', 'write', 'Create/update products'),
    ('products', 'delete', 'Delete products'),
    ('orders', 'read', 'View orders'),
    ('orders', 'write', 'Create/update orders'),
    ('orders', 'cancel', 'Cancel orders'),
    ('payments', 'read', 'View payments'),
    ('payments', 'refund', 'Process refunds'),
    ('inventory', 'read', 'View inventory'),
    ('inventory', 'write', 'Update inventory'),
    ('reports', 'read', 'View reports'),
    ('admin', 'access', 'Admin panel access')
ON CONFLICT (resource, action) DO NOTHING;

-- Assign permissions to roles
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r CROSS JOIN permissions p
WHERE r.name = 'SUPER_ADMIN';

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.name = 'ADMIN'
  AND p.resource IN ('users', 'products', 'orders', 'payments');

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.name = 'SELLER'
  AND ((p.resource = 'products' AND p.action IN ('read', 'write'))
    OR (p.resource = 'orders' AND p.action IN ('read'))
    OR (p.resource = 'inventory' AND p.action IN ('read', 'write')));

INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.name = 'BUYER'
  AND ((p.resource = 'products' AND p.action = 'read')
    OR (p.resource = 'orders' AND p.action IN ('read', 'write', 'cancel')));
