package service

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestListOmnioProQuotaProgressReturnsOneSharedLevelPool(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	mock.ExpectQuery(`(?s)FROM effective_level e.*LEFT JOIN omnio_pro_level_shared_quota_usage u.*u\.user_id = \$1 AND u\.level_id = e\.id`).
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{
			"level_id",
			"level_name",
			"daily_limit_usd",
			"daily_used_usd",
			"monthly_limit_usd",
			"monthly_used_usd",
		}).AddRow(int64(4), "Omnio Pro Max", 10.0, 5.0, 100.0, 25.0))

	items, err := NewMembershipService(db).ListOmnioProQuotaProgress(context.Background(), 7)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, int64(4), items[0].LevelID)
	require.InDelta(t, 5.0, items[0].DailyRemainingUSD, 0.000001)
	require.InDelta(t, 75.0, items[0].MonthlyRemainingUSD, 0.000001)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestHasAvailableOmnioProQuotaRequiresEligibleGroupAndSharedPool(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer func() { _ = db.Close() }()

	mock.ExpectQuery(`(?s)JOIN membership_level_group_benefits b.*LEFT JOIN omnio_pro_level_shared_quota_usage u.*b\.group_id = \$2.*b\.allow_access = TRUE`).
		WithArgs(int64(7), int64(9)).
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

	available := NewMembershipService(db).HasAvailableOmnioProQuota(context.Background(), 7, 9)
	require.True(t, available)
	require.NoError(t, mock.ExpectationsWereMet())
}
