-- FORK-002: per-group Omnio Pro free quota with atomic wallet overflow billing.

ALTER TABLE omnio_pro_group_settings
    ADD COLUMN IF NOT EXISTS daily_free_usd DECIMAL(20, 10),
    ADD COLUMN IF NOT EXISTS monthly_free_usd DECIMAL(20, 10);

ALTER TABLE omnio_pro_group_settings
    DROP CONSTRAINT IF EXISTS omnio_pro_group_settings_daily_free_usd_check,
    DROP CONSTRAINT IF EXISTS omnio_pro_group_settings_monthly_free_usd_check;

ALTER TABLE omnio_pro_group_settings
    ADD CONSTRAINT omnio_pro_group_settings_daily_free_usd_check
        CHECK (daily_free_usd IS NULL OR daily_free_usd >= 0),
    ADD CONSTRAINT omnio_pro_group_settings_monthly_free_usd_check
        CHECK (monthly_free_usd IS NULL OR monthly_free_usd >= 0);

CREATE TABLE IF NOT EXISTS omnio_pro_quota_usage (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    daily_window_start DATE NOT NULL,
    monthly_window_start DATE NOT NULL,
    daily_used_usd DECIMAL(20, 10) NOT NULL DEFAULT 0,
    monthly_used_usd DECIMAL(20, 10) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, group_id),
    CHECK (daily_used_usd >= 0),
    CHECK (monthly_used_usd >= 0)
);

CREATE INDEX IF NOT EXISTS idx_omnio_pro_quota_usage_group
    ON omnio_pro_quota_usage(group_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS omnio_pro_quota_events (
    id BIGSERIAL PRIMARY KEY,
    request_id VARCHAR(255) NOT NULL,
    api_key_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    total_cost_usd DECIMAL(20, 10) NOT NULL,
    free_cost_usd DECIMAL(20, 10) NOT NULL,
    wallet_cost_usd DECIMAL(20, 10) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(request_id, api_key_id),
    CHECK (total_cost_usd >= 0),
    CHECK (free_cost_usd >= 0),
    CHECK (wallet_cost_usd >= 0),
    CHECK (ABS(total_cost_usd - free_cost_usd - wallet_cost_usd) < 0.0000001)
);

CREATE INDEX IF NOT EXISTS idx_omnio_pro_quota_events_user_group_time
    ON omnio_pro_quota_events(user_id, group_id, created_at DESC);

COMMENT ON TABLE omnio_pro_quota_usage IS
    'FORK-002 current calendar-day and calendar-month Omnio Pro free quota usage (Asia/Shanghai)';
COMMENT ON TABLE omnio_pro_quota_events IS
    'FORK-002 immutable per-request split between Omnio Pro free quota and wallet';
