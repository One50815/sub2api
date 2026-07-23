package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const (
	MembershipGrantStatusActive  = "active"
	MembershipGrantStatusRevoked = "revoked"
	SubscriptionOverageBlock     = "block"
	SubscriptionOverageWallet    = "wallet"
)

var (
	ErrMembershipLevelNotFound = infraerrors.NotFound("MEMBERSHIP_LEVEL_NOT_FOUND", "membership level not found")
	ErrMembershipOfferNotFound = infraerrors.NotFound("MEMBERSHIP_OFFER_NOT_FOUND", "membership offer not found")
	ErrMembershipGrantNotFound = infraerrors.NotFound("MEMBERSHIP_GRANT_NOT_FOUND", "membership grant not found")
)

type MembershipLevel struct {
	ID               int64                    `json:"id"`
	Name             string                   `json:"name"`
	Slug             string                   `json:"slug"`
	Description      string                   `json:"description"`
	Rank             int                      `json:"rank"`
	BadgeColor       string                   `json:"badge_color"`
	ConcurrencyBonus int                      `json:"concurrency_bonus"`
	PrioritySupport  bool                     `json:"priority_support"`
	Active           bool                     `json:"active"`
	SortOrder        int                      `json:"sort_order"`
	GroupBenefits    []MembershipGroupBenefit `json:"group_benefits,omitempty"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
}

type MembershipGroupBenefit struct {
	ID             int64    `json:"id"`
	LevelID        int64    `json:"level_id"`
	GroupID        int64    `json:"group_id"`
	GroupName      string   `json:"group_name"`
	AllowAccess    bool     `json:"allow_access"`
	ProOnly        bool     `json:"pro_only"`
	RateMultiplier *float64 `json:"rate_multiplier"`
	RPMLimit       *int     `json:"rpm_limit"`
	DailyFreeUSD   *float64 `json:"daily_free_usd"`
	MonthlyFreeUSD *float64 `json:"monthly_free_usd"`
}

// OmnioProGroupSetting is retained for the legacy group-settings admin API.
// New management surfaces use MembershipGroupBenefit as the source of truth.
type OmnioProGroupSetting struct {
	GroupID        int64    `json:"group_id"`
	RateMultiplier *float64 `json:"rate_multiplier"`
	ProOnly        bool     `json:"pro_only"`
	DailyFreeUSD   *float64 `json:"daily_free_usd"`
	MonthlyFreeUSD *float64 `json:"monthly_free_usd"`
}

type OmnioProQuotaProgress struct {
	LevelID             int64   `json:"level_id"`
	LevelName           string  `json:"level_name"`
	GroupID             int64   `json:"group_id"`
	GroupName           string  `json:"group_name"`
	DailyLimitUSD       float64 `json:"daily_limit_usd"`
	DailyUsedUSD        float64 `json:"daily_used_usd"`
	DailyRemainingUSD   float64 `json:"daily_remaining_usd"`
	MonthlyLimitUSD     float64 `json:"monthly_limit_usd"`
	MonthlyUsedUSD      float64 `json:"monthly_used_usd"`
	MonthlyRemainingUSD float64 `json:"monthly_remaining_usd"`
}

// EffectiveMembershipGroupBenefit is the highest active Omnio Pro benefit
// resolved for a user. It is deliberately separate from the public group
// multiplier and from manually configured user/group rates.
type EffectiveMembershipGroupBenefit struct {
	MembershipGroupBenefit
	LevelName string `json:"level_name"`
	LevelRank int    `json:"level_rank"`
}

type MembershipOffer struct {
	ID            int64    `json:"id"`
	LevelID       int64    `json:"level_id"`
	LevelName     string   `json:"level_name"`
	LevelRank     int      `json:"level_rank"`
	BadgeColor    string   `json:"badge_color"`
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Price         float64  `json:"price"`
	OriginalPrice *float64 `json:"original_price"`
	Currency      string   `json:"currency"`
	DurationDays  int      `json:"duration_days"`
	ForSale       bool     `json:"for_sale"`
	SortOrder     int      `json:"sort_order"`
}

type MembershipGrant struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	LevelID    int64     `json:"level_id"`
	LevelName  string    `json:"level_name"`
	LevelRank  int       `json:"level_rank"`
	BadgeColor string    `json:"badge_color"`
	SourceType string    `json:"source_type"`
	SourceID   string    `json:"source_id"`
	StartsAt   time.Time `json:"starts_at"`
	ExpiresAt  time.Time `json:"expires_at"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes"`
}

type SubscriptionEntitlement struct {
	ID              int64      `json:"id"`
	OrderID         *int64     `json:"order_id"`
	SubscriptionID  *int64     `json:"subscription_id"`
	PlanID          *int64     `json:"plan_id"`
	GroupID         int64      `json:"group_id"`
	PlanName        string     `json:"plan_name"`
	DailyLimitUSD   *float64   `json:"daily_limit_usd"`
	WeeklyLimitUSD  *float64   `json:"weekly_limit_usd"`
	MonthlyLimitUSD *float64   `json:"monthly_limit_usd"`
	RateMultiplier  float64    `json:"rate_multiplier"`
	OveragePolicy   string     `json:"overage_policy"`
	StartsAt        *time.Time `json:"starts_at"`
	ExpiresAt       *time.Time `json:"expires_at"`
	Status          string     `json:"status"`
}

type MembershipSummary struct {
	EffectiveLevel *MembershipLevel         `json:"effective_level"`
	Grants         []MembershipGrant        `json:"grants"`
	Benefits       []MembershipGroupBenefit `json:"benefits"`
	Offers         []MembershipOffer        `json:"offers"`
	QuotaProgress  []OmnioProQuotaProgress  `json:"quota_progress"`
}

type MembershipCatalog struct {
	Levels       []MembershipLevel       `json:"levels"`
	Offers       []MembershipOffer       `json:"offers"`
	PlanBenefits []PlanMembershipBenefit `json:"plan_benefits"`
}

type PlanMembershipBenefit struct {
	PlanID       int64  `json:"plan_id"`
	PlanName     string `json:"plan_name"`
	LevelID      int64  `json:"level_id"`
	LevelName    string `json:"level_name"`
	DurationDays *int   `json:"duration_days"`
}

