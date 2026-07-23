package admin

import (
	"strconv"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	service *service.TicketService
}

func NewTicketHandler(ticketService *service.TicketService) *TicketHandler {
	return &TicketHandler{service: ticketService}
}

type updateTicketConfigRequest struct {
	UserCenterEnabled bool `json:"user_center_enabled"`
	AcceptNewTickets  bool `json:"accept_new_tickets"`
}

type adminReplyTicketRequest struct {
	Content   string `json:"content" binding:"required"`
	RequestID string `json:"request_id"`
}

type adminUpdateTicketRequest struct {
	Status     *string `json:"status" binding:"omitempty,oneof=pending_admin pending_user resolved closed"`
	Priority   *string `json:"priority" binding:"omitempty,oneof=low normal high urgent"`
	AssigneeID *int64  `json:"assignee_id"`
}

func (h *TicketHandler) GetConfig(c *gin.Context) {
	config, err := h.service.GetConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, config)
}

func (h *TicketHandler) UpdateConfig(c *gin.Context) {
	var req updateTicketConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	config, err := h.service.UpdateConfig(c.Request.Context(), service.TicketConfig{
		UserCenterEnabled: req.UserCenterEnabled,
		AcceptNewTickets:  req.AcceptNewTickets,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, config)
}

func (h *TicketHandler) List(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "Administrator not found in context")
		return
	}
	page, pageSize := response.ParsePagination(c)
	filter := service.TicketListFilter{
		Status:           strings.TrimSpace(c.Query("status")),
		Category:         strings.TrimSpace(c.Query("category")),
		Priority:         strings.TrimSpace(c.Query("priority")),
		RelatedRequestID: strings.TrimSpace(c.Query("related_request_id")),
	}
	if value := strings.TrimSpace(c.Query("user_id")); value != "" {
		filter.UserID, _ = strconv.ParseInt(value, 10, 64)
	}
	if value := strings.TrimSpace(c.Query("assignee_id")); value != "" {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id < 0 {
			response.BadRequest(c, "Invalid assignee ID")
			return
		}
		filter.AssigneeID = &id
	}
	result, err := h.service.ListAdmin(c.Request.Context(), subject.UserID, filter, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

func (h *TicketHandler) Summary(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "Administrator not found in context")
		return
	}
	summary, err := h.service.AdminSummary(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, summary)
}

func (h *TicketHandler) Assignees(c *gin.Context) {
	items, err := h.service.Assignees(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, items)
}

func (h *TicketHandler) Get(c *gin.Context) {
	subject, ticketID, ok := adminAndTicketID(c)
	if !ok {
		return
	}
	detail, err := h.service.GetAdmin(c.Request.Context(), subject.UserID, ticketID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, detail)
}

func (h *TicketHandler) MarkRead(c *gin.Context) {
	subject, ticketID, ok := adminAndTicketID(c)
	if !ok {
		return
	}
	read, err := h.service.MarkAdminRead(c.Request.Context(), subject.UserID, ticketID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, read)
}

func (h *TicketHandler) Reply(c *gin.Context) {
	subject, ticketID, ok := adminAndTicketID(c)
	if !ok {
		return
	}
	var req adminReplyTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	message, err := h.service.ReplyAdmin(c.Request.Context(), subject.UserID, ticketID, req.Content, req.RequestID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, message)
}

func (h *TicketHandler) Update(c *gin.Context) {
	_, ticketID, ok := adminAndTicketID(c)
	if !ok {
		return
	}
	var req adminUpdateTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	ticket, err := h.service.UpdateAdmin(c.Request.Context(), ticketID, service.TicketAdminUpdate{
		Status: req.Status, Priority: req.Priority, AssigneeID: req.AssigneeID,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, ticket)
}

func adminAndTicketID(c *gin.Context) (middleware2.AuthSubject, int64, bool) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "Administrator not found in context")
		return middleware2.AuthSubject{}, 0, false
	}
	ticketID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || ticketID <= 0 {
		response.BadRequest(c, "Invalid ticket ID")
		return middleware2.AuthSubject{}, 0, false
	}
	return subject, ticketID, true
}
