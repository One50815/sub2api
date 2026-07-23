-- FORK-002 correction: Omnio Pro pricing and visibility belong to each group.
-- Membership levels only determine whether a user currently has Pro access.

CREATE TABLE IF NOT EXISTS omnio_pro_group_settings (
    group_id BIGINT PRIMARY KEY REFERENCES groups(id) ON DELETE CASCADE,
    rate_multiplier DECIMAL(12, 6),
    pro_only BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT omnio_pro_group_rate_non_negative
        CHECK (rate_multiplier IS NULL OR rate_multiplier >= 0)
);

-- Preserve already configured Pro rates while moving to the group-scoped model.
INSERT INTO omnio_pro_group_settings (group_id, rate_multiplier, pro_only)
SELECT DISTINCT ON (group_id)
       group_id,
       rate_multiplier,
       pro_only
FROM membership_level_group_benefits
WHERE rate_multiplier IS NOT NULL OR pro_only = TRUE
ORDER BY group_id, pro_only DESC, updated_at DESC, id DESC
ON CONFLICT (group_id) DO NOTHING;

COMMENT ON TABLE omnio_pro_group_settings IS
    'Group-scoped Omnio Pro multiplier and Pro-only visibility (FORK-002)';
COMMENT ON COLUMN omnio_pro_group_settings.rate_multiplier IS
    'Final multiplier used by active Omnio Pro users; NULL falls back to the public group multiplier';
COMMENT ON COLUMN omnio_pro_group_settings.pro_only IS
    'When true, only users with an active Omnio Pro grant can see and bind the group';
