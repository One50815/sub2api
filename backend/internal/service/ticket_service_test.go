package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type ticketSettingRepositoryStub struct {
	values map[string]string
}

func (s *ticketSettingRepositoryStub) Get(context.Context, string) (*Setting, error) {
	return nil, ErrSettingNotFound
}
func (s *ticketSettingRepositoryStub) GetValue(_ context.Context, key string) (string, error) {
	value, ok := s.values[key]
	if !ok {
		return "", ErrSettingNotFound
	}
	return value, nil
}
func (s *ticketSettingRepositoryStub) Set(_ context.Context, key, value string) error {
	if s.values == nil {
		s.values = map[string]string{}
	}
	s.values[key] = value
	return nil
}
func (s *ticketSettingRepositoryStub) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	result := map[string]string{}
	for _, key := range keys {
		if value, ok := s.values[key]; ok {
			result[key] = value
		}
	}
	return result, nil
}
func (s *ticketSettingRepositoryStub) SetMultiple(_ context.Context, values map[string]string) error {
	if s.values == nil {
		s.values = map[string]string{}
	}
	for key, value := range values {
		s.values[key] = value
	}
	return nil
}
func (s *ticketSettingRepositoryStub) GetAll(context.Context) (map[string]string, error) {
	return s.values, nil
}
func (s *ticketSettingRepositoryStub) Delete(_ context.Context, key string) error {
	delete(s.values, key)
	return nil
}

type ticketRepositoryStub struct {
	createCalled bool
	listCalled   bool
	created      *TicketDetail
}

func (s *ticketRepositoryStub) Create(_ context.Context, userID int64, subject, category, content, requestID string) (*TicketDetail, error) {
	s.createCalled = true
	if s.created != nil {
		return s.created, nil
	}
	return &TicketDetail{Ticket: &Ticket{ID: 1, UserID: userID, Subject: subject, Category: category, RelatedRequestID: requestID}, Messages: []TicketMessage{{Content: content}}}, nil
}
func (s *ticketRepositoryStub) List(context.Context, int64, *int64, TicketListFilter, int, int) ([]Ticket, int64, error) {
	s.listCalled = true
	return []Ticket{}, 0, nil
}
func (s *ticketRepositoryStub) Get(context.Context, int64, int64, *int64) (*TicketDetail, error) {
	return nil, ErrTicketNotFound
}
func (s *ticketRepositoryStub) AddMessage(context.Context, int64, int64, string, string, string, *int64) (*TicketMessage, error) {
	return nil, nil
}
func (s *ticketRepositoryStub) SetUserStatus(context.Context, int64, int64, string) (*Ticket, error) {
	return nil, nil
}
func (s *ticketRepositoryStub) UpdateAdmin(context.Context, int64, TicketAdminUpdate) (*Ticket, error) {
	return nil, nil
}
func (s *ticketRepositoryStub) MarkRead(context.Context, int64, int64, *int64) (*TicketRead, error) {
	return nil, nil
}
func (s *ticketRepositoryStub) Summary(context.Context, *int64, int64, string) (*TicketSummary, error) {
	return &TicketSummary{}, nil
}
func (s *ticketRepositoryStub) Assignees(context.Context) ([]TicketAssignee, error) {
	return []TicketAssignee{}, nil
}

func TestTicketConfigDefaultsAndPersistence(t *testing.T) {
	settings := &ticketSettingRepositoryStub{values: map[string]string{}}
	svc := NewTicketService(&ticketRepositoryStub{}, settings)

	config, err := svc.GetConfig(context.Background())
	require.NoError(t, err)
	require.True(t, config.UserCenterEnabled)
	require.True(t, config.AcceptNewTickets)

	updated, err := svc.UpdateConfig(context.Background(), TicketConfig{UserCenterEnabled: false, AcceptNewTickets: true})
	require.NoError(t, err)
	require.False(t, updated.UserCenterEnabled)
	require.Equal(t, "false", settings.values[SettingKeyTicketUserCenterEnabled])
	require.Equal(t, "true", settings.values[SettingKeyTicketAcceptNewEnabled])
}

func TestTicketUserCenterAndIntakeAreIndependentGates(t *testing.T) {
	ctx := context.Background()
	repo := &ticketRepositoryStub{}
	settings := &ticketSettingRepositoryStub{values: map[string]string{
		SettingKeyTicketUserCenterEnabled: "false",
		SettingKeyTicketAcceptNewEnabled:  "true",
	}}
	svc := NewTicketService(repo, settings)

	_, err := svc.Create(ctx, 7, "Subject", TicketCategoryOther, "Body", "")
	require.ErrorIs(t, err, ErrTicketCenterDisabled)
	require.False(t, repo.createCalled)

	settings.values[SettingKeyTicketUserCenterEnabled] = "true"
	settings.values[SettingKeyTicketAcceptNewEnabled] = "false"
	_, err = svc.Create(ctx, 7, "Subject", TicketCategoryOther, "Body", "")
	require.ErrorIs(t, err, ErrTicketIntakeDisabled)
	require.False(t, repo.createCalled)

	_, err = svc.ListUser(ctx, 7, TicketListFilter{}, 1, 20)
	require.NoError(t, err)
	require.True(t, repo.listCalled, "history remains readable while new-ticket intake is disabled")
}

func TestTicketCreateValidationAndUserStatusBoundary(t *testing.T) {
	ctx := context.Background()
	repo := &ticketRepositoryStub{}
	settings := &ticketSettingRepositoryStub{values: map[string]string{}}
	svc := NewTicketService(repo, settings)

	detail, err := svc.Create(ctx, 9, "  Login issue  ", TicketCategoryAccount, "  Cannot sign in  ", " req-123 ")
	require.NoError(t, err)
	require.Equal(t, "Login issue", detail.Ticket.Subject)
	require.Equal(t, "req-123", detail.Ticket.RelatedRequestID)
	require.Equal(t, "Cannot sign in", detail.Messages[0].Content)

	_, err = svc.SetUserStatus(ctx, 9, 1, TicketStatusPendingUser)
	require.True(t, errors.Is(err, ErrInvalidTicketStatus))
}