type MembershipAuditLog struct {
	ID         int64          `json:"id"`
	UserID     *int64         `json:"user_id"`
	LevelID    *int64         `json:"level_id"`
	Action     string         `json:"action"`
	SourceType string         `json:"source_type"`
	SourceID   string         `json:"source_id"`
	Operator   string         `json:"operator"`
	Detail     map[string]any `json:"detail"`
	CreatedAt  time.Time      `json:"created_at"`
}

type membershipSQLExecutor interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}

type MembershipService struct {
	db *sql.DB
}

func NewMembershipService(db *sql.DB) *MembershipService {
	return &MembershipService{db: db}
}

func (s *MembershipService) ListOffers(ctx context.Context, forSaleOnly bool) ([]MembershipOffer, error) {
	query := `
		SELECT o.id, o.level_id, l.name, l.rank, l.badge_color, o.name, o.description,
		       o.price, o.original_price, o.currency, o.duration_days, o.for_sale, o.sort_order
		FROM membership_offers o
		JOIN membership_levels l ON l.id = o.level_id
		WHERE ($1 = FALSE OR (o.for_sale = TRUE AND l.active = TRUE))
		ORDER BY l.rank, o.sort_order, o.id`
	rows, err := s.db.QueryContext(ctx, query, forSaleOnly)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	offers := make([]MembershipOffer, 0)
	for rows.Next() {
		var item MembershipOffer
		var original sql.NullFloat64
		if err := rows.Scan(&item.ID, &item.LevelID, &item.LevelName, &item.LevelRank, &item.BadgeColor,
			&item.Name, &item.Description, &item.Price, &original, &item.Currency, &item.DurationDays,
			&item.ForSale, &item.SortOrder); err != nil {
			return nil, err
		}
		if original.Valid {
			v := original.Float64
			item.OriginalPrice = &v
		}
		offers = append(offers, item)
	}
	return offers, rows.Err()
}

func (s *MembershipService) GetOffer(ctx context.Context, id int64, requireForSale bool) (*MembershipOffer, error) {
	query := `
		SELECT o.id, o.level_id, l.name, l.rank, l.badge_color, o.name, o.description,
		       o.price, o.original_price, o.currency, o.duration_days, o.for_sale, o.sort_order
		FROM membership_offers o JOIN membership_levels l ON l.id = o.level_id
		WHERE o.id = $1 AND ($2 = FALSE OR (o.for_sale = TRUE AND l.active = TRUE))`
	var item MembershipOffer
	var original sql.NullFloat64
	err := s.db.QueryRowContext(ctx, query, id, requireForSale).Scan(
		&item.ID, &item.LevelID, &item.LevelName, &item.LevelRank, &item.BadgeColor,
		&item.Name, &item.Description, &item.Price, &original, &item.Currency,
		&item.DurationDays, &item.ForSale, &item.SortOrder,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrMembershipOfferNotFound
	}
	if err != nil {
		return nil, err
	}
	if original.Valid {
		v := original.Float64
		item.OriginalPrice = &v
	}
	return &item, nil
}

func (s *MembershipService) GetSummary(ctx context.Context, userID int64) (*MembershipSummary, error) {
	grants, err := s.listUserGrants(ctx, userID)
	if err != nil {
		return nil, err
	}
	effective, err := s.getEffectiveLevel(ctx, userID)
	if err != nil {
		return nil, err
	}
	benefits := make([]MembershipGroupBenefit, 0)
	if effective != nil {
		benefits, err = s.listLevelBenefits(ctx, effective.ID)
		if err != nil {
			return nil, err
		}
		effective.GroupBenefits = benefits
	}
	offers, err := s.ListOffers(ctx, true)
	if err != nil {
		return nil, err
	}
	quotaProgress := make([]OmnioProQuotaProgress, 0)
	if effective != nil {
		quotaProgress, err = s.ListOmnioProQuotaProgress(ctx, userID)
		if err != nil {
			return nil, err
		}
	}
	return &MembershipSummary{
		EffectiveLevel: effective,
		Grants:         grants,
		Benefits:       benefits,
		Offers:         offers,
		QuotaProgress:  quotaProgress,
	}, nil
}

func (s *MembershipService) ListOmnioProQuotaProgress(ctx context.Context, userID int64) ([]OmnioProQuotaProgress, error) {
	rows, err := s.db.QueryContext(ctx, `
		WITH effective_level AS (
			SELECT ml.id, ml.name
			FROM user_membership_grants mg
			JOIN membership_levels ml ON ml.id = mg.level_id
			WHERE mg.user_id = $1 AND mg.status = 'active'
			  AND mg.starts_at <= NOW() AND mg.expires_at > NOW() AND ml.active = TRUE
			ORDER BY ml.rank DESC, mg.expires_at DESC, ml.id DESC
			LIMIT 1
		), period AS (
			SELECT
				(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Shanghai')::DATE AS day_start,
				date_trunc('month', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Shanghai')::DATE AS month_start
		)
		SELECT e.id, e.name, b.group_id, g.name,
		       COALESCE(b.daily_free_usd, 0),
		       CASE WHEN u.daily_window_start = p.day_start THEN COALESCE(u.daily_used_usd, 0) ELSE 0 END,
		       COALESCE(b.monthly_free_usd, 0),
		       CASE WHEN u.monthly_window_start = p.month_start THEN COALESCE(u.monthly_used_usd, 0) ELSE 0 END
		FROM effective_level e
		JOIN membership_level_group_benefits b ON b.level_id = e.id
		JOIN groups g ON g.id = b.group_id AND g.deleted_at IS NULL
		CROSS JOIN period p
		LEFT JOIN omnio_pro_level_quota_usage u
		  ON u.user_id = $1 AND u.level_id = e.id AND u.group_id = b.group_id
		WHERE (COALESCE(b.daily_free_usd, 0) > 0 OR COALESCE(b.monthly_free_usd, 0) > 0)
		ORDER BY g.sort_order, g.id`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]OmnioProQuotaProgress, 0)
	for rows.Next() {
		var item OmnioProQuotaProgress
		if err := rows.Scan(
			&item.LevelID,
			&item.LevelName,
			&item.GroupID,
			&item.GroupName,
			&item.DailyLimitUSD,
			&item.DailyUsedUSD,
			&item.MonthlyLimitUSD,
			&item.MonthlyUsedUSD,
		); err != nil {
			return nil, err
		}
		item.DailyRemainingUSD = remainingQuota(item.DailyLimitUSD, item.DailyUsedUSD)
		item.MonthlyRemainingUSD = remainingQuota(item.MonthlyLimitUSD, item.MonthlyUsedUSD)
		items = append(items, item)
	}
	return items, rows.Err()
}

