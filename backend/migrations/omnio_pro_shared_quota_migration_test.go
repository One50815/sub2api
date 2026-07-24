package migrations

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOmnioProSharedQuotaMigrationPromotesLimitsAndAggregatesUsage(t *testing.T) {
	content, err := FS.ReadFile("192_omnio_pro_shared_free_quota.sql")
	require.NoError(t, err)

	sql := string(content)
	require.Contains(t, sql, "ADD COLUMN IF NOT EXISTS daily_free_usd")
	require.Contains(t, sql, "MAX(COALESCE(daily_free_usd, 0))")
	require.Contains(t, sql, "CREATE TABLE IF NOT EXISTS omnio_pro_level_shared_quota_usage")
	require.Contains(t, sql, "PRIMARY KEY (user_id, level_id)")
	require.Contains(t, sql, "SUM(")
	require.NotContains(t, sql, "PRIMARY KEY (user_id, level_id, group_id)")
}
