-- FORK-002 correction: Omnio Pro benefits and free quota are level-scoped.
-- Legacy group-scoped tables remain intact so an older image can still roll back.

ALTER TABLE membership_level_group_benefits
    ADD COLUMN IF NOT EXISTS daily_free_usd DECIMAL(20, 10),
    ADD COLUMN IF NOT EXISTS monthly_free_usd DECIMAL(20, 10);

ALTER TABLE membership_level_group_benefits
    DROP CONSTRAINT IF EXISTS membership_level_group_benefits_daily_free_usd_check,
    DROP CONSTRAINT IF EXISTS membership_level_group_benefits_monthly_free_usd_check;

ALTER TABLE membership_level_group_benefits
    ADD CONSTRAINT membership_level_group_benefits_daily_free_usd_check
        CHECK (daily_free_usd IS NULL OR daily_free_usd >= 0),
    ADD CONSTRAINT membership_level_group_benefits_monthly_free_usd_check
        CHECK (monthly_free_usd IS NULL OR monthly_free_usd >= 0);

-- Preserve the currently effective group-scoped behavior on upgrade. Every
-- existing level starts with the same legacy values and can then be edited
-- independently. Existing explicit access remains enabled.
INSERT INTO membership_level_group_benefits (
    level_id,
    group_id,
    allow_access,
    rate_multiplier,
    rpm_limit,
    pro_only,
    daily_free_usd,
    monthly_free_usd
)
SELECT
    l.id,
    s.group_id,
    s.pro_only,
    s.rate_multiplier,
    NULL,
    s.pro_only,
    s.daily_free_usd,
    s.monthly_free_usd
FROM membership_levels l
CROSS JOIN omnio_pro_group_settings s
WHERE TRUE
ON CONFLICT (level_id, group_id) DO UPDATE SET
    allow_access = membership_level_group_benefits.allow_access OR EXCLUDED.pro_only,
    rate_multiplier = EXCLUDED.rate_multiplier,
    pro_only = EXCLUDED.pro_only,
    daily_free_usd = EXCLUDED.daily_free_usd,
    monthly_free_usd = EXCLUDED.monthly_free_usd,
    updated_at = NOW();

CREATE TABLE IF NOT EXISTS omnio_pro_level_quota_usage (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    level_id BIGINT NOT NULL REFERENCES membership_levels(id) ON DELETE CASCADE,
    group_id BIGINT NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    daily_window_start DATE NOT NULL,
    monthly_window_start DATE NOT NULL,
    daily_used_usd DECIMAL(20, 10) NOT NULL DEFAULT 0,
    monthly_used_usd DECIMAL(20, 10) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, level_id, group_id),
    CHECK (daily_used_usd >= 0),
    CHECK (monthly_used_usd >= 0)
);

CREATE INDEX IF NOT EXISTS idx_omnio_pro_level_quota_usage_group
    ON omnio_pro_level_quota_usage(level_id, group_id, updated_at DESC);

CREATE TABLE IF NOT EXISTS omnio_pro_level_quota_events (
    id BIGSERIAL PRIMARY KEY,
    request_id VARCHAR(255) NOT NULL,
    api_key_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    level_id BIGINT NOT NULL,
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

CREATE INDEX IF NOT EXISTS idx_omnio_pro_level_quota_events_user_level_group_time
    ON omnio_pro_level_quota_events(user_id, level_id, group_id, created_at DESC);

-- Seed each configured level with the legacy aggregate. This avoids granting a
-- fresh quota merely because the deployment introduced level-scoped counters.
INSERT INTO omnio_pro_level_quota_usage (
    user_id,
    level_id,
    group_id,
    daily_window_start,
    monthly_window_start,
    daily_used_usd,
    monthly_used_usd,
    created_at,
    updated_at
)
SELECT
    u.user_id,
    b.level_id,
    u.group_id,
    u.daily_window_start,
    u.monthly_window_start,
    u.daily_used_usd,
    u.monthly_used_usd,
    u.created_at,
    u.updated_at
FROM omnio_pro_quota_usage u
JOIN membership_level_group_benefits b ON b.group_id = u.group_id
WHERE TRUE
ON CONFLICT (user_id, level_id, group_id) DO NOTHING;

COMMENT ON COLUMN membership_level_group_benefits.daily_free_usd IS
    'Calendar-day free USD quota for this Omnio Pro level and group';
COMMENT ON COLUMN membership_level_group_benefits.monthly_free_usd IS
    'Calendar-month free USD quota for this Omnio Pro level and group';
COMMENT ON TABLE omnio_pro_level_quota_usage IS
    'FORK-002 level-scoped current daily and monthly Omnio Pro quota usage';
COMMENT ON TABLE omnio_pro_level_quota_events IS
    'FORK-002 immutable level-scoped split between Omnio Pro quota and wallet';