func remainingQuota(limit, used float64) float64 {
	if limit <= 0 {
		return 0
	}
	if used >= limit {
		return 0
	}
	return limit - used
}

func (s *MembershipService) HasAvailableOmnioProQuota(ctx context.Context, userID, groupID int64) bool {
	if s == nil || s.db == nil || userID <= 0 || groupID <= 0 {
		return false
	}
	var available bool
	err := s.db.QueryRowContext(ctx, `
		WITH effective_level AS (
			SELECT ml.id
			FROM user_membership_grants mg
			JOIN membership_levels ml ON ml.id = mg.level_id
			WHERE mg.user_id = $1 AND mg.status = 'active'
			  AND mg.starts_at <= NOW() AND mg.expires_at > NOW() AND ml.active = TRUE
			ORDER BY ml.rank DESC, mg.expires_at DESC, ml.id DESC
			LIMIT 1
		), period AS (
			SELECT
				(CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Shanghai')::DATE AS day_start,
				date_trunc('month', CURRENT_TIMESTAMP AT TIME ZONE 'Asia/Shanghai')::DATE AS month_start
		)
		SELECT EXISTS (
			SELECT 1
			FROM effective_level e
			JOIN membership_level_group_benefits b ON b.level_id = e.id
			CROSS JOIN period p
			LEFT JOIN omnio_pro_level_quota_usage u
			  ON u.user_id = $1 AND u.level_id = e.id AND u.group_id = b.group_id
			WHERE b.group_id = $2
			  AND (
				(COALESCE(b.daily_free_usd, 0) > 0 AND
				 CASE WHEN u.daily_window_start = p.day_start THEN COALESCE(u.daily_used_usd, 0) ELSE 0 END < b.daily_free_usd)
				OR
				(COALESCE(b.monthly_free_usd, 0) > 0 AND
				 CASE WHEN u.monthly_window_start = p.month_start THEN COALESCE(u.monthly_used_usd, 0) ELSE 0 END < b.monthly_free_usd)
			  )
			  AND (COALESCE(b.daily_free_usd, 0) <= 0 OR
				   CASE WHEN u.daily_window_start = p.day_start THEN COALESCE(u.daily_used_usd, 0) ELSE 0 END < b.daily_free_usd)
			  AND (COALESCE(b.monthly_free_usd, 0) <= 0 OR
				   CASE WHEN u.monthly_window_start = p.month_start THEN COALESCE(u.monthly_used_usd, 0) ELSE 0 END < b.monthly_free_usd)
		)`, userID, groupID).Scan(&available)
	return err == nil && available
}

