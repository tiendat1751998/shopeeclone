-- Payment Ledger Service - Initial Schema
-- Double-entry ledger, journal entries, and financial reconciliation

CREATE TABLE IF NOT EXISTS ledger_accounts (
    id VARCHAR(36) PRIMARY KEY,
    account_code VARCHAR(36) NOT NULL UNIQUE,
    account_name VARCHAR(255) NOT NULL,
    account_type ENUM('asset','liability','equity','revenue','expense') NOT NULL,
    parent_account_id VARCHAR(36) DEFAULT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'SGD',
    balance BIGINT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    description TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_la_type (account_type, is_active),
    INDEX idx_la_parent (parent_account_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS journal_entries (
    id VARCHAR(36) PRIMARY KEY,
    entry_number VARCHAR(36) NOT NULL UNIQUE,
    transaction_id VARCHAR(36) NOT NULL,
    reference_type ENUM('payment','refund','settlement','adjustment','fee','commission','withdrawal','deposit','transfer') NOT NULL,
    reference_id VARCHAR(36) NOT NULL,
    debit_account_id VARCHAR(36) NOT NULL,
    credit_account_id VARCHAR(36) NOT NULL,
    amount BIGINT NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'SGD',
    description TEXT DEFAULT NULL,
    status ENUM('pending','posted','reversed','failed') NOT NULL DEFAULT 'pending',
    posted_at TIMESTAMP NULL DEFAULT NULL,
    reversed_at TIMESTAMP NULL DEFAULT NULL,
    reversal_entry_id VARCHAR(36) DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_je_transaction (transaction_id),
    INDEX idx_je_reference (reference_type, reference_id),
    INDEX idx_je_debit (debit_account_id),
    INDEX idx_je_credit (credit_account_id),
    INDEX idx_je_status (status),
    INDEX idx_je_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS wallet_accounts (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    wallet_type ENUM('buyer','seller','platform','escrow','promotion') NOT NULL DEFAULT 'buyer',
    currency VARCHAR(3) NOT NULL DEFAULT 'SGD',
    balance BIGINT NOT NULL DEFAULT 0,
    frozen_balance BIGINT NOT NULL DEFAULT 0,
    status ENUM('active','frozen','suspended','closed') NOT NULL DEFAULT 'active',
    last_transaction_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_wa_user_type (user_id, wallet_type, currency),
    INDEX idx_wa_user (user_id),
    INDEX idx_wa_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS wallet_transactions (
    id VARCHAR(36) PRIMARY KEY,
    wallet_id VARCHAR(36) NOT NULL,
    transaction_type ENUM('credit','debit','freeze','unfreeze','adjustment') NOT NULL,
    amount BIGINT NOT NULL,
    balance_before BIGINT NOT NULL,
    balance_after BIGINT NOT NULL,
    reference_type ENUM('payment','refund','settlement','withdrawal','deposit','fee','commission','promotion','adjustment') NOT NULL,
    reference_id VARCHAR(36) DEFAULT NULL,
    journal_entry_id VARCHAR(36) DEFAULT NULL,
    description TEXT DEFAULT NULL,
    status ENUM('pending','completed','failed','reversed') NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_wt_wallet (wallet_id),
    INDEX idx_wt_reference (reference_type, reference_id),
    INDEX idx_wt_status (status),
    INDEX idx_wt_created (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS settlement_batches (
    id VARCHAR(36) PRIMARY KEY,
    batch_number VARCHAR(36) NOT NULL UNIQUE,
    seller_id VARCHAR(36) NOT NULL,
    period_start DATE NOT NULL,
    period_end DATE NOT NULL,
    total_sales BIGINT NOT NULL DEFAULT 0,
    total_fees BIGINT NOT NULL DEFAULT 0,
    total_commissions BIGINT NOT NULL DEFAULT 0,
    total_adjustments BIGINT NOT NULL DEFAULT 0,
    net_settlement BIGINT NOT NULL DEFAULT 0,
    currency VARCHAR(3) NOT NULL DEFAULT 'SGD',
    status ENUM('pending','processing','completed','failed','cancelled') NOT NULL DEFAULT 'pending',
    payment_method VARCHAR(64) DEFAULT NULL,
    payment_reference VARCHAR(255) DEFAULT NULL,
    paid_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_sb_seller (seller_id),
    INDEX idx_sb_status (status),
    INDEX idx_sb_period (period_start, period_end)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS reconciliation_records (
    id VARCHAR(36) PRIMARY KEY,
    reconciliation_date DATE NOT NULL,
    account_id VARCHAR(36) NOT NULL,
    ledger_balance BIGINT NOT NULL DEFAULT 0,
    external_balance BIGINT NOT NULL DEFAULT 0,
    difference BIGINT NOT NULL DEFAULT 0,
    status ENUM('matched','unmatched','investigating','resolved') NOT NULL DEFAULT 'unmatched',
    discrepancy_details JSON DEFAULT NULL,
    resolved_by VARCHAR(36) DEFAULT NULL,
    resolved_at TIMESTAMP NULL DEFAULT NULL,
    notes TEXT DEFAULT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_rr_date_account (reconciliation_date, account_id),
    INDEX idx_rr_status (status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
