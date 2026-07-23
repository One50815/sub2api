package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const (
	TicketStatusPendingAdmin = "pending_admin"
	TicketStatusPendingUser  = "pending_user"
	TicketStatusResolved     = "resolved"
	TicketStatusClosed       = "closed"

	TicketCategoryAccount  = "account"
	TicketCategoryBilling  = "billing"
	TicketCategoryAPIModel = "api_model"
	TicketCategoryIncident = "incident"
	TicketCategoryFeature  = "feature"
	TicketCategoryOther    = "other"

	TicketPriorityLow    = "low"
	TicketPriorityNormal = "normal"
	TicketPriorityHigh   = "high"
	TicketPriorityUrgent = "urgent"

	TicketSenderUser  = "user"
	TicketSenderAdmin = "admin"

	TicketPermissionRead   = "read"
	TicketPermissionReply  = "reply"
	TicketPermissionManage = "manage"

	TicketSubjectMaxLength   = 200
	TicketContentMaxLength   = 10000
	TicketRequestIDMaxLength = 128
)

var (
	ErrTicketNotFound         = infraerrors.NotFound("TICKET_NOT_FOUND", "ticket not found")
	ErrTicketClosed           = infraerrors.Conflict("TICKET_CLOSED", "closed tickets cannot be changed")
	ErrTicketCenterDisabled   = infraerrors.Forbidden("TICKET_CENTER_DISABLED", "ticket center is disabled")
	ErrTicketIntakeDisabled   = infraerrors.Forbidden("TICKET_INTAKE_DISABLED", "new ticket intake is disabled")
	ErrInvalidTicketSubject   = infraerrors.BadRequest("INVALID_TICKET_SUBJECT", "ticket subject is required and must be at most 200 characters")
	ErrInvalidTicketContent   = infraerrors.BadRequest("INVALID_TICKET_CONTENT", "ticket content is required and must be at most 10000 characters")
	ErrInvalidTicketCategory  = infraerrors.BadRequest("INVALID_TICKET_CATEGORY", "invalid ticket category")
	ErrInvalidTicketStatus    = infraerrors.BadRequest("INVALID_TICKET_STATUS", "invalid ticket status")
	ErrInvalidTicketPriority  = infraerrors.BadRequest("INVALID_TICKET_PRIORITY", "invalid ticket priority")
	ErrInvalidTicketRequestID = infraerrors.BadRequest("INVALID_TICKET_REQUEST_ID", "related request ID must be at most 128 characters")
	ErrInvalidTicketAssignee  = infraerrors.BadRequest("INVALID_TICKET_ASSIGNEE", "assignee must be an active administrator")
)

