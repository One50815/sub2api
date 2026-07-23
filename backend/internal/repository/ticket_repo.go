package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type ticketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) service.TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, userID int64, subject, category, content, requestID string) (*service.TicketDetail, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	ticket := &service.Ticket{}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO tickets (user_id, subject, category, status, priority, related_request_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, subject, category, status, priority, related_request_id,
		          assignee_id, created_at, updated_at, last_message_at, closed_at`,
		userID, subject, category, service.TicketStatusPendingAdmin, service.TicketPriorityNormal, requestID,
	).Scan(
		&ticket.ID, &ticket.UserID, &ticket.Subject, &ticket.Category, &ticket.Status,
		&ticket.Priority, &ticket.RelatedRequestID, &ticket.AssigneeID, &ticket.CreatedAt,
		&ticket.UpdatedAt, &ticket.LastMessageAt, &ticket.ClosedAt,
	)
	if err != nil {
		return nil, err
	}

	message := service.TicketMessage{}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO ticket_messages (ticket_id, sender_id, sender_type, content, request_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, ticket_id, sender_id, sender_type, content, request_id, created_at`,
		ticket.ID, userID, service.TicketSenderUser, content, requestID,
	).Scan(&message.ID, &message.TicketID, &message.SenderID, &message.SenderType, &message.Content, &message.RequestID, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	if err := upsertTicketRead(ctx, tx, ticket.ID, userID, message.ID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &service.TicketDetail{Ticket: ticket, Messages: []service.TicketMessage{message}}, nil
}

func (r *ticketRepository) List(ctx context.Context, readerID int64, ownerID *int64, filter service.TicketListFilter, offset, limit int) ([]service.Ticket, int64, error) {
	where, filterArgs := ticketFilterSQL(ownerID, filter, 1)
	var total int64
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM tickets t"+where, filterArgs...).Scan(&total); err != nil {
		return nil, 0, err
	}

	unreadSender := service.TicketSenderUser
	if ownerID != nil {
		unreadSender = service.TicketSenderAdmin
	}
	where, filterArgs = ticketFilterSQL(ownerID, filter, 3)
	args := []any{readerID, unreadSender}
	args = append(args, filterArgs...)
	args = append(args, limit, offset)
	limitPlaceholder := len(args) - 1
	offsetPlaceholder := len(args)
	query := fmt.Sprintf(`
		SELECT t.id, t.user_id, t.subject, t.category, t.status, t.priority,
		       t.related_request_id, t.assignee_id, t.created_at, t.updated_at,
		       t.last_message_at, t.closed_at,
		       COALESCE(NULLIF(u.username, ''), u.email, ''),
		       COALESCE(NULLIF(a.username, ''), a.email, ''),
		       COALESCE((
		         SELECT COUNT(*)
		         FROM ticket_messages tm
		         LEFT JOIN ticket_reads tr ON tr.ticket_id = tm.ticket_id AND tr.reader_id = $1
		         WHERE tm.ticket_id = t.id AND tm.sender_type = $2
		           AND tm.id > COALESCE(tr.last_read_message_id, 0)
		       ), 0)
		FROM tickets t
		LEFT JOIN users u ON u.id = t.user_id
		LEFT JOIN users a ON a.id = t.assignee_id
		%s
		ORDER BY t.last_message_at DESC, t.id DESC
		LIMIT $%d OFFSET $%d`, where, limitPlaceholder, offsetPlaceholder)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	items := make([]service.Ticket, 0)
	for rows.Next() {
		item, err := scanTicket(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *ticketRepository) Get(ctx context.Context, ticketID, readerID int64, ownerID *int64) (*service.TicketDetail, error) {
	unreadSender := service.TicketSenderUser
	ownerClause := ""
	args := []any{readerID, unreadSender, ticketID}
	if ownerID != nil {
		unreadSender = service.TicketSenderAdmin
		args[1] = unreadSender
		args = append(args, *ownerID)
		ownerClause = " AND t.user_id = $4"
	}
	row := r.db.QueryRowContext(ctx, `
		SELECT t.id, t.user_id, t.subject, t.category, t.status, t.priority,
		       t.related_request_id, t.assignee_id, t.created_at, t.updated_at,
		       t.last_message_at, t.closed_at,
		       COALESCE(NULLIF(u.username, ''), u.email, ''),
		       COALESCE(NULLIF(a.username, ''), a.email, ''),
		       COALESCE((
		         SELECT COUNT(*)
		         FROM ticket_messages tm
		         LEFT JOIN ticket_reads tr ON tr.ticket_id = tm.ticket_id AND tr.reader_id = $1
		         WHERE tm.ticket_id = t.id AND tm.sender_type = $2
		           AND tm.id > COALESCE(tr.last_read_message_id, 0)
		       ), 0)
		FROM tickets t
		LEFT JOIN users u ON u.id = t.user_id
		LEFT JOIN users a ON a.id = t.assignee_id
		WHERE t.id = $3`+ownerClause, args...)
	ticket, err := scanTicket(row)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, ticket_id, sender_id, sender_type, content, request_id, created_at
		FROM ticket_messages WHERE ticket_id = $1 ORDER BY id ASC`, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	messages := make([]service.TicketMessage, 0)
	for rows.Next() {
		var message service.TicketMessage
		if err := rows.Scan(&message.ID, &message.TicketID, &message.SenderID, &message.SenderType, &message.Content, &message.RequestID, &message.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return &service.TicketDetail{Ticket: ticket, Messages: messages}, rows.Err()
}

func (r *ticketRepository) AddMessage(ctx context.Context, senderID, ticketID int64, senderType, content, requestID string, ownerID *int64) (*service.TicketMessage, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()

	query := "SELECT status FROM tickets WHERE id = $1"
	args := []any{ticketID}
	if ownerID != nil {
		query += " AND user_id = $2"
		args = append(args, *ownerID)
	}
	query += " FOR UPDATE"
	var status string
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&status); err != nil {
		return nil, err
	}
	if status == service.TicketStatusClosed {
		return nil, service.ErrTicketClosed
	}

	message := &service.TicketMessage{}
	err = tx.QueryRowContext(ctx, `
		INSERT INTO ticket_messages (ticket_id, sender_id, sender_type, content, request_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, ticket_id, sender_id, sender_type, content, request_id, created_at`,
		ticketID, senderID, senderType, content, requestID,
	).Scan(&message.ID, &message.TicketID, &message.SenderID, &message.SenderType, &message.Content, &message.RequestID, &message.CreatedAt)
	if err != nil {
		return nil, err
	}
	nextStatus := service.TicketStatusPendingUser
	if senderType == service.TicketSenderUser {
		nextStatus = service.TicketStatusPendingAdmin
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE tickets SET status = $1, updated_at = NOW(), last_message_at = NOW(), closed_at = NULL
		WHERE id = $2`, nextStatus, ticketID); err != nil {
		return nil, err
	}
	if err := upsertTicketRead(ctx, tx, ticketID, senderID, message.ID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return message, nil
}

func (r *ticketRepository) SetUserStatus(ctx context.Context, userID, ticketID int64, status string) (*service.Ticket, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	var current string
	if err := tx.QueryRowContext(ctx, "SELECT status FROM tickets WHERE id = $1 AND user_id = $2 FOR UPDATE", ticketID, userID).Scan(&current); err != nil {
		return nil, err
	}
	if current == service.TicketStatusClosed {
		return nil, service.ErrTicketClosed
	}
	if _, err := tx.ExecContext(ctx, `
		UPDATE tickets
		SET status = $1, updated_at = NOW(), closed_at = CASE WHEN $1 = 'closed' THEN NOW() ELSE NULL END
		WHERE id = $2 AND user_id = $3`, status, ticketID, userID); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return r.getTicketWithoutUnread(ctx, ticketID)
}

func (r *ticketRepository) UpdateAdmin(ctx context.Context, ticketID int64, update service.TicketAdminUpdate) (*service.Ticket, error) {
	if update.AssigneeID != nil && *update.AssigneeID > 0 {
		var exists bool
		if err := r.db.QueryRowContext(ctx, `
			SELECT EXISTS(SELECT 1 FROM users WHERE id = $1 AND role = 'admin' AND status = 'active' AND deleted_at IS NULL)`,
			*update.AssigneeID,
		).Scan(&exists); err != nil {
			return nil, err
		}
		if !exists {
			return nil, service.ErrInvalidTicketAssignee
		}
	}

	sets := make([]string, 0, 5)
	args := make([]any, 0, 4)
	add := func(expression string, value any) {
		args = append(args, value)
		sets = append(sets, fmt.Sprintf(expression, len(args)))
	}
	if update.Status != nil {
		add("status = $%d", *update.Status)
		if *update.Status == service.TicketStatusClosed {
			sets = append(sets, "closed_at = NOW()")
		} else {
			sets = append(sets, "closed_at = NULL")
		}
	}
	if update.Priority != nil {
		add("priority = $%d", *update.Priority)
	}
	if update.AssigneeID != nil {
		if *update.AssigneeID == 0 {
			sets = append(sets, "assignee_id = NULL")
		} else {
			add("assignee_id = $%d", *update.AssigneeID)
		}
	}
	if len(sets) > 0 {
		sets = append(sets, "updated_at = NOW()")
		args = append(args, ticketID)
		result, err := r.db.ExecContext(ctx, fmt.Sprintf("UPDATE tickets SET %s WHERE id = $%d", strings.Join(sets, ", "), len(args)), args...)
		if err != nil {
			return nil, err
		}
		if count, err := result.RowsAffected(); err == nil && count == 0 {
			return nil, sql.ErrNoRows
		}
	}
	return r.getTicketWithoutUnread(ctx, ticketID)
}

func (r *ticketRepository) MarkRead(ctx context.Context, readerID, ticketID int64, ownerID *int64) (*service.TicketRead, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() { _ = tx.Rollback() }()
	query := "SELECT id FROM tickets WHERE id = $1"
	args := []any{ticketID}
	if ownerID != nil {
		query += " AND user_id = $2"
		args = append(args, *ownerID)
	}
	var exists int64
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&exists); err != nil {
		return nil, err
	}
	var lastMessageID int64
	if err := tx.QueryRowContext(ctx, "SELECT COALESCE(MAX(id), 0) FROM ticket_messages WHERE ticket_id = $1", ticketID).Scan(&lastMessageID); err != nil {
		return nil, err
	}
	if err := upsertTicketRead(ctx, tx, ticketID, readerID, lastMessageID); err != nil {
		return nil, err
	}
	read := &service.TicketRead{}
	if err := tx.QueryRowContext(ctx, `
		SELECT ticket_id, reader_id, last_read_message_id, updated_at
		FROM ticket_reads WHERE ticket_id = $1 AND reader_id = $2`, ticketID, readerID,
	).Scan(&read.TicketID, &read.ReaderID, &read.LastReadMessageID, &read.UpdatedAt); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return read, nil
}

func (r *ticketRepository) Summary(ctx context.Context, ownerID *int64, readerID int64, unreadSenderType string) (*service.TicketSummary, error) {
	where := ""
	args := []any{}
	if ownerID != nil {
		where = " WHERE user_id = $1"
		args = append(args, *ownerID)
	}
	rows, err := r.db.QueryContext(ctx, "SELECT status, COUNT(*) FROM tickets"+where+" GROUP BY status", args...)
	if err != nil {
		return nil, err
	}
	summary := &service.TicketSummary{}
	for rows.Next() {
		var status string
		var count int64
		if err := rows.Scan(&status, &count); err != nil {
			rows.Close()
			return nil, err
		}
		summary.Total += count
		switch status {
		case service.TicketStatusPendingAdmin:
			summary.PendingAdmin = count
		case service.TicketStatusPendingUser:
			summary.PendingUser = count
		case service.TicketStatusResolved:
			summary.Resolved = count
		case service.TicketStatusClosed:
			summary.Closed = count
		}
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	summary.PendingAdminCount = summary.PendingAdmin
	summary.OpenCount = summary.Total - summary.Closed

	unreadArgs := []any{readerID, unreadSenderType}
	unreadOwner := ""
	if ownerID != nil {
		unreadArgs = append(unreadArgs, *ownerID)
		unreadOwner = " AND t.user_id = $3"
	}
	err = r.db.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM ticket_messages tm
		JOIN tickets t ON t.id = tm.ticket_id
		LEFT JOIN ticket_reads tr ON tr.ticket_id = tm.ticket_id AND tr.reader_id = $1
		WHERE tm.sender_type = $2 AND tm.id > COALESCE(tr.last_read_message_id, 0)`+unreadOwner,
		unreadArgs...,
	).Scan(&summary.UnreadCount)
	return summary, err
}

func (r *ticketRepository) Assignees(ctx context.Context) ([]service.TicketAssignee, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, COALESCE(NULLIF(username, ''), email)
		FROM users WHERE role = 'admin' AND status = 'active' AND deleted_at IS NULL
		ORDER BY id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]service.TicketAssignee, 0)
	for rows.Next() {
		var item service.TicketAssignee
		if err := rows.Scan(&item.ID, &item.Username); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *ticketRepository) getTicketWithoutUnread(ctx context.Context, ticketID int64) (*service.Ticket, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT t.id, t.user_id, t.subject, t.category, t.status, t.priority,
		       t.related_request_id, t.assignee_id, t.created_at, t.updated_at,
		       t.last_message_at, t.closed_at,
		       COALESCE(NULLIF(u.username, ''), u.email, ''),
		       COALESCE(NULLIF(a.username, ''), a.email, ''), 0
		FROM tickets t
		LEFT JOIN users u ON u.id = t.user_id
		LEFT JOIN users a ON a.id = t.assignee_id
		WHERE t.id = $1`, ticketID)
	return scanTicket(row)
}

func upsertTicketRead(ctx context.Context, tx *sql.Tx, ticketID, readerID, messageID int64) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO ticket_reads (ticket_id, reader_id, last_read_message_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (ticket_id, reader_id)
		DO UPDATE SET last_read_message_id = EXCLUDED.last_read_message_id, updated_at = NOW()`,
		ticketID, readerID, messageID)
	return err
}

type ticketScanner interface {
	Scan(dest ...any) error
}

func scanTicket(scanner ticketScanner) (*service.Ticket, error) {
	ticket := &service.Ticket{}
	var assigneeID sql.NullInt64
	var closedAt sql.NullTime
	if err := scanner.Scan(
		&ticket.ID, &ticket.UserID, &ticket.Subject, &ticket.Category, &ticket.Status,
		&ticket.Priority, &ticket.RelatedRequestID, &assigneeID, &ticket.CreatedAt,
		&ticket.UpdatedAt, &ticket.LastMessageAt, &closedAt, &ticket.Username,
		&ticket.AssigneeUsername, &ticket.UnreadCount,
	); err != nil {
		return nil, err
	}
	if assigneeID.Valid {
		value := assigneeID.Int64
		ticket.AssigneeID = &value
	}
	if closedAt.Valid {
		value := closedAt.Time
		ticket.ClosedAt = &value
	}
	return ticket, nil
}

func ticketFilterSQL(ownerID *int64, filter service.TicketListFilter, start int) (string, []any) {
	clauses := make([]string, 0, 7)
	args := make([]any, 0, 7)
	add := func(expression string, value any) {
		args = append(args, value)
		clauses = append(clauses, fmt.Sprintf(expression, start+len(args)-1))
	}
	if ownerID != nil {
		add("t.user_id = $%d", *ownerID)
	}
	if filter.Status != "" {
		add("t.status = $%d", filter.Status)
	}
	if filter.Category != "" {
		add("t.category = $%d", filter.Category)
	}
	if filter.Priority != "" {
		add("t.priority = $%d", filter.Priority)
	}
	if filter.UserID > 0 {
		add("t.user_id = $%d", filter.UserID)
	}
	if filter.AssigneeID != nil {
		if *filter.AssigneeID == 0 {
			clauses = append(clauses, "t.assignee_id IS NULL")
		} else {
			add("t.assignee_id = $%d", *filter.AssigneeID)
		}
	}
	if filter.RelatedRequestID != "" {
		add("t.related_request_id ILIKE '%%' || $%d || '%%'", strings.TrimSpace(filter.RelatedRequestID))
	}
	if len(clauses) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(clauses, " AND "), args
}
