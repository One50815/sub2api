package service

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestListMembershipGrantsSupportsAdminFilters(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"id",
		"user_id",
		"email",
		"level_id",
		"level_name",
		"level_rank",
		"badge_color",
		"source_type",
		"source_id",
		"starts_at",
		"expires_at",
		"status",
		"notes",
	}).AddRow(
		int64(7),
		int64(42),
		"user@example.com",
		int64(3),
		"Omnio Pro Max",
		20,
		"#2563eb",
		"manual",
		"manual:42:1",
		now,
		now.Add(30*24*time.Hour),
		MembershipGrantStatusActive,
		"活动赠送",
	)
	mock.ExpectQuery(regexp.QuoteMeta("FROM user_membership_grants g")).
		WithArgs(int64(42), MembershipGrantStatusActive, 50).
		WillReturnRows(rows)

	userID := int64(42)
	items, err := NewMembershipService(db).ListGrants(
		context.Background(),
		&userID,
		MembershipGrantStatusActive,
		50,
	)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, int64(7), items[0].ID)
	require.Equal(t, "user@example.com", items[0].UserEmail)
	require.Equal(t, "Omnio Pro Max", items[0].LevelName)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListMembershipGrantsRejectsInvalidStatus(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	_, err = NewMembershipService(db).ListGrants(context.Background(), nil, "deleted", 100)
	require.Error(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRevokeMembershipGrantWritesAudit(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	mock.ExpectQuery(regexp.QuoteMeta(
		"UPDATE user_membership_grants SET status='revoked', revoked_at=NOW(), updated_at=NOW()",
	)).
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"user_id", "level_id"}).AddRow(int64(42), int64(3)))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO membership_audit_logs")).
		WithArgs(int64(42), int64(3), "GRANT_REVOKED", "grant", "7", "admin:1", "{}").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = NewMembershipService(db).RevokeGrant(context.Background(), 7, "admin:1")
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
