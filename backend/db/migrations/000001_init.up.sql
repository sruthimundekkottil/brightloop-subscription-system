-- Brightloop Subscription Management System — initial schema
-- Target: MySQL 8.0.16+ (CHECK constraints enforced; SKIP LOCKED available since 8.0.1)
-- Validated on MySQL 8.0.46. Apply with golang-migrate (DSN needs multiStatements=true).
-- All DATETIME values are UTC (application writes UTC; DSN sets loc=UTC).

-- =====================================================================
-- A — Platform & Access
-- =====================================================================

CREATE TABLE users (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    email         VARCHAR(254) NOT NULL COLLATE utf8mb4_0900_ai_ci, -- case-insensitive collation → UNIQUE is case-insensitive
    password_hash VARCHAR(255) NOT NULL,                              -- bcrypt
    role          ENUM('admin', 'staff', 'customer') NOT NULL,
    is_active     BOOLEAN NOT NULL DEFAULT TRUE,
    created_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uq_users_email (email)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE customers (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id         BIGINT UNSIGNED NULL,                 -- NULL if staff created the customer without a login
    name            VARCHAR(150) NOT NULL,
    company_name    VARCHAR(150) NULL,
    phone           VARCHAR(30) NULL,
    billing_address TEXT NULL,
    created_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uq_customers_user (user_id),
    KEY idx_customers_name (name),
    CONSTRAINT fk_customers_user FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE staff_profiles (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id    BIGINT UNSIGNED NOT NULL,
    full_name  VARCHAR(150) NOT NULL,
    job_title  VARCHAR(100) NULL,
    created_by BIGINT UNSIGNED NOT NULL,                  -- must be an admin (checked in the service layer)
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uq_staff_user (user_id),
    CONSTRAINT fk_staff_user       FOREIGN KEY (user_id)    REFERENCES users (id),
    CONSTRAINT fk_staff_created_by FOREIGN KEY (created_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE password_reset_tokens (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    user_id    BIGINT UNSIGNED NOT NULL,
    token_hash CHAR(64) NOT NULL,                         -- SHA-256 hex of the emailed token
    expires_at DATETIME(3) NOT NULL,
    used_at    DATETIME(3) NULL,                          -- single use
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uq_reset_token_hash (token_hash),
    KEY idx_reset_tokens_user (user_id),
    CONSTRAINT fk_reset_user FOREIGN KEY (user_id) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE company_settings (
    id                              TINYINT UNSIGNED NOT NULL DEFAULT 1,
    company_name                    VARCHAR(150) NOT NULL DEFAULT 'Brightloop Software',
    address                         TEXT NULL,
    currency_code                   CHAR(3) NOT NULL DEFAULT 'INR',
    invoice_prefix                  VARCHAR(10) NOT NULL DEFAULT 'INV',
    default_quotation_validity_days INT UNSIGNED NOT NULL DEFAULT 30,
    auto_confirm_recurring_invoices BOOLEAN NOT NULL DEFAULT FALSE,
    test_clock_date                 DATE NULL,           -- honoured only when APP_ENV=development; forward-only
    updated_at                      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    CONSTRAINT chk_company_settings_single_row CHECK (id = 1)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Gapless document numbers (MySQL has no sequences; a locked counter row also rolls back cleanly)
CREATE TABLE number_sequences (
    name       VARCHAR(40) NOT NULL,                      -- e.g. 'subscription-2026', 'invoice-2026'
    next_value BIGINT UNSIGNED NOT NULL DEFAULT 1,
    PRIMARY KEY (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- DB-backed job queue / transactional outbox (Kill Bill notifications pattern)
CREATE TABLE job_queue (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    queue_name   VARCHAR(64) NOT NULL,                    -- 'email', 'overdue-reminder'
    dedupe_key   VARCHAR(150) NULL,                       -- e.g. 'invoice-sent:42'; NULLs allowed many times
    payload      JSON NOT NULL,
    effective_at DATETIME(3) NOT NULL,                    -- run no earlier than this
    state        ENUM('available', 'processing', 'done', 'failed') NOT NULL DEFAULT 'available',
    attempts     INT UNSIGNED NOT NULL DEFAULT 0,
    max_attempts INT UNSIGNED NOT NULL DEFAULT 5,
    locked_by    VARCHAR(64) NULL,
    locked_until DATETIME(3) NULL,
    last_error   TEXT NULL,
    created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    processed_at DATETIME(3) NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_jobs_dedupe (dedupe_key),
    KEY idx_jobs_ready (state, effective_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- Audit trail (Kill Bill audit_log / UniBee merchant_operation_log)
CREATE TABLE audit_log (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    table_name VARCHAR(64) NOT NULL,
    record_id  BIGINT UNSIGNED NOT NULL,
    action     VARCHAR(40) NOT NULL,                      -- 'create', 'update', 'confirm', 'cancel', ...
    changed_by BIGINT UNSIGNED NULL,                      -- NULL = system job
    reason     VARCHAR(255) NULL,
    changes    JSON NULL,                                 -- {"field": [old, new]}
    created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_audit_target (table_name, record_id),
    KEY idx_audit_user (changed_by),
    CONSTRAINT fk_audit_user FOREIGN KEY (changed_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- =====================================================================
-- B — Catalog & Pricing
-- =====================================================================

CREATE TABLE products (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name         VARCHAR(150) NOT NULL,
    product_type VARCHAR(50) NOT NULL,                    -- 'software', 'addon', 'service'
    sales_price  DECIMAL(12,2) NOT NULL,
    cost_price   DECIMAL(12,2) NOT NULL,
    notes        TEXT NULL,
    is_recurring BOOLEAN NOT NULL DEFAULT TRUE,           -- FALSE = one-time item, billed on the first invoice only
    is_active    BOOLEAN NOT NULL DEFAULT TRUE,           -- archive, never delete
    created_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at   DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_products_active_name (is_active, name),
    CONSTRAINT chk_products_sales_price CHECK (sales_price >= 0),
    CONSTRAINT chk_products_cost_price  CHECK (cost_price >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE product_attributes (
    id   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(80) NOT NULL,                            -- 'Brand'
    PRIMARY KEY (id),
    UNIQUE KEY uq_attribute_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE attribute_values (
    id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    attribute_id BIGINT UNSIGNED NOT NULL,
    value        VARCHAR(80) NOT NULL,                    -- 'Odoo'
    PRIMARY KEY (id),
    UNIQUE KEY uq_attribute_value (attribute_id, value),
    CONSTRAINT fk_attr_values_attribute FOREIGN KEY (attribute_id) REFERENCES product_attributes (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- price_extra per product × value (as Odoo's product.template.attribute.value)
CREATE TABLE product_attribute_values (
    id                 BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    product_id         BIGINT UNSIGNED NOT NULL,
    attribute_value_id BIGINT UNSIGNED NOT NULL,
    price_extra        DECIMAL(12,2) NOT NULL DEFAULT 0, -- Brand = Odoo → 560.00 (may be negative; final price checked ≥ 0 in service)
    is_active          BOOLEAN NOT NULL DEFAULT TRUE,     -- archived, not deleted, once used
    PRIMARY KEY (id),
    UNIQUE KEY uq_product_attr_value (product_id, attribute_value_id),
    CONSTRAINT fk_pav_product FOREIGN KEY (product_id)         REFERENCES products (id) ON DELETE CASCADE,
    CONSTRAINT fk_pav_value   FOREIGN KEY (attribute_value_id) REFERENCES attribute_values (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE product_variants (
    id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    product_id BIGINT UNSIGNED NOT NULL,
    name       VARCHAR(150) NOT NULL,                     -- 'CRM Suite (Odoo)'
    is_active  BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (id),
    KEY idx_variants_product (product_id),
    CONSTRAINT fk_variants_product FOREIGN KEY (product_id) REFERENCES products (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE product_variant_values (
    variant_id                 BIGINT UNSIGNED NOT NULL,
    product_attribute_value_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (variant_id, product_attribute_value_id),
    CONSTRAINT fk_pvv_variant FOREIGN KEY (variant_id)                 REFERENCES product_variants (id)         ON DELETE CASCADE,
    CONSTRAINT fk_pvv_pav     FOREIGN KEY (product_attribute_value_id) REFERENCES product_attribute_values (id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE plans (
    id                     BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name                   VARCHAR(150) NOT NULL,
    price                  DECIMAL(12,2) NOT NULL,
    billing_period         ENUM('daily', 'weekly', 'monthly', 'yearly') NOT NULL,
    billing_interval_count SMALLINT UNSIGNED NOT NULL DEFAULT 1,  -- e.g. 3 × monthly = quarterly (UniBee interval_count)
    min_quantity           INT UNSIGNED NOT NULL DEFAULT 1,
    available_from         DATE NOT NULL,
    available_until        DATE NULL,
    auto_close             BOOLEAN NOT NULL DEFAULT FALSE,
    can_close_early        BOOLEAN NOT NULL DEFAULT TRUE,
    can_pause              BOOLEAN NOT NULL DEFAULT FALSE,
    can_renew              BOOLEAN NOT NULL DEFAULT TRUE,
    is_active              BOOLEAN NOT NULL DEFAULT TRUE,
    created_at             DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at             DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    CONSTRAINT chk_plans_price          CHECK (price >= 0),
    CONSTRAINT chk_plans_interval_count CHECK (billing_interval_count >= 1),
    CONSTRAINT chk_plans_min_quantity   CHECK (min_quantity >= 1),
    CONSTRAINT chk_plans_availability   CHECK (available_until IS NULL OR available_until >= available_from)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE tax_rules (
    id             BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name           VARCHAR(80) NOT NULL,                  -- 'GST 18%'
    rate_percent   DECIMAL(5,2) NOT NULL,
    is_default     BOOLEAN NOT NULL DEFAULT FALSE,
    is_active      BOOLEAN NOT NULL DEFAULT TRUE,
    default_marker TINYINT UNSIGNED GENERATED ALWAYS AS (IF(is_default, 1, NULL)) STORED, -- MySQL has no partial indexes
    PRIMARY KEY (id),
    UNIQUE KEY uq_tax_name (name),
    UNIQUE KEY uq_tax_single_default (default_marker),
    CONSTRAINT chk_tax_rate CHECK (rate_percent >= 0 AND rate_percent <= 100)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE product_tax_rules (
    product_id  BIGINT UNSIGNED NOT NULL,
    tax_rule_id BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (product_id, tax_rule_id),
    CONSTRAINT fk_ptr_product FOREIGN KEY (product_id)  REFERENCES products (id) ON DELETE CASCADE,
    CONSTRAINT fk_ptr_tax     FOREIGN KEY (tax_rule_id) REFERENCES tax_rules (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE discount_rules (
    id                          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name                        VARCHAR(120) NOT NULL,
    discount_type               ENUM('fixed', 'percentage') NOT NULL,
    value                       DECIMAL(12,2) NOT NULL,
    scope                       ENUM('subscription', 'products') NOT NULL,
    customer_id                 BIGINT UNSIGNED NULL,     -- NULL = any customer
    min_purchase_amount         DECIMAL(12,2) NULL,
    min_quantity                INT UNSIGNED NULL,
    starts_on                   DATE NOT NULL,
    ends_on                     DATE NOT NULL,
    max_uses                    INT UNSIGNED NULL,        -- total uses; NULL = unlimited
    max_uses_per_customer       INT UNSIGNED NULL,        -- UniBee user_limit
    max_cycles_per_subscription INT UNSIGNED NULL,        -- UniBee cycle_limit: first N billing cycles only
    max_discount_amount         DECIMAL(12,2) NULL,       -- cap per invoice (Odoo discount_max_amount)
    times_used                  INT UNSIGNED NOT NULL DEFAULT 0,  -- incremented/decremented by D in the billing transaction
    is_active                   BOOLEAN NOT NULL DEFAULT TRUE,
    created_by                  BIGINT UNSIGNED NOT NULL,
    created_at                  DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_discounts_active_dates (is_active, starts_on, ends_on),
    CONSTRAINT fk_discount_customer   FOREIGN KEY (customer_id) REFERENCES customers (id),
    CONSTRAINT fk_discount_created_by FOREIGN KEY (created_by)  REFERENCES users (id),
    CONSTRAINT chk_discount_value      CHECK (value > 0),
    CONSTRAINT chk_discount_percentage CHECK (discount_type <> 'percentage' OR value <= 100),
    CONSTRAINT chk_discount_dates      CHECK (ends_on >= starts_on),
    CONSTRAINT chk_discount_min_amount CHECK (min_purchase_amount IS NULL OR min_purchase_amount >= 0),
    CONSTRAINT chk_discount_min_qty    CHECK (min_quantity IS NULL OR min_quantity >= 1),
    CONSTRAINT chk_discount_max_uses   CHECK (max_uses IS NULL OR max_uses >= 1),
    CONSTRAINT chk_discount_usage      CHECK (max_uses IS NULL OR times_used <= max_uses),
    CONSTRAINT chk_discount_cap        CHECK (max_discount_amount IS NULL OR max_discount_amount > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE discount_rule_products (
    discount_rule_id BIGINT UNSIGNED NOT NULL,
    product_id       BIGINT UNSIGNED NOT NULL,
    PRIMARY KEY (discount_rule_id, product_id),
    CONSTRAINT fk_drp_rule    FOREIGN KEY (discount_rule_id) REFERENCES discount_rules (id) ON DELETE CASCADE,
    CONSTRAINT fk_drp_product FOREIGN KEY (product_id)       REFERENCES products (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- =====================================================================
-- C — Sales & Subscriptions
-- =====================================================================

CREATE TABLE payment_terms (
    id       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name     VARCHAR(60) NOT NULL,                        -- 'Immediate', 'Net 15', 'Net 30'
    due_days INT UNSIGNED NOT NULL,
    PRIMARY KEY (id),
    UNIQUE KEY uq_payment_terms_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE close_reasons (
    id       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name     VARCHAR(80) NOT NULL,                        -- 'Too expensive', 'Renewed', 'Ended', 'Quotation expired'
    is_churn BOOLEAN NOT NULL DEFAULT TRUE,
    PRIMARY KEY (id),
    UNIQUE KEY uq_close_reasons_name (name)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE quotation_templates (
    id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    name             VARCHAR(150) NOT NULL,
    plan_id          BIGINT UNSIGNED NOT NULL,
    validity_days    INT UNSIGNED NOT NULL,
    payment_terms_id BIGINT UNSIGNED NULL,
    terms            TEXT NULL,
    notes            TEXT NULL,
    is_active        BOOLEAN NOT NULL DEFAULT TRUE,
    created_by       BIGINT UNSIGNED NOT NULL,
    created_at       DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    CONSTRAINT fk_templates_plan       FOREIGN KEY (plan_id)          REFERENCES plans (id),
    CONSTRAINT fk_templates_terms      FOREIGN KEY (payment_terms_id) REFERENCES payment_terms (id),
    CONSTRAINT fk_templates_created_by FOREIGN KEY (created_by)       REFERENCES users (id),
    CONSTRAINT chk_templates_validity  CHECK (validity_days > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE quotation_template_lines (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    template_id BIGINT UNSIGNED NOT NULL,
    product_id  BIGINT UNSIGNED NOT NULL,
    variant_id  BIGINT UNSIGNED NULL,
    quantity    INT UNSIGNED NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_tlines_template FOREIGN KEY (template_id) REFERENCES quotation_templates (id) ON DELETE CASCADE,
    CONSTRAINT fk_tlines_product  FOREIGN KEY (product_id)  REFERENCES products (id),
    CONSTRAINT fk_tlines_variant  FOREIGN KEY (variant_id)  REFERENCES product_variants (id),
    CONSTRAINT chk_tlines_quantity CHECK (quantity > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE subscriptions (
    id                    BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    number                VARCHAR(30) NOT NULL,           -- 'SUB-2026-00042' from number_sequences
    customer_id           BIGINT UNSIGNED NOT NULL,
    plan_id               BIGINT UNSIGNED NOT NULL,
    template_id           BIGINT UNSIGNED NULL,
    payment_terms_id      BIGINT UNSIGNED NOT NULL,
    status                ENUM('draft', 'quotation', 'confirmed', 'active', 'paused', 'closed') NOT NULL DEFAULT 'draft',
    start_date            DATE NOT NULL,
    expiry_date           DATE NULL,                      -- NULL = open-ended
    billing_anchor_date   DATE NULL,                      -- periods are computed from here (Kill Bill BCD / UniBee billing_cycle_anchor)
    next_billing_date     DATE NULL,                      -- read by D's billing engine
    quotation_valid_until DATE NULL,
    scheduled_close_date  DATE NULL,                      -- close at end of term (Kill Bill END_OF_TERM / UniBee cancel_at_period_end)
    recurring_total       DECIMAL(12,2) NOT NULL DEFAULT 0, -- per billing period, before tax and discounts
    recurring_monthly     DECIMAL(12,2) NOT NULL DEFAULT 0, -- normalised to one month → MRR
    renewed_from_id       BIGINT UNSIGNED NULL,
    terms                 TEXT NULL,
    salesperson_id        BIGINT UNSIGNED NULL,
    paused_at             DATETIME(3) NULL,
    closed_at             DATETIME(3) NULL,
    close_reason_id       BIGINT UNSIGNED NULL,
    close_note            VARCHAR(200) NULL,
    created_at            DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at            DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    open_renewal_marker   BIGINT UNSIGNED GENERATED ALWAYS AS
                          (IF(renewed_from_id IS NOT NULL AND status <> 'closed', renewed_from_id, NULL)) STORED,
    PRIMARY KEY (id),
    UNIQUE KEY uq_subscriptions_number (number),
    UNIQUE KEY uq_subscriptions_one_open_renewal (open_renewal_marker),
    KEY idx_subs_status_next_billing (status, next_billing_date),
    KEY idx_subs_status_expiry (status, expiry_date),
    KEY idx_subs_status_scheduled_close (status, scheduled_close_date),
    CONSTRAINT fk_subs_customer     FOREIGN KEY (customer_id)      REFERENCES customers (id),
    CONSTRAINT fk_subs_plan         FOREIGN KEY (plan_id)          REFERENCES plans (id),
    CONSTRAINT fk_subs_template     FOREIGN KEY (template_id)      REFERENCES quotation_templates (id),
    CONSTRAINT fk_subs_terms        FOREIGN KEY (payment_terms_id) REFERENCES payment_terms (id),
    CONSTRAINT fk_subs_renewed_from FOREIGN KEY (renewed_from_id)  REFERENCES subscriptions (id),
    CONSTRAINT fk_subs_salesperson  FOREIGN KEY (salesperson_id)   REFERENCES users (id),
    CONSTRAINT fk_subs_close_reason FOREIGN KEY (close_reason_id)  REFERENCES close_reasons (id),
    CONSTRAINT chk_subs_expiry      CHECK (expiry_date IS NULL OR expiry_date >= start_date),
    CONSTRAINT chk_subs_recurring   CHECK (recurring_total >= 0 AND recurring_monthly >= 0),
    CONSTRAINT chk_subs_closed      CHECK ((status = 'closed') = (closed_at IS NOT NULL))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE subscription_lines (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    subscription_id BIGINT UNSIGNED NOT NULL,
    product_id      BIGINT UNSIGNED NOT NULL,
    variant_id      BIGINT UNSIGNED NULL,
    quantity        INT UNSIGNED NOT NULL,
    unit_price      DECIMAL(12,2) NOT NULL,               -- agreed price, from B's GetUnitPrice
    amount          DECIMAL(14,2) GENERATED ALWAYS AS (quantity * unit_price) STORED,
    PRIMARY KEY (id),
    KEY idx_sub_lines_subscription (subscription_id),
    CONSTRAINT fk_sub_lines_subscription FOREIGN KEY (subscription_id) REFERENCES subscriptions (id) ON DELETE CASCADE,
    CONSTRAINT fk_sub_lines_product      FOREIGN KEY (product_id)      REFERENCES products (id),
    CONSTRAINT fk_sub_lines_variant      FOREIGN KEY (variant_id)      REFERENCES product_variants (id),
    CONSTRAINT chk_sub_lines_quantity    CHECK (quantity > 0),
    CONSTRAINT chk_sub_lines_price       CHECK (unit_price >= 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE subscription_status_history (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    subscription_id BIGINT UNSIGNED NOT NULL,
    from_status     ENUM('draft', 'quotation', 'confirmed', 'active', 'paused', 'closed') NULL,
    to_status       ENUM('draft', 'quotation', 'confirmed', 'active', 'paused', 'closed') NOT NULL,
    changed_by      BIGINT UNSIGNED NULL,                 -- NULL = scheduled job
    note            VARCHAR(255) NULL,
    changed_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_sub_history (subscription_id, changed_at),
    CONSTRAINT fk_sub_history_subscription FOREIGN KEY (subscription_id) REFERENCES subscriptions (id) ON DELETE CASCADE,
    CONSTRAINT fk_sub_history_user         FOREIGN KEY (changed_by)      REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- =====================================================================
-- D — Billing, Payments & Reporting
-- =====================================================================

CREATE TABLE invoices (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    number          VARCHAR(30) NULL,                     -- NULL while draft; assigned on confirm
    customer_id     BIGINT UNSIGNED NOT NULL,
    subscription_id BIGINT UNSIGNED NULL,
    status          ENUM('draft', 'confirmed', 'paid', 'cancelled') NOT NULL DEFAULT 'draft',
    invoice_date    DATE NOT NULL,
    due_date        DATE NOT NULL,
    period_start    DATE NULL,
    period_end      DATE NULL,
    billing_name    VARCHAR(150) NOT NULL,                -- customer snapshot
    billing_address TEXT NULL,
    subtotal        DECIMAL(12,2) NOT NULL DEFAULT 0,     -- recurring + fixed items
    discount_total  DECIMAL(12,2) NOT NULL DEFAULT 0,     -- positive number (sum of negative discount items, sign flipped)
    tax_total       DECIMAL(12,2) NOT NULL DEFAULT 0,
    total           DECIMAL(12,2) NOT NULL DEFAULT 0,     -- = SUM(invoice_items.amount)
    amount_paid     DECIMAL(12,2) NOT NULL DEFAULT 0,     -- maintained in the payment transaction
    payment_status  VARCHAR(10) GENERATED ALWAYS AS (
                        CASE WHEN amount_paid = 0 THEN 'not_paid'
                             WHEN amount_paid < total THEN 'partial'
                             ELSE 'paid' END) STORED,
    confirmed_at    DATETIME(3) NULL,
    sent_at         DATETIME(3) NULL,
    cancelled_at    DATETIME(3) NULL,
    cancel_reason   VARCHAR(200) NULL,
    created_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    updated_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uq_invoices_number (number),               -- multiple NULLs allowed (drafts)
    KEY idx_invoices_status_due (status, due_date),
    KEY idx_invoices_date (invoice_date),
    CONSTRAINT fk_invoices_customer     FOREIGN KEY (customer_id)     REFERENCES customers (id),
    CONSTRAINT fk_invoices_subscription FOREIGN KEY (subscription_id) REFERENCES subscriptions (id),
    CONSTRAINT chk_invoices_total       CHECK (total >= 0),
    CONSTRAINT chk_invoices_paid_range  CHECK (amount_paid >= 0 AND amount_paid <= total),
    CONSTRAINT chk_invoices_due         CHECK (due_date >= invoice_date),
    CONSTRAINT chk_invoices_numbered    CHECK (status IN ('draft', 'cancelled') OR number IS NOT NULL),
    CONSTRAINT chk_invoices_paid_status CHECK (status <> 'paid' OR amount_paid = total),
    CONSTRAINT chk_invoices_period      CHECK (period_end IS NULL OR period_end >= period_start)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- One table for every amount on an invoice (Kill Bill invoice_items pattern).
-- Items are append-only and invoices are never deleted (drafts are cancelled), so no ON DELETE CASCADE here.
CREATE TABLE invoice_items (
    id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    invoice_id       BIGINT UNSIGNED NOT NULL,
    item_type        ENUM('recurring', 'fixed', 'discount', 'tax') NOT NULL,
    linked_item_id   BIGINT UNSIGNED NULL,                -- discount/tax items point at the item they apply to
    product_id       BIGINT UNSIGNED NULL,                -- NULL for the plan line, discounts and taxes
    plan_id          BIGINT UNSIGNED NULL,
    discount_rule_id BIGINT UNSIGNED NULL,
    tax_rule_id      BIGINT UNSIGNED NULL,
    description      VARCHAR(255) NOT NULL,               -- snapshot, e.g. 'CRM Suite (Odoo) · 1 Mar – 31 Mar 2026'
    service_start    DATE NULL,                           -- service period of recurring items
    service_end      DATE NULL,
    quantity         INT UNSIGNED NOT NULL DEFAULT 1,
    unit_price       DECIMAL(12,2) NULL,                  -- snapshot for recurring/fixed items
    rate_percent     DECIMAL(5,2) NULL,                   -- snapshot: tax rate or discount percentage
    amount           DECIMAL(12,2) NOT NULL,              -- negative only for discount items
    created_at       DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    KEY idx_items_invoice (invoice_id),
    KEY idx_items_linked (linked_item_id),
    CONSTRAINT fk_items_invoice  FOREIGN KEY (invoice_id)       REFERENCES invoices (id),   -- invoices are never deleted
    CONSTRAINT fk_items_linked   FOREIGN KEY (linked_item_id)   REFERENCES invoice_items (id),
    CONSTRAINT fk_items_product  FOREIGN KEY (product_id)       REFERENCES products (id),
    CONSTRAINT fk_items_plan     FOREIGN KEY (plan_id)          REFERENCES plans (id),
    CONSTRAINT fk_items_discount FOREIGN KEY (discount_rule_id) REFERENCES discount_rules (id),
    CONSTRAINT fk_items_tax      FOREIGN KEY (tax_rule_id)      REFERENCES tax_rules (id),
    CONSTRAINT chk_items_link     CHECK ((item_type IN ('discount', 'tax')) = (linked_item_id IS NOT NULL)),
    CONSTRAINT chk_items_discount CHECK ((item_type = 'discount') = (discount_rule_id IS NOT NULL)),
    CONSTRAINT chk_items_tax      CHECK ((item_type = 'tax') = (tax_rule_id IS NOT NULL)),
    CONSTRAINT chk_items_sign     CHECK ((item_type = 'discount' AND amount < 0) OR (item_type <> 'discount' AND amount >= 0)),
    CONSTRAINT chk_items_quantity CHECK (quantity > 0),
    CONSTRAINT chk_items_service  CHECK (service_end IS NULL OR service_end >= service_start)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- The duplicate-billing guard: one invoice per subscription per period
CREATE TABLE billing_records (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    subscription_id BIGINT UNSIGNED NOT NULL,
    period_start    DATE NOT NULL,
    period_end      DATE NOT NULL,
    invoice_id      BIGINT UNSIGNED NOT NULL,
    inputs          JSON NULL,                            -- snapshot of what billing used (Kill Bill invoice_billing_events)
    created_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uq_billing_period (subscription_id, period_start),
    UNIQUE KEY uq_billing_invoice (invoice_id),
    CONSTRAINT fk_billing_subscription FOREIGN KEY (subscription_id) REFERENCES subscriptions (id),
    CONSTRAINT fk_billing_invoice      FOREIGN KEY (invoice_id)      REFERENCES invoices (id),
    CONSTRAINT chk_billing_period      CHECK (period_end >= period_start)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE discount_applications (
    id               BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    discount_rule_id BIGINT UNSIGNED NOT NULL,
    invoice_id       BIGINT UNSIGNED NOT NULL,
    customer_id      BIGINT UNSIGNED NOT NULL,
    subscription_id  BIGINT UNSIGNED NULL,
    amount           DECIMAL(12,2) NOT NULL,              -- positive: total discount granted on this invoice
    reversed_at      DATETIME(3) NULL,                    -- set when the invoice is cancelled (UniBee 'rollback')
    applied_at       DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uq_discount_once_per_invoice (discount_rule_id, invoice_id),
    KEY idx_da_rule_customer (discount_rule_id, customer_id),
    KEY idx_da_rule_subscription (discount_rule_id, subscription_id),
    CONSTRAINT fk_da_rule         FOREIGN KEY (discount_rule_id) REFERENCES discount_rules (id),
    CONSTRAINT fk_da_invoice      FOREIGN KEY (invoice_id)       REFERENCES invoices (id),
    CONSTRAINT fk_da_customer     FOREIGN KEY (customer_id)      REFERENCES customers (id),
    CONSTRAINT fk_da_subscription FOREIGN KEY (subscription_id)  REFERENCES subscriptions (id),
    CONSTRAINT chk_da_amount      CHECK (amount > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE payments (
    id              BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    invoice_id      BIGINT UNSIGNED NOT NULL,
    method          ENUM('bank_transfer', 'card', 'cash', 'cheque', 'upi', 'other') NOT NULL,
    amount          DECIMAL(12,2) NOT NULL,
    payment_date    DATE NOT NULL,
    note            VARCHAR(250) NULL,
    idempotency_key CHAR(36) NOT NULL,                    -- UUID generated by the form; a double submit can't record twice
    recorded_by     BIGINT UNSIGNED NOT NULL,
    created_at      DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    PRIMARY KEY (id),
    UNIQUE KEY uq_payments_idempotency (idempotency_key),
    KEY idx_payments_invoice (invoice_id),
    KEY idx_payments_date (payment_date),
    CONSTRAINT fk_payments_invoice FOREIGN KEY (invoice_id)  REFERENCES invoices (id),
    CONSTRAINT fk_payments_user    FOREIGN KEY (recorded_by) REFERENCES users (id),
    CONSTRAINT chk_payments_amount CHECK (amount > 0)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

CREATE TABLE billing_runs (
    id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    as_of_date    DATE NOT NULL,
    triggered_by  BIGINT UNSIGNED NULL,                   -- NULL = scheduler
    dry_run       BOOLEAN NOT NULL DEFAULT FALSE,
    started_at    DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3),
    finished_at   DATETIME(3) NULL,
    due_count     INT UNSIGNED NOT NULL DEFAULT 0,
    created_count INT UNSIGNED NOT NULL DEFAULT 0,
    skipped_count INT UNSIGNED NOT NULL DEFAULT 0,
    failed_count  INT UNSIGNED NOT NULL DEFAULT 0,
    errors        JSON NULL,                              -- [{"subscription_id": .., "message": ..}]
    PRIMARY KEY (id),
    KEY idx_billing_runs_started (started_at),
    CONSTRAINT fk_billing_runs_user FOREIGN KEY (triggered_by) REFERENCES users (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
