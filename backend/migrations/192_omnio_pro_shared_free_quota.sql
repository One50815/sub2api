-- FORK-002 correction: each membership level owns one free-credit pool shared
-- by every eligible group. Legacy group-scoped columns and counters remain for
-- rollback compatibility, but the runtime moves to the shared table below.

ALTER TABLE membership_levels
    ADD COLUMN IF NOT EXISTS daily_free_usd DECIMAL(20, 10) NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS monthly_free_usd DECIMAL(20, 10) NOT NULL DEFAULT 0;

ALTER TABLE membership_levels
    DROP CONSTRAINT IF EXISTS membership_levels_daily_free_usd_check,
    DROP CONSTRAINT IF EXISTS membership_levels_monthly_free_usd_check;

ALTER TABLE membership_levels
    ADD CONSTRAINT membership_levels_daily_free_usd_check CHECK (daily_free_usd >= 0),
    ADD CONSTRAINT membership_levels_monthly_free_usd_check CHECK (monthly_free_usd >= 0);

-- A level previously had one quota setting per group. Promote the largest
-- configured value to the level so an upgrade never silently lowers a user's
-- available allowance when those group values differ.
WITH legacy_limits AS (
    SELECT
        level_id,
        MAX(COALESCE(daily_free_usd, 0)) AS daily_free_usd,
        MAX(COALESCE(monthly_free_usd, 0)) AS monthly_free_usd
    FROM membership_level_group_benefits
    GROUP BY level_id
)
UPDATE membership_levels AS levels
SET
    daily_free_usd = legacy_limits.daily_free_usd,
    monthly_free_usd = legacy_limits.monthly_free_usd,
    updated_at = NOW()
FROM legacy_limits
WHERE levels.id = legacy_limits.level_id;

CREATE TABLE IF NOT EXISTS omnio_pro_level_shared_quota_usage (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    level_id BIGINT NOT NULL REFERENCES membership_levels(id) ON DELETE CASCADE,
    daily_window_start DATE NOT NULL,
    monthly_window_start DATE NOT NULL,
    daily_used_usd DECIMAL(20, 10) NOT NULL DEFAULT 0,
    monthly_used_usd DECIMAL(20, 10) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (user_id, level_id),
    CHECK (daily_used_usd >= 0),
    CHECK (monthly_used_usd >= 0)
);

CREATE INDEX IF NOT EXISTS idx_omnio_pro_level_shared_quota_usage_level
    ON omnio_pro_level_shared_quota_usage(level_id, updated_at DESC);

-- Preserve already consumed credit. Current-day/current-month consumption is
-- summed across the old group rows into the single shared user-level pool.
WITH period AS (
    SELECT
        (CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Shanghai')::DATE AS day_start,
        date_trunc('month', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Shanghai')::DATE AS month_start
), legacy_usage AS (
    SELECT
        usage.user_id,
        usage.level_id,
        period.day_start,
        period.month_start,
        SUM(
            CASE WHEN usage.daily_window_start = period.day_start
                THEN usage.daily_used_usd ELSE 0 END
        ) AS daily_used_usd,
        SUM(
            CASE WHEN usage.monthly_window_start = period.month_start
                THEN usage.monthly_used_usd ELSE 0 END
        ) AS monthly_used_usd,
        MIN(usage.created_at) AS created_at,
        MAX(usage.updated_at) AS updated_at
    FROM omnio_pro_level_quota_usage AS usage
    CROSS JOIN period
    GROUP BY usage.user_id, usage.level_id, period.day_start, period.month_start
)
INSERT INTO omnio_pro_level_shared_quota_usage (
    user_id,
    level_id,
    daily_window_start,
    monthly_window_start,
    daily_used_usd,
    monthly_used_usd,
    created_at,
    updated_at
)
SELECT
    user_id,
    level_id,
    day_start,
    month_start,
    daily_used_usd,
    monthly_used_usd,
    created_at,
    updated_at
FROM legacy_usage
ON CONFLICT (user_id, level_id) DO UPDATE SET
    daily_window_start = EXCLUDED.daily_window_start,
    monthly_window_start = EXCLUDED.monthly_window_start,
    daily_used_usd = GREATEST(
        omnio_pro_level_shared_quota_usage.daily_used_usd,
        EXCLUDED.daily_used_usd
    ),
    monthly_used_usd = GREATEST(
        omnio_pro_level_shared_quota_usage.monthly_used_usd,
        EXCLUDED.monthly_used_usd
    ),
    updated_at = GREATEST(
        omnio_pro_level_shared_quota_usage.updated_at,
        EXCLUDED.updated_at
    );

COMMENT ON COLUMN membership_levels.daily_free_usd IS
    'Daily USD free-credit limit shared by all eligible groups for this level';
COMMENT ON COLUMN membership_levels.monthly_free_usd IS
    'Monthly USD free-credit limit shared by all eligible groups for this level';
COMMENT ON TABLE omnio_pro_level_shared_quota_usage IS
    'FORK-002 shared Omnio Pro quota usage keyed by user and membership level';