type Ticket struct {
	ID               int64      `json:"id"`
	UserID           int64      `json:"user_id"`
	Subject          string     `json:"subject"`
	Category         string     `json:"category"`
	Status           string     `json:"status"`
	Priority         string     `json:"priority"`
	RelatedRequestID string     `json:"related_request_id"`
	AssigneeID       *int64     `json:"assignee_id,omitempty"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	LastMessageAt    time.Time  `json:"last_message_at"`
	ClosedAt         *time.Time `json:"closed_at,omitempty"`
	Username         string     `json:"username,omitempty"`
	AssigneeUsername string     `json:"assignee_username,omitempty"`
	UnreadCount      int64      `json:"unread_count"`
}

type TicketMessage struct {
	ID         int64     `json:"id"`
	TicketID   int64     `json:"ticket_id"`
	SenderID   int64     `json:"sender_id"`
	SenderType string    `json:"sender_type"`
	Content    string    `json:"content"`
	RequestID  string    `json:"request_id"`
	CreatedAt  time.Time `json:"created_at"`
}

type TicketRead struct {
	TicketID          int64     `json:"ticket_id"`
	ReaderID          int64     `json:"reader_id"`
	LastReadMessageID int64     `json:"last_read_message_id"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type TicketDetail struct {
	Ticket   *Ticket         `json:"ticket"`
	Messages []TicketMessage `json:"messages"`
}

type TicketPage struct {
	Items    []Ticket `json:"items"`
	Total    int64    `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
}

type TicketListFilter struct {
	Status           string
	Category         string
	Priority         string
	UserID           int64
	AssigneeID       *int64
	RelatedRequestID string
}

type TicketAdminUpdate struct {
	Status     *string
	Priority   *string
	AssigneeID *int64
}

type TicketSummary struct {
	Total             int64 `json:"total"`
	PendingAdmin      int64 `json:"pending_admin"`
	PendingUser       int64 `json:"pending_user"`
	Resolved          int64 `json:"resolved"`
	Closed            int64 `json:"closed"`
	UnreadCount       int64 `json:"unread_count"`
	PendingAdminCount int64 `json:"pending_admin_count"`
	OpenCount         int64 `json:"open_count"`
}

type TicketAssignee struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

type TicketConfig struct {
	UserCenterEnabled bool `json:"user_center_enabled"`
	AcceptNewTickets  bool `json:"accept_new_tickets"`
}

type TicketRepository interface {
	Create(ctx context.Context, userID int64, subject, category, content, requestID string) (*TicketDetail, error)
	List(ctx context.Context, readerID int64, ownerID *int64, filter TicketListFilter, offset, limit int) ([]Ticket, int64, error)
	Get(ctx context.Context, ticketID, readerID int64, ownerID *int64) (*TicketDetail, error)
	AddMessage(ctx context.Context, senderID, ticketID int64, senderType, content, requestID string, ownerID *int64) (*TicketMessage, error)
	SetUserStatus(ctx context.Context, userID, ticketID int64, status string) (*Ticket, error)
	UpdateAdmin(ctx context.Context, ticketID int64, update TicketAdminUpdate) (*Ticket, error)
	MarkRead(ctx context.Context, readerID, ticketID int64, ownerID *int64) (*TicketRead, error)
	Summary(ctx context.Context, ownerID *int64, readerID int64, unreadSenderType string) (*TicketSummary, error)
	Assignees(ctx context.Context) ([]TicketAssignee, error)
}

type TicketService struct {
	repo        TicketRepository
	settingRepo SettingRepository
}

func NewTicketService(repo TicketRepository, settingRepo SettingRepository) *TicketService {
	return &TicketService{repo: repo, settingRepo: settingRepo}
}

func (s *TicketService) GetConfig(ctx context.Context) (*TicketConfig, error) {
	values, err := s.settingRepo.GetMultiple(ctx, []string{SettingKeyTicketUserCenterEnabled, SettingKeyTicketAcceptNewEnabled})
	if err != nil {
		return nil, fmt.Errorf("get ticket settings: %w", err)
	}
	return &TicketConfig{
		UserCenterEnabled: values[SettingKeyTicketUserCenterEnabled] != "false",
		AcceptNewTickets:  values[SettingKeyTicketAcceptNewEnabled] != "false",
	}, nil
}

func (s *TicketService) UpdateConfig(ctx context.Context, config TicketConfig) (*TicketConfig, error) {
	values := map[string]string{
		SettingKeyTicketUserCenterEnabled: fmt.Sprintf("%t", config.UserCenterEnabled),
		SettingKeyTicketAcceptNewEnabled:  fmt.Sprintf("%t", config.AcceptNewTickets),
	}
	if err := s.settingRepo.SetMultiple(ctx, values); err != nil {
		return nil, fmt.Errorf("update ticket settings: %w", err)
	}
	return &config, nil
}

func (s *TicketService) Create(ctx context.Context, userID int64, subject, category, content, requestID string) (*TicketDetail, error) {
	config, err := s.GetConfig(ctx)
	if err != nil {
		return nil, err
	}
	if !config.UserCenterEnabled {
		return nil, ErrTicketCenterDisabled
	}
	if !config.AcceptNewTickets {
		return nil, ErrTicketIntakeDisabled
	}
	subject, content, requestID, err = validateTicketCreate(subject, category, content, requestID)
	if err != nil {
		return nil, err
	}
	return s.repo.Create(ctx, userID, subject, strings.TrimSpace(category), content, requestID)
}

func (s *TicketService) ListUser(ctx context.Context, userID int64, filter TicketListFilter, page, pageSize int) (*TicketPage, error) {
	if err := s.requireUserCenter(ctx); err != nil {
		return nil, err
	}
	if err := validateTicketFilter(filter); err != nil {
		return nil, err
	}
	page, pageSize = normalizeTicketPage(page, pageSize)
	items, total, err := s.repo.List(ctx, userID, &userID, filter, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, err
	}
	redactTicketsForUser(items)
	return &TicketPage{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *TicketService) ListAdmin(ctx context.Context, readerID int64, filter TicketListFilter, page, pageSize int) (*TicketPage, error) {
	if err := validateTicketFilter(filter); err != nil {
		return nil, err
	}
	page, pageSize = normalizeTicketPage(page, pageSize)
	items, total, err := s.repo.List(ctx, readerID, nil, filter, (page-1)*pageSize, pageSize)
	if err != nil {
		return nil, err
	}
	return &TicketPage{Items: items, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *TicketService) GetUser(ctx context.Context, userID, ticketID int64) (*TicketDetail, error) {
	if err := s.requireUserCenter(ctx); err != nil {
		return nil, err
	}
	detail, err := s.repo.Get(ctx, ticketID, userID, &userID)
	if err != nil {
		return nil, mapTicketNotFound(err)
	}
	redactTicketForUser(detail.Ticket)
	for i := range detail.Messages {
		if detail.Messages[i].SenderType == TicketSenderAdmin {
			detail.Messages[i].SenderID = 0
		}
	}
	return detail, nil
}

func (s *TicketService) GetAdmin(ctx context.Context, readerID, ticketID int64) (*TicketDetail, error) {
	detail, err := s.repo.Get(ctx, ticketID, readerID, nil)
	if err != nil {
		return nil, mapTicketNotFound(err)
	}
	return detail, nil
}

func (s *TicketService) ReplyUser(ctx context.Context, userID, ticketID int64, content, requestID string) (*TicketMessage, error) {
	if err := s.requireUserCenter(ctx); err != nil {
		return nil, err
	}
	content, requestID, err := validateTicketMessage(content, requestID)
	if err != nil {
		return nil, err
	}
	message, err := s.repo.AddMessage(ctx, userID, ticketID, TicketSenderUser, content, requestID, &userID)
	return message, mapTicketNotFound(err)
}

func (s *TicketService) ReplyAdmin(ctx context.Context, adminID, ticketID int64, content, requestID string) (*TicketMessage, error) {
	content, requestID, err := validateTicketMessage(content, requestID)
	if err != nil {
		return nil, err
	}
	message, err := s.repo.AddMessage(ctx, adminID, ticketID, TicketSenderAdmin, content, requestID, nil)
	return message, mapTicketNotFound(err)
}

func (s *TicketService) SetUserStatus(ctx context.Context, userID, ticketID int64, status string) (*Ticket, error) {
	if err := s.requireUserCenter(ctx); err != nil {
		return nil, err
	}
	status = strings.TrimSpace(status)
	if status != TicketStatusResolved && status != TicketStatusClosed {
		return nil, ErrInvalidTicketStatus
	}
	ticket, err := s.repo.SetUserStatus(ctx, userID, ticketID, status)
	if err != nil {
		return nil, mapTicketNotFound(err)
	}
	redactTicketForUser(ticket)
	return ticket, nil
}

func (s *TicketService) UpdateAdmin(ctx context.Context, ticketID int64, update TicketAdminUpdate) (*Ticket, error) {
	if update.Status != nil {
		value := strings.TrimSpace(*update.Status)
		if !isValidTicketStatus(value) {
			return nil, ErrInvalidTicketStatus
		}
		update.Status = &value
	}
	if update.Priority != nil {
		value := strings.TrimSpace(*update.Priority)
		if !isValidTicketPriority(value) {
			return nil, ErrInvalidTicketPriority
		}
		update.Priority = &value
	}
	if update.AssigneeID != nil && *update.AssigneeID < 0 {
		return nil, ErrInvalidTicketAssignee
	}
	ticket, err := s.repo.UpdateAdmin(ctx, ticketID, update)
	return ticket, mapTicketNotFound(err)
}

func (s *TicketService) MarkUserRead(ctx context.Context, userID, ticketID int64) (*TicketRead, error) {
	if err := s.requireUserCenter(ctx); err != nil {
		return nil, err
	}
	read, err := s.repo.MarkRead(ctx, userID, ticketID, &userID)
	return read, mapTicketNotFound(err)
}

func (s *TicketService) MarkAdminRead(ctx context.Context, adminID, ticketID int64) (*TicketRead, error) {
	read, err := s.repo.MarkRead(ctx, adminID, ticketID, nil)
	return read, mapTicketNotFound(err)
}

func (s *TicketService) UserSummary(ctx context.Context, userID int64) (*TicketSummary, error) {
	if err := s.requireUserCenter(ctx); err != nil {
		return nil, err
	}
	return s.repo.Summary(ctx, &userID, userID, TicketSenderAdmin)
}

func (s *TicketService) AdminSummary(ctx context.Context, adminID int64) (*TicketSummary, error) {
	return s.repo.Summary(ctx, nil, adminID, TicketSenderUser)
}

func (s *TicketService) Assignees(ctx context.Context) ([]TicketAssignee, error) {
	return s.repo.Assignees(ctx)
}

func (s *TicketService) requireUserCenter(ctx context.Context) error {
	config, err := s.GetConfig(ctx)
	if err != nil {
		return err
	}
	if !config.UserCenterEnabled {
		return ErrTicketCenterDisabled
	}
	return nil
}

func validateTicketCreate(subject, category, content, requestID string) (string, string, string, error) {
	subject = strings.TrimSpace(subject)
	if subject == "" || utf8.RuneCountInString(subject) > TicketSubjectMaxLength {
		return "", "", "", ErrInvalidTicketSubject
	}
	if !isValidTicketCategory(strings.TrimSpace(category)) {
		return "", "", "", ErrInvalidTicketCategory
	}
	content, requestID, err := validateTicketMessage(content, requestID)
	return subject, content, requestID, err
}

func validateTicketMessage(content, requestID string) (string, string, error) {
	content = strings.TrimSpace(content)
	if content == "" || utf8.RuneCountInString(content) > TicketContentMaxLength {
		return "", "", ErrInvalidTicketContent
	}
	requestID = strings.TrimSpace(requestID)
	if utf8.RuneCountInString(requestID) > TicketRequestIDMaxLength {
		return "", "", ErrInvalidTicketRequestID
	}
	return content, requestID, nil
}

func validateTicketFilter(filter TicketListFilter) error {
	if filter.Status != "" && !isValidTicketStatus(filter.Status) {
		return ErrInvalidTicketStatus
	}
	if filter.Category != "" && !isValidTicketCategory(filter.Category) {
		return ErrInvalidTicketCategory
	}
	if filter.Priority != "" && !isValidTicketPriority(filter.Priority) {
		return ErrInvalidTicketPriority
	}
	if utf8.RuneCountInString(strings.TrimSpace(filter.RelatedRequestID)) > TicketRequestIDMaxLength {
		return ErrInvalidTicketRequestID
	}
	return nil
}

func isValidTicketStatus(value string) bool {
	switch value {
	case TicketStatusPendingAdmin, TicketStatusPendingUser, TicketStatusResolved, TicketStatusClosed:
		return true
	default:
		return false
	}
}

func isValidTicketCategory(value string) bool {
	switch value {
	case TicketCategoryAccount, TicketCategoryBilling, TicketCategoryAPIModel, TicketCategoryIncident, TicketCategoryFeature, TicketCategoryOther:
		return true
	default:
		return false
	}
}

func isValidTicketPriority(value string) bool {
	switch value {
	case TicketPriorityLow, TicketPriorityNormal, TicketPriorityHigh, TicketPriorityUrgent:
		return true
	default:
		return false
	}
}

func normalizeTicketPage(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func mapTicketNotFound(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ErrTicketNotFound
	}
	return err
}

func redactTicketsForUser(items []Ticket) {
	for i := range items {
		redactTicketForUser(&items[i])
	}
}

func redactTicketForUser(ticket *Ticket) {
	if ticket == nil {
		return
	}
	ticket.AssigneeID = nil
	ticket.AssigneeUsername = ""
}
