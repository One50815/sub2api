-- FORK-002: independent VIP membership and immutable subscription entitlement snapshots.

CREATE TABLE IF NOT EXISTS membership_levels (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(80) NOT NULL,
    slug VARCHAR(80) NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    rank INT NOT NULL DEFAULT 0,
    badge_color VARCHAR(32) NOT NULL DEFAULT '#111827',
    concurrency_bonus INT NOT NULL DEFAULT 0,
    priority_support BOOLEAN NOT NULL DEFAULT FALSE,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (rank >= 0),
    CHECK (concurrency_bonus >= 0)
);

CREATE TABLE IF NOT EXISTS membership_level_group_benefits (
    id BIGSERIAL PRIMARY KEY,
    level_id BIGINT NOT NULL REFERENCES membership_levels(id) ON DELETE CASCADE,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    allow_access BOOLEAN NOT NULL DEFAULT FALSE,
    rate_multiplier DECIMAL(12, 6),
    rpm_limit INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(level_id, group_id),
    CHECK (rate_multiplier IS NULL OR rate_multiplier >= 0),
    CHECK (rpm_limit IS NULL OR rpm_limit >= 0)
);

CREATE TABLE IF NOT EXISTS membership_offers (
    id BIGSERIAL PRIMARY KEY,
    level_id BIGINT NOT NULL REFERENCES membership_levels(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    price DECIMAL(20, 2) NOT NULL,
    original_price DECIMAL(20, 2),
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    duration_days INT NOT NULL,
    for_sale BOOLEAN NOT NULL DEFAULT TRUE,
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (price >= 0),
    CHECK (original_price IS NULL OR original_price >= 0),
    CHECK (duration_days > 0)
);

CREATE TABLE IF NOT EXISTS user_membership_grants (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    level_id BIGINT NOT NULL REFERENCES membership_levels(id) ON DELETE RESTRICT,
    source_type VARCHAR(30) NOT NULL,
    source_id VARCHAR(80) NOT NULL,
    starts_at TIMESTAMPTZ NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    notes TEXT NOT NULL DEFAULT '',
    revoked_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(source_type, source_id, level_id),
    CHECK (status IN ('active', 'revoked')),
    CHECK (expires_at > starts_at)
);

CREATE INDEX IF NOT EXISTS idx_membership_grants_user_active
    ON user_membership_grants(user_id, status, expires_at);

CREATE TABLE IF NOT EXISTS subscription_grants (
    id BIGSERIAL PRIMARY KEY,
    order_id BIGINT UNIQUE REFERENCES payment_orders(id) ON DELETE SET NULL,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    subscription_id BIGINT REFERENCES user_subscriptions(id) ON DELETE SET NULL,
    plan_id BIGINT REFERENCES subscription_plans(id) ON DELETE SET NULL,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE RESTRICT,
    plan_name VARCHAR(100) NOT NULL DEFAULT '',
    daily_limit_usd DECIMAL(20, 10),
    weekly_limit_usd DECIMAL(20, 10),
    monthly_limit_usd DECIMAL(20, 10),
    rate_multiplier DECIMAL(12, 6) NOT NULL DEFAULT 1,
    overage_policy VARCHAR(20) NOT NULL DEFAULT 'block',
    starts_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (daily_limit_usd IS NULL OR daily_limit_usd >= 0),
    CHECK (weekly_limit_usd IS NULL OR weekly_limit_usd >= 0),
    CHECK (monthly_limit_usd IS NULL OR monthly_limit_usd >= 0),
    CHECK (rate_multiplier >= 0),
    CHECK (overage_policy IN ('block', 'wallet')),
    CHECK (status IN ('pending', 'active', 'revoked', 'expired'))
);

CREATE INDEX IF NOT EXISTS idx_subscription_grants_active
    ON subscription_grants(user_id, group_id, status, expires_at);

CREATE TABLE IF NOT EXISTS subscription_plan_membership_benefits (
    plan_id BIGINT PRIMARY KEY REFERENCES subscription_plans(id) ON DELETE CASCADE,
    level_id BIGINT NOT NULL REFERENCES membership_levels(id) ON DELETE CASCADE,
    duration_days INT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (duration_days IS NULL OR duration_days > 0)
);

CREATE TABLE IF NOT EXISTS subscription_change_schedules (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    current_subscription_id BIGINT REFERENCES user_subscriptions(id) ON DELETE CASCADE,
    target_plan_id BIGINT NOT NULL REFERENCES subscription_plans(id) ON DELETE CASCADE,
    effective_at TIMESTAMPTZ NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'scheduled',
    created_by BIGINT REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CHECK (status IN ('scheduled', 'applied', 'cancelled'))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_subscription_change_one_pending
    ON subscription_change_schedules(user_id, current_subscription_id)
    WHERE status = 'scheduled';

CREATE TABLE IF NOT EXISTS membership_audit_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    level_id BIGINT REFERENCES membership_levels(id) ON DELETE SET NULL,
    action VARCHAR(60) NOT NULL,
    source_type VARCHAR(30) NOT NULL DEFAULT '',
    source_id VARCHAR(80) NOT NULL DEFAULT '',
    operator VARCHAR(80) NOT NULL DEFAULT 'system',
    detail JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_membership_audit_user_created
    ON membership_audit_logs(user_id, created_at DESC);

COMMENT ON TABLE membership_levels IS 'FORK-002 VIP membership levels';
COMMENT ON TABLE subscription_grants IS 'FORK-002 immutable subscription entitlement snapshots';