func (s *MembershipService) getEffectiveLevel(ctx context.Context, userID int64) (*MembershipLevel, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT l.id, l.name, l.slug, l.description, l.rank, l.badge_color,
		       l.concurrency_bonus, l.priority_support, l.active, l.sort_order, l.created_at, l.updated_at
		FROM user_membership_grants g
		JOIN membership_levels l ON l.id = g.level_id
		WHERE g.user_id = $1 AND g.status = 'active' AND g.starts_at <= NOW()
		  AND g.expires_at > NOW() AND l.active = TRUE
		ORDER BY l.rank DESC, g.expires_at DESC, l.id DESC LIMIT 1`, userID)
	var level MembershipLevel
	err := row.Scan(&level.ID, &level.Name, &level.Slug, &level.Description, &level.Rank,
		&level.BadgeColor, &level.ConcurrencyBonus, &level.PrioritySupport, &level.Active,
		&level.SortOrder, &level.CreatedAt, &level.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return &level, err
}

func (s *MembershipService) listUserGrants(ctx context.Context, userID int64) ([]MembershipGrant, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT g.id, g.user_id, g.level_id, l.name, l.rank, l.badge_color, g.source_type,
		       g.source_id, g.starts_at, g.expires_at, g.status, g.notes
		FROM user_membership_grants g JOIN membership_levels l ON l.id = g.level_id
		WHERE g.user_id = $1 ORDER BY g.expires_at DESC, g.id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	items := make([]MembershipGrant, 0)
	for rows.Next() {
		var item MembershipGrant
		if err := rows.Scan(&item.ID, &item.UserID, &item.LevelID, &item.LevelName, &item.LevelRank,
			&item.BadgeColor, &item.SourceType, &item.SourceID, &item.StartsAt, &item.ExpiresAt,
			&item.Status, &item.Notes); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MembershipService) listLevelBenefits(ctx context.Context, levelID int64) ([]MembershipGroupBenefit, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT b.id, b.level_id, b.group_id, g.name, b.allow_access, b.pro_only,
		       b.rate_multiplier, b.rpm_limit, b.daily_free_usd, b.monthly_free_usd
		FROM membership_level_group_benefits b
		JOIN groups g ON g.id = b.group_id
		WHERE b.level_id = $1
		ORDER BY g.sort_order, g.id`, levelID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	items := make([]MembershipGroupBenefit, 0)
	for rows.Next() {
		var item MembershipGroupBenefit
		var rate, dailyFree, monthlyFree sql.NullFloat64
		var rpm sql.NullInt32
		if err := rows.Scan(
			&item.ID,
			&item.LevelID,
			&item.GroupID,
			&item.GroupName,
			&item.AllowAccess,
			&item.ProOnly,
			&rate,
			&rpm,
			&dailyFree,
			&monthlyFree,
		); err != nil {
			return nil, err
		}
		if rate.Valid {
			v := rate.Float64
			item.RateMultiplier = &v
		}
		if rpm.Valid {
			v := int(rpm.Int32)
			item.RPMLimit = &v
		}
		if dailyFree.Valid {
			v := dailyFree.Float64
			item.DailyFreeUSD = &v
		}
		if monthlyFree.Valid {
			v := monthlyFree.Float64
			item.MonthlyFreeUSD = &v
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

// GetEffectiveGroupBenefits returns the highest active Omnio Pro level's
// group benefits for a user. A missing map entry means the user has no Pro
// benefit for that group; the normal group multiplier remains authoritative.
func (s *MembershipService) GetEffectiveGroupBenefits(ctx context.Context, userID int64) (map[int64]EffectiveMembershipGroupBenefit, error) {
	rows, err := s.db.QueryContext(ctx, `
		WITH effective_level AS (
			SELECT ml.id, ml.name, ml.rank
			FROM user_membership_grants mg
			JOIN membership_levels ml ON ml.id = mg.level_id
			WHERE mg.user_id = $1 AND mg.status = 'active'
			  AND mg.starts_at <= NOW() AND mg.expires_at > NOW() AND ml.active = TRUE
			ORDER BY ml.rank DESC, mg.expires_at DESC, ml.id DESC
			LIMIT 1
		)
		SELECT b.id, b.level_id, b.group_id, g.name, b.allow_access, b.pro_only,
		       b.rate_multiplier, b.rpm_limit, b.daily_free_usd, b.monthly_free_usd,
		       e.name, e.rank
		FROM effective_level e
		JOIN membership_level_group_benefits b ON b.level_id = e.id
		JOIN groups g ON g.id = b.group_id
		ORDER BY g.sort_order, g.id`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	result := make(map[int64]EffectiveMembershipGroupBenefit)
	for rows.Next() {
		var item EffectiveMembershipGroupBenefit
		var rate, dailyFree, monthlyFree sql.NullFloat64
		var rpm sql.NullInt32
		if err := rows.Scan(&item.ID, &item.LevelID, &item.GroupID, &item.GroupName, &item.AllowAccess,
			&item.ProOnly, &rate, &rpm, &dailyFree, &monthlyFree, &item.LevelName, &item.LevelRank); err != nil {
			return nil, err
		}
		if rate.Valid {
			v := rate.Float64
			item.RateMultiplier = &v
		}
		if rpm.Valid {
			v := int(rpm.Int32)
			item.RPMLimit = &v
		}
		if dailyFree.Valid {
			v := dailyFree.Float64
			item.DailyFreeUSD = &v
		}
		if monthlyFree.Valid {
			v := monthlyFree.Float64
			item.MonthlyFreeUSD = &v
		}
		result[item.GroupID] = item
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MembershipService) GetEffectiveGroupBenefit(ctx context.Context, userID, groupID int64) (EffectiveMembershipGroupBenefit, bool, error) {
	benefits, err := s.GetEffectiveGroupBenefits(ctx, userID)
	if err != nil {
		return EffectiveMembershipGroupBenefit{}, false, err
	}
	benefit, ok := benefits[groupID]
	return benefit, ok, nil
}

// GetProOnlyGroupIDs returns groups reserved by at least one Omnio Pro level.
// A user's effective level still needs its own allow_access benefit.
func (s *MembershipService) GetProOnlyGroupIDs(ctx context.Context) (map[int64]struct{}, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT DISTINCT group_id
		FROM membership_level_group_benefits
		WHERE pro_only = TRUE`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	result := make(map[int64]struct{})
	for rows.Next() {
		var groupID int64
		if err := rows.Scan(&groupID); err != nil {
			return nil, err
		}
		result[groupID] = struct{}{}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func (s *MembershipService) GetOmnioProGroupSetting(ctx context.Context, groupID int64) (OmnioProGroupSetting, error) {
	setting := OmnioProGroupSetting{GroupID: groupID}
	var rate, dailyFree, monthlyFree sql.NullFloat64
	err := s.db.QueryRowContext(ctx, `
		SELECT rate_multiplier, pro_only, daily_free_usd, monthly_free_usd
		FROM omnio_pro_group_settings
		WHERE group_id = $1`, groupID).Scan(&rate, &setting.ProOnly, &dailyFree, &monthlyFree)
	if errors.Is(err, sql.ErrNoRows) {
		return setting, nil
	}
	if err != nil {
		return OmnioProGroupSetting{}, err
	}
	if rate.Valid {
		value := rate.Float64
		setting.RateMultiplier = &value
	}
	if dailyFree.Valid {
		value := dailyFree.Float64
		setting.DailyFreeUSD = &value
	}
	if monthlyFree.Valid {
		value := monthlyFree.Float64
		setting.MonthlyFreeUSD = &value
	}
	return setting, nil
}

func (s *MembershipService) UpsertOmnioProGroupSetting(ctx context.Context, setting OmnioProGroupSetting) (OmnioProGroupSetting, error) {
	if setting.GroupID <= 0 {
		return OmnioProGroupSetting{}, infraerrors.BadRequest("INVALID_OMNIO_PRO_GROUP", "group_id is required")
	}
	if setting.RateMultiplier != nil && *setting.RateMultiplier < 0 {
		return OmnioProGroupSetting{}, infraerrors.BadRequest("INVALID_OMNIO_PRO_GROUP", "rate_multiplier must be non-negative")
	}
	if setting.DailyFreeUSD != nil && *setting.DailyFreeUSD < 0 {
		return OmnioProGroupSetting{}, infraerrors.BadRequest("INVALID_OMNIO_PRO_GROUP", "daily_free_usd must be non-negative")
	}
	if setting.MonthlyFreeUSD != nil && *setting.MonthlyFreeUSD < 0 {
		return OmnioProGroupSetting{}, infraerrors.BadRequest("INVALID_OMNIO_PRO_GROUP", "monthly_free_usd must be non-negative")
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return OmnioProGroupSetting{}, err
	}
	defer func() { _ = tx.Rollback() }()

	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM groups WHERE id=$1)`, setting.GroupID).Scan(&exists); err != nil {
		return OmnioProGroupSetting{}, err
	}
	if !exists {
		return OmnioProGroupSetting{}, ErrGroupNotFound
	}
	if setting.RateMultiplier == nil && !setting.ProOnly && setting.DailyFreeUSD == nil && setting.MonthlyFreeUSD == nil {
		if _, err := tx.ExecContext(ctx, `
			UPDATE membership_level_group_benefits
			SET rate_multiplier=NULL, pro_only=FALSE, daily_free_usd=NULL,
			    monthly_free_usd=NULL, updated_at=NOW()
			WHERE group_id=$1`, setting.GroupID); err != nil {
			return OmnioProGroupSetting{}, err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM omnio_pro_group_settings WHERE group_id=$1`, setting.GroupID); err != nil {
			return OmnioProGroupSetting{}, err
		}
		if err := tx.Commit(); err != nil {
			return OmnioProGroupSetting{}, err
		}
		return setting, nil
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO membership_level_group_benefits (
			level_id, group_id, allow_access, rate_multiplier, rpm_limit, pro_only,
			daily_free_usd, monthly_free_usd
		)
		SELECT id, $1, $3, $2, NULL, $3, $4, $5
		FROM membership_levels
		ON CONFLICT (level_id, group_id) DO UPDATE SET
			allow_access=membership_level_group_benefits.allow_access OR EXCLUDED.pro_only,
			rate_multiplier=EXCLUDED.rate_multiplier,
			pro_only=EXCLUDED.pro_only,
			daily_free_usd=EXCLUDED.daily_free_usd,
			monthly_free_usd=EXCLUDED.monthly_free_usd,
			updated_at=NOW()`,
		setting.GroupID,
		setting.RateMultiplier,
		setting.ProOnly,
		setting.DailyFreeUSD,
		setting.MonthlyFreeUSD,
	); err != nil {
		return OmnioProGroupSetting{}, err
	}
	_, err = tx.ExecContext(ctx, `
		INSERT INTO omnio_pro_group_settings (
			group_id, rate_multiplier, pro_only, daily_free_usd, monthly_free_usd
		)
		VALUES ($1,$2,$3,$4,$5)
		ON CONFLICT (group_id) DO UPDATE SET
			rate_multiplier=EXCLUDED.rate_multiplier,
			pro_only=EXCLUDED.pro_only,
			daily_free_usd=EXCLUDED.daily_free_usd,
			monthly_free_usd=EXCLUDED.monthly_free_usd,
			updated_at=NOW()`,
		setting.GroupID,
		setting.RateMultiplier,
		setting.ProOnly,
		setting.DailyFreeUSD,
		setting.MonthlyFreeUSD,
	)
	if err != nil {
		return OmnioProGroupSetting{}, err
	}
	if err := tx.Commit(); err != nil {
		return OmnioProGroupSetting{}, err
	}
	return s.GetOmnioProGroupSetting(ctx, setting.GroupID)
}

func (s *MembershipService) GetCatalog(ctx context.Context) (*MembershipCatalog, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, slug, description, rank, badge_color, concurrency_bonus,
		       priority_support, active, sort_order, created_at, updated_at
		FROM membership_levels ORDER BY rank, sort_order, id`)
	if err != nil {
		return nil, err
	}
	levels := make([]MembershipLevel, 0)
	for rows.Next() {
		var level MembershipLevel
		if err := rows.Scan(&level.ID, &level.Name, &level.Slug, &level.Description, &level.Rank,
			&level.BadgeColor, &level.ConcurrencyBonus, &level.PrioritySupport, &level.Active,
			&level.SortOrder, &level.CreatedAt, &level.UpdatedAt); err != nil {
			_ = rows.Close()
			return nil, err
		}
		level.GroupBenefits, err = s.listLevelBenefits(ctx, level.ID)
		if err != nil {
			_ = rows.Close()
			return nil, err
		}
		levels = append(levels, level)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	offers, err := s.ListOffers(ctx, false)
	if err != nil {
		return nil, err
	}
	planBenefits, err := s.listPlanBenefits(ctx)
	if err != nil {
		return nil, err
	}
	return &MembershipCatalog{Levels: levels, Offers: offers, PlanBenefits: planBenefits}, nil
}

func (s *MembershipService) UpsertLevel(ctx context.Context, level MembershipLevel) (*MembershipLevel, error) {
	level.Name = strings.TrimSpace(level.Name)
	level.Slug = strings.ToLower(strings.TrimSpace(level.Slug))
	if level.Name == "" || level.Slug == "" || level.Rank < 0 || level.ConcurrencyBonus < 0 {
		return nil, infraerrors.BadRequest("INVALID_MEMBERSHIP_LEVEL", "name, slug and non-negative limits are required")
	}
	if strings.TrimSpace(level.BadgeColor) == "" {
		level.BadgeColor = "#111827"
	}
	if level.ID == 0 {
		err := s.db.QueryRowContext(ctx, `
			INSERT INTO membership_levels (name, slug, description, rank, badge_color, concurrency_bonus, priority_support, active, sort_order)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
			RETURNING id, created_at, updated_at`, level.Name, level.Slug, level.Description, level.Rank,
			level.BadgeColor, level.ConcurrencyBonus, level.PrioritySupport, level.Active, level.SortOrder).
			Scan(&level.ID, &level.CreatedAt, &level.UpdatedAt)
		if err != nil {
			return nil, err
		}
	} else {
		result, err := s.db.ExecContext(ctx, `
			UPDATE membership_levels SET name=$2, slug=$3, description=$4, rank=$5, badge_color=$6,
			concurrency_bonus=$7, priority_support=$8, active=$9, sort_order=$10, updated_at=NOW()
			WHERE id=$1`, level.ID, level.Name, level.Slug, level.Description, level.Rank,
			level.BadgeColor, level.ConcurrencyBonus, level.PrioritySupport, level.Active, level.SortOrder)
		if err != nil {
			return nil, err
		}
		if count, _ := result.RowsAffected(); count == 0 {
			return nil, ErrMembershipLevelNotFound
		}
	}
	s.writeAudit(ctx, nil, &level.ID, "LEVEL_UPSERTED", "admin", fmt.Sprintf("%d", level.ID), "admin", map[string]any{"name": level.Name})
	return &level, nil
}

func (s *MembershipService) DeleteLevel(ctx context.Context, id int64) error {
	if id <= 0 {
		return infraerrors.BadRequest("INVALID_MEMBERSHIP_LEVEL", "level id is required")
	}
	var offers, grants, planBenefits int64
	if err := s.db.QueryRowContext(ctx, `
		SELECT
			(SELECT COUNT(*) FROM membership_offers WHERE level_id=$1),
			(SELECT COUNT(*) FROM user_membership_grants WHERE level_id=$1),
			(SELECT COUNT(*) FROM subscription_plan_membership_benefits WHERE level_id=$1)`, id).
		Scan(&offers, &grants, &planBenefits); err != nil {
		return err
	}
	if offers > 0 || grants > 0 || planBenefits > 0 {
		return infraerrors.Conflict(
			"MEMBERSHIP_LEVEL_IN_USE",
			fmt.Sprintf("level is still in use (offers=%d, grants=%d, plan_benefits=%d)", offers, grants, planBenefits),
		)
	}
	result, err := s.db.ExecContext(ctx, `DELETE FROM membership_levels WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return ErrMembershipLevelNotFound
	}
	s.writeAudit(ctx, nil, nil, "LEVEL_DELETED", "admin", fmt.Sprintf("%d", id), "admin", map[string]any{"level_id": id})
	return nil
}

func (s *MembershipService) UpsertOffer(ctx context.Context, offer MembershipOffer) (*MembershipOffer, error) {
	if offer.LevelID <= 0 || strings.TrimSpace(offer.Name) == "" || offer.Price < 0 || offer.DurationDays <= 0 {
		return nil, infraerrors.BadRequest("INVALID_MEMBERSHIP_OFFER", "level, name, price and duration are required")
	}
	offer.Currency = strings.ToUpper(strings.TrimSpace(offer.Currency))
	if offer.Currency == "" {
		offer.Currency = "USD"
	}
	if offer.ID == 0 {
		err := s.db.QueryRowContext(ctx, `
			INSERT INTO membership_offers (level_id, name, description, price, original_price, currency, duration_days, for_sale, sort_order)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`, offer.LevelID, strings.TrimSpace(offer.Name),
			offer.Description, offer.Price, offer.OriginalPrice, offer.Currency, offer.DurationDays, offer.ForSale, offer.SortOrder).
			Scan(&offer.ID)
		if err != nil {
			return nil, err
		}
	} else {
		result, err := s.db.ExecContext(ctx, `
			UPDATE membership_offers SET level_id=$2, name=$3, description=$4, price=$5,
			original_price=$6, currency=$7, duration_days=$8, for_sale=$9, sort_order=$10, updated_at=NOW()
			WHERE id=$1`, offer.ID, offer.LevelID, strings.TrimSpace(offer.Name), offer.Description,
			offer.Price, offer.OriginalPrice, offer.Currency, offer.DurationDays, offer.ForSale, offer.SortOrder)
		if err != nil {
			return nil, err
		}
		if count, _ := result.RowsAffected(); count == 0 {
			return nil, ErrMembershipOfferNotFound
		}
	}
	s.writeAudit(ctx, nil, &offer.LevelID, "OFFER_UPSERTED", "admin", fmt.Sprintf("%d", offer.ID), "admin", map[string]any{"name": offer.Name})
	return s.GetOffer(ctx, offer.ID, false)
}

func (s *MembershipService) DeleteOffer(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM membership_offers WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return ErrMembershipOfferNotFound
	}
	return nil
}

func (s *MembershipService) UpsertGroupBenefit(ctx context.Context, benefit MembershipGroupBenefit) error {
	if benefit.LevelID <= 0 || benefit.GroupID <= 0 {
		return infraerrors.BadRequest("INVALID_MEMBERSHIP_BENEFIT", "level_id and group_id are required")
	}
	if benefit.RateMultiplier != nil && *benefit.RateMultiplier < 0 {
		return infraerrors.BadRequest("INVALID_MEMBERSHIP_BENEFIT", "rate_multiplier must be non-negative")
	}
	if benefit.RPMLimit != nil && *benefit.RPMLimit < 0 {
		return infraerrors.BadRequest("INVALID_MEMBERSHIP_BENEFIT", "rpm_limit must be non-negative")
	}
	if benefit.DailyFreeUSD != nil && *benefit.DailyFreeUSD < 0 {
		return infraerrors.BadRequest("INVALID_MEMBERSHIP_BENEFIT", "daily_free_usd must be non-negative")
	}
	if benefit.MonthlyFreeUSD != nil && *benefit.MonthlyFreeUSD < 0 {
		return infraerrors.BadRequest("INVALID_MEMBERSHIP_BENEFIT", "monthly_free_usd must be non-negative")
	}
	if benefit.ProOnly && !benefit.AllowAccess {
		return infraerrors.BadRequest("INVALID_MEMBERSHIP_BENEFIT", "pro_only requires allow_access")
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO membership_level_group_benefits (
			level_id, group_id, allow_access, rate_multiplier, rpm_limit, pro_only,
			daily_free_usd, monthly_free_usd
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		ON CONFLICT (level_id, group_id) DO UPDATE SET allow_access=EXCLUDED.allow_access,
			rate_multiplier=EXCLUDED.rate_multiplier, rpm_limit=EXCLUDED.rpm_limit,
			pro_only=EXCLUDED.pro_only, daily_free_usd=EXCLUDED.daily_free_usd,
			monthly_free_usd=EXCLUDED.monthly_free_usd, updated_at=NOW()`,
		benefit.LevelID,
		benefit.GroupID,
		benefit.AllowAccess,
		benefit.RateMultiplier,
		benefit.RPMLimit,
		benefit.ProOnly,
		benefit.DailyFreeUSD,
		benefit.MonthlyFreeUSD,
	)
	return err
}

func (s *MembershipService) DeleteGroupBenefit(ctx context.Context, levelID, groupID int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM membership_level_group_benefits WHERE level_id=$1 AND group_id=$2`, levelID, groupID)
	return err
}

func (s *MembershipService) SetPlanBenefit(ctx context.Context, benefit PlanMembershipBenefit) error {
	if benefit.PlanID <= 0 || benefit.LevelID <= 0 {
		return infraerrors.BadRequest("INVALID_PLAN_MEMBERSHIP_BENEFIT", "plan_id and level_id are required")
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO subscription_plan_membership_benefits (plan_id, level_id, duration_days)
		VALUES ($1,$2,$3)
		ON CONFLICT (plan_id) DO UPDATE SET level_id=EXCLUDED.level_id,
			duration_days=EXCLUDED.duration_days, updated_at=NOW()`, benefit.PlanID, benefit.LevelID, benefit.DurationDays)
	return err
}

func (s *MembershipService) DeletePlanBenefit(ctx context.Context, planID int64) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM subscription_plan_membership_benefits WHERE plan_id=$1`, planID)
	return err
}

func (s *MembershipService) listPlanBenefits(ctx context.Context) ([]PlanMembershipBenefit, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT b.plan_id, p.name, b.level_id, l.name, b.duration_days
		FROM subscription_plan_membership_benefits b
		JOIN subscription_plans p ON p.id=b.plan_id JOIN membership_levels l ON l.id=b.level_id
		ORDER BY p.sort_order, p.id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	items := make([]PlanMembershipBenefit, 0)
	for rows.Next() {
		var item PlanMembershipBenefit
		var duration sql.NullInt32
		if err := rows.Scan(&item.PlanID, &item.PlanName, &item.LevelID, &item.LevelName, &duration); err != nil {
			return nil, err
		}
		if duration.Valid {
			v := int(duration.Int32)
			item.DurationDays = &v
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MembershipService) ListPlanBenefits(ctx context.Context) ([]PlanMembershipBenefit, error) {
	return s.listPlanBenefits(ctx)
}

func (s *MembershipService) GrantManual(ctx context.Context, userID, levelID int64, days int, notes, operator string) (*MembershipGrant, error) {
	if userID <= 0 || levelID <= 0 || days <= 0 {
		return nil, infraerrors.BadRequest("INVALID_MEMBERSHIP_GRANT", "user_id, level_id and positive days are required")
	}
	sourceID := fmt.Sprintf("manual:%d:%d", userID, time.Now().UnixNano())
	return s.createGrant(ctx, userID, levelID, "manual", sourceID, days, notes, operator)
}

func (s *MembershipService) GrantOfferFromOrder(ctx context.Context, orderID, userID, offerID int64) error {
	offer, err := s.GetOffer(ctx, offerID, false)
	if err != nil {
		return err
	}
	_, err = s.createGrant(ctx, userID, offer.LevelID, "payment_order", fmt.Sprintf("%d", orderID), offer.DurationDays, offer.Name, "system")
	return err
}

func (s *MembershipService) createGrant(ctx context.Context, userID, levelID int64, sourceType, sourceID string, days int, notes, operator string) (*MembershipGrant, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	var starts time.Time
	err = tx.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(expires_at), NOW()) FROM user_membership_grants
		WHERE user_id=$1 AND level_id=$2 AND status='active' AND expires_at > NOW()`, userID, levelID).Scan(&starts)
	if err != nil {
		return nil, err
	}
	expires := starts.AddDate(0, 0, days)
	var grantID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO user_membership_grants (user_id, level_id, source_type, source_id, starts_at, expires_at, status, notes)
		VALUES ($1,$2,$3,$4,$5,$6,'active',$7)
		ON CONFLICT (source_type, source_id, level_id) DO UPDATE SET updated_at=NOW()
		RETURNING id`, userID, levelID, sourceType, sourceID, starts, expires, notes).Scan(&grantID)
	if err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, &userID, &levelID, "GRANT_CREATED", sourceType, sourceID, operator, map[string]any{"days": days})
	grants, err := s.listUserGrants(ctx, userID)
	if err != nil {
		return nil, err
	}
	for i := range grants {
		if grants[i].ID == grantID {
			return &grants[i], nil
		}
	}
	return nil, ErrMembershipGrantNotFound
}

func (s *MembershipService) RevokeGrant(ctx context.Context, id int64, operator string) error {
	var userID, levelID int64
	err := s.db.QueryRowContext(ctx, `
		UPDATE user_membership_grants SET status='revoked', revoked_at=NOW(), updated_at=NOW()
		WHERE id=$1 AND status='active' RETURNING user_id, level_id`, id).Scan(&userID, &levelID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrMembershipGrantNotFound
	}
	if err != nil {
		return err
	}
	s.writeAudit(ctx, &userID, &levelID, "GRANT_REVOKED", "grant", fmt.Sprintf("%d", id), operator, nil)
	return nil
}

func (s *MembershipService) RevokeOrderGrants(ctx context.Context, orderID int64) error {
	sourceID := fmt.Sprintf("%d", orderID)
	_, err := s.db.ExecContext(ctx, `
		UPDATE user_membership_grants SET status='revoked', revoked_at=NOW(), updated_at=NOW()
		WHERE source_id=$1 AND source_type IN ('payment_order','subscription_order') AND status='active'`, sourceID)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `UPDATE subscription_grants SET status='revoked', updated_at=NOW() WHERE order_id=$1`, orderID)
	return err
}

func (s *MembershipService) RestoreOrderGrants(ctx context.Context, orderID int64) error {
	sourceID := fmt.Sprintf("%d", orderID)
	_, err := s.db.ExecContext(ctx, `
		UPDATE user_membership_grants SET status='active', revoked_at=NULL, updated_at=NOW()
		WHERE source_id=$1 AND source_type IN ('payment_order','subscription_order') AND status='revoked'`, sourceID)
	if err != nil {
		return err
	}
	_, err = s.db.ExecContext(ctx, `
		UPDATE subscription_grants SET status=CASE WHEN expires_at > NOW() THEN 'active' ELSE 'expired' END, updated_at=NOW()
		WHERE order_id=$1 AND status='revoked'`, orderID)
	return err
}

func (s *MembershipService) PrepareSubscriptionGrant(ctx context.Context, exec membershipSQLExecutor, orderID, userID, planID, groupID int64) error {
	if exec == nil {
		exec = s.db
	}
	_, err := exec.ExecContext(ctx, `
		INSERT INTO subscription_grants (
			order_id, user_id, plan_id, group_id, plan_name, daily_limit_usd, weekly_limit_usd,
			monthly_limit_usd, rate_multiplier, overage_policy, status
		)
		SELECT $1, $2, p.id, g.id, p.name, g.daily_limit_usd, g.weekly_limit_usd,
		       g.monthly_limit_usd, g.rate_multiplier, 'block', 'pending'
		FROM subscription_plans p JOIN groups g ON g.id=p.group_id
		WHERE p.id=$3 AND g.id=$4
		ON CONFLICT (order_id) DO NOTHING`, orderID, userID, planID, groupID)
	return err
}

func (s *MembershipService) ActivateSubscriptionGrant(ctx context.Context, orderID, userID, subscriptionID, planID int64, startsAt, expiresAt time.Time, fallbackDays int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()
	result, err := tx.ExecContext(ctx, `
		UPDATE subscription_grants SET subscription_id=$2, starts_at=$3, expires_at=$4,
			status='active', updated_at=NOW() WHERE order_id=$1`, orderID, subscriptionID, startsAt, expiresAt)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("subscription entitlement snapshot is missing")
	}
	var levelID int64
	var duration sql.NullInt32
	err = tx.QueryRowContext(ctx, `
		SELECT level_id, duration_days FROM subscription_plan_membership_benefits WHERE plan_id=$1`, planID).Scan(&levelID, &duration)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if err == nil {
		days := fallbackDays
		if duration.Valid && duration.Int32 > 0 {
			days = int(duration.Int32)
		}
		if days <= 0 {
			days = 30
		}
		var grantStart time.Time
		if err := tx.QueryRowContext(ctx, `
			SELECT COALESCE(MAX(expires_at), NOW()) FROM user_membership_grants
			WHERE user_id=$1 AND level_id=$2 AND status='active' AND expires_at > NOW()`, userID, levelID).Scan(&grantStart); err != nil {
			return err
		}
		grantExpiry := grantStart.AddDate(0, 0, days)
		_, err = tx.ExecContext(ctx, `
			INSERT INTO user_membership_grants (user_id, level_id, source_type, source_id, starts_at, expires_at, status, notes)
			VALUES ($1,$2,'subscription_order',$3,$4,$5,'active','Included with subscription')
			ON CONFLICT (source_type, source_id, level_id) DO NOTHING`, userID, levelID, fmt.Sprintf("%d", orderID), grantStart, grantExpiry)
		if err != nil {
			return err
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	s.writeAudit(ctx, &userID, nil, "SUBSCRIPTION_ENTITLEMENT_ACTIVATED", "payment_order", fmt.Sprintf("%d", orderID), "system", map[string]any{"subscription_id": subscriptionID})
	return nil
}

func (s *MembershipService) GetSubscriptionEntitlement(ctx context.Context, subscriptionID int64) (*SubscriptionEntitlement, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, order_id, subscription_id, plan_id, group_id, plan_name, daily_limit_usd,
		       weekly_limit_usd, monthly_limit_usd, rate_multiplier, overage_policy,
		       starts_at, expires_at, status
		FROM subscription_grants
		WHERE subscription_id=$1 AND status='active' AND (expires_at IS NULL OR expires_at > NOW())
		ORDER BY expires_at DESC NULLS LAST, id DESC LIMIT 1`, subscriptionID)
	var item SubscriptionEntitlement
	var orderID, subID, planID sql.NullInt64
	var daily, weekly, monthly sql.NullFloat64
	var starts, expires sql.NullTime
	err := row.Scan(&item.ID, &orderID, &subID, &planID, &item.GroupID, &item.PlanName,
		&daily, &weekly, &monthly, &item.RateMultiplier, &item.OveragePolicy, &starts, &expires, &item.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	item.OrderID = nullInt64Ptr(orderID)
	item.SubscriptionID = nullInt64Ptr(subID)
	item.PlanID = nullInt64Ptr(planID)
	item.DailyLimitUSD = nullFloat64Ptr(daily)
	item.WeeklyLimitUSD = nullFloat64Ptr(weekly)
	item.MonthlyLimitUSD = nullFloat64Ptr(monthly)
	item.StartsAt = nullTimePtr(starts)
	item.ExpiresAt = nullTimePtr(expires)
	return &item, nil
}

func (s *MembershipService) SetOveragePolicy(ctx context.Context, userID, entitlementID int64, policy string) error {
	if policy != SubscriptionOverageBlock && policy != SubscriptionOverageWallet {
		return infraerrors.BadRequest("INVALID_OVERAGE_POLICY", "overage policy must be block or wallet")
	}
	result, err := s.db.ExecContext(ctx, `
		UPDATE subscription_grants SET overage_policy=$3, updated_at=NOW()
		WHERE id=$1 AND user_id=$2 AND status='active'`, entitlementID, userID, policy)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return infraerrors.NotFound("SUBSCRIPTION_ENTITLEMENT_NOT_FOUND", "subscription entitlement not found")
	}
	return nil
}

func (s *MembershipService) IsGroupAllowed(ctx context.Context, userID, groupID int64) bool {
	var allowed bool
	err := s.db.QueryRowContext(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM user_membership_grants mg
			JOIN membership_levels ml ON ml.id=mg.level_id AND ml.active=TRUE
			JOIN membership_level_group_benefits b ON b.level_id=ml.id
			WHERE mg.user_id=$1 AND mg.status='active' AND mg.starts_at<=NOW() AND mg.expires_at>NOW()
			  AND b.group_id=$2 AND b.allow_access=TRUE
		)`, userID, groupID).Scan(&allowed)
	return err == nil && allowed
}

func (s *MembershipService) ResolveConcurrency(ctx context.Context, userID int64, base int) int {
	var bonus int
	err := s.db.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(ml.concurrency_bonus),0)
		FROM user_membership_grants mg JOIN membership_levels ml ON ml.id=mg.level_id
		WHERE mg.user_id=$1 AND mg.status='active' AND mg.starts_at<=NOW() AND mg.expires_at>NOW() AND ml.active=TRUE`, userID).Scan(&bonus)
	if err != nil || bonus <= 0 {
		return base
	}
	return base + bonus
}

func (s *MembershipService) ListAuditLogs(ctx context.Context, limit int) ([]MembershipAuditLog, error) {
	if limit <= 0 || limit > 200 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, user_id, level_id, action, source_type, source_id, operator, detail, created_at
		FROM membership_audit_logs ORDER BY created_at DESC LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	items := make([]MembershipAuditLog, 0)
	for rows.Next() {
		var item MembershipAuditLog
		var userID, levelID sql.NullInt64
		var raw []byte
		if err := rows.Scan(&item.ID, &userID, &levelID, &item.Action, &item.SourceType,
			&item.SourceID, &item.Operator, &raw, &item.CreatedAt); err != nil {
			return nil, err
		}
		item.UserID = nullInt64Ptr(userID)
		item.LevelID = nullInt64Ptr(levelID)
		item.Detail = map[string]any{}
		_ = json.Unmarshal(raw, &item.Detail)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *MembershipService) writeAudit(ctx context.Context, userID, levelID *int64, action, sourceType, sourceID, operator string, detail map[string]any) {
	if detail == nil {
		detail = map[string]any{}
	}
	raw, _ := json.Marshal(detail)
	_, _ = s.db.ExecContext(ctx, `
		INSERT INTO membership_audit_logs (user_id, level_id, action, source_type, source_id, operator, detail)
		VALUES ($1,$2,$3,$4,$5,$6,$7::jsonb)`, userID, levelID, action, sourceType, sourceID, operator, string(raw))
}

func nullInt64Ptr(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	v := value.Int64
	return &v
}

func nullFloat64Ptr(value sql.NullFloat64) *float64 {
	if !value.Valid {
		return nil
	}
	v := value.Float64
	return &v
}

func nullTimePtr(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}
	v := value.Time
	return &v
}
