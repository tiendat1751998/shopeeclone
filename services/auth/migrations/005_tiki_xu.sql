-- TikiXu loyalty points tables
CREATE TABLE IF NOT EXISTS tiki_xu_accounts (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL UNIQUE,
    balance BIGINT NOT NULL DEFAULT 0,
    lifetime_earned BIGINT NOT NULL DEFAULT 0,
    lifetime_spent BIGINT NOT NULL DEFAULT 0,
    tier ENUM('bronze','silver','gold','platinum','diamond') NOT NULL DEFAULT 'bronze',
    tier_points INT NOT NULL DEFAULT 0,
    status ENUM('active','frozen','closed') NOT NULL DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_txa_tier (tier),
    INDEX idx_txa_status (status),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS tiki_xu_transactions (
    id VARCHAR(36) PRIMARY KEY,
    account_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    type ENUM('earn','spend','refund','expire','adjust','signup_bonus','purchase','review','referral') NOT NULL,
    reference_type VARCHAR(50) DEFAULT NULL,
    reference_id VARCHAR(36) DEFAULT NULL,
    description VARCHAR(500) DEFAULT NULL,
    expires_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_txt_account (account_id),
    INDEX idx_txt_user (user_id, created_at),
    INDEX idx_txt_type (type),
    INDEX idx_txt_reference (reference_type, reference_id),
    FOREIGN KEY (account_id) REFERENCES tiki_xu_accounts(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS tiki_xu_tier_benefits (
    id VARCHAR(36) PRIMARY KEY,
    tier ENUM('bronze','silver','gold','platinum','diamond') NOT NULL UNIQUE,
    min_points INT NOT NULL DEFAULT 0,
    earn_rate DECIMAL(5,2) NOT NULL DEFAULT 1.00,
    free_shipping BOOLEAN NOT NULL DEFAULT FALSE,
    exclusive_deals BOOLEAN NOT NULL DEFAULT FALSE,
    priority_support BOOLEAN NOT NULL DEFAULT FALSE,
    birthday_gift BOOLEAN NOT NULL DEFAULT FALSE,
    description TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
