-- FORK-003: Omnio Pro is the only public membership product.
-- Keep subscription tables and paid-order fulfillment for historical data,
-- but migrate an empty Pro catalog from existing saleable plans on upgrade.

DO $$
DECLARE
    pro_level_id BIGINT;
BEGIN
    IF NOT EXISTS (SELECT 1 FROM membership_offers) THEN
        INSERT INTO membership_levels
            (name, slug, description, rank, badge_color, active, sort_order)
        VALUES
            ('Omnio Pro', 'omnio-pro', 'Omnio Pro 统一权益方案', 100, '#2563eb', TRUE, 0)
        ON CONFLICT (slug) DO UPDATE SET active = TRUE
        RETURNING id INTO pro_level_id;

        IF pro_level_id IS NULL THEN
            SELECT id INTO pro_level_id FROM membership_levels WHERE slug = 'omnio-pro';
        END IF;

        INSERT INTO membership_offers
            (level_id, name, description, price, original_price, currency,
             duration_days, for_sale, sort_order)
        SELECT pro_level_id,
               COALESCE(NULLIF(TRIM(sp.name), ''), 'Omnio Pro'),
               COALESCE(NULLIF(TRIM(sp.description), ''), '迁移自历史订阅方案'),
               sp.price,
               sp.original_price,
               COALESCE(NULLIF(UPPER(TRIM(sp.currency)), ''), 'USD'),
               CASE LOWER(COALESCE(sp.validity_unit, 'day'))
                   WHEN 'year' THEN GREATEST(sp.validity_days * 365, 1)
                   WHEN 'month' THEN GREATEST(sp.validity_days * 30, 1)
                   ELSE GREATEST(sp.validity_days, 1)
               END,
               sp.for_sale,
               sp.sort_order
        FROM subscription_plans sp
        WHERE sp.for_sale = TRUE;

        INSERT INTO membership_level_group_benefits
            (level_id, group_id, allow_access, rate_multiplier, pro_only)
SELECT DISTINCT pro_level_id, sp.group_id, TRUE, NULL::DECIMAL(12, 6), FALSE
        FROM subscription_plans sp
        WHERE sp.for_sale = TRUE
        ON CONFLICT (level_id, group_id) DO NOTHING;

        INSERT INTO membership_audit_logs
            (level_id, action, source_type, source_id, operator, detail)
        VALUES
            (pro_level_id, 'SUBSCRIPTIONS_RETIRED', 'migration', '188', 'system',
             jsonb_build_object('source', 'subscription_plans'));
    END IF;
END $$;

COMMENT ON TABLE membership_offers IS
    'Omnio Pro offers; public subscription purchases are retired (FORK-003)';
