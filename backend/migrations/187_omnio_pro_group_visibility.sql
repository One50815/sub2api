-- FORK-002 extension: separate Omnio Pro-only group visibility from access.
-- Existing membership benefits remain valid; this flag only controls whether
-- a group is hidden from users without the corresponding active Pro benefit.
ALTER TABLE membership_level_group_benefits
    ADD COLUMN IF NOT EXISTS pro_only BOOLEAN NOT NULL DEFAULT FALSE;

COMMENT ON COLUMN membership_level_group_benefits.pro_only IS
    'When true, the group is only listed/bindable for users with this active Omnio Pro benefit';
