package service

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestUpsertGroupBenefitOnlyWritesSelectedLevel(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO membership_level_group_benefits")).
		WithArgs(int64(3), int64(9), true, 1.25, 120, true, 5.0, 50.0).
		WillReturnResult(sqlmock.NewResult(1, 1))

	daily, monthly := 5.0, 50.0
	err = NewMembershipService(db).UpsertGroupBenefit(context.Background(), MembershipGroupBenefit{
		LevelID:        3,
		GroupID:        9,
		AllowAccess:    true,
		ProOnly:        true,
		RateMultiplier: membershipFloat64Ptr(1.25),
		RPMLimit:       membershipIntPtr(120),
		DailyFreeUSD:   &daily,
		MonthlyFreeUSD: &monthly,
	})
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteGroupBenefitDoesNotDeleteOtherLevels(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectExec(regexp.QuoteMeta(
		"DELETE FROM membership_level_group_benefits WHERE level_id=$1 AND group_id=$2",
	)).WithArgs(int64(3), int64(9)).WillReturnResult(sqlmock.NewResult(0, 1))

	err = NewMembershipService(db).DeleteGroupBenefit(context.Background(), 3, 9)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func membershipFloat64Ptr(value float64) *float64 {
	return &value
}

func membershipIntPtr(value int) *int {
	return &value
}
