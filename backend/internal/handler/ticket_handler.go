package handler

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

type createTicketRequest struct {
	Subject          string `json:"subject" binding:"required"`
	Category         string `json:"category" binding:"required"`
	Content          string `json:"content" binding:"required"`
	RelatedRequestID string `json:"related_request_id"`
}

type replyTicketRequest struct {
	Content   string `json:"content" binding:"required"`
	RequestID string `json:"request_id"`
}

type updateUserTicketStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=resolved closed"`
}

// Config is intentionally reachable while the feature is disabled so the UI
// can remove the entry and redirect cleanly.
func (h *TicketHandler) Config(c *gin.Context) {
	config, err := h.service.GetConfig(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, config)
}

func (h *TicketHandler) List(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	page, pageSize := response.ParsePagination(c)
	filter := service.TicketListFilter{
		Status:   strings.TrimSpace(c.Query("status")),
		Category: strings.TrimSpace(c.Query("category")),
	}
	result, err := h.service.ListUser(c.Request.Context(), subject.UserID, filter, page, pageSize)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

func (h *TicketHandler) Create(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	var req createTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	detail, err := h.service.Create(c.Request.Context(), subject.UserID, req.Subject, req.Category, req.Content, req.RelatedRequestID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, detail)
}

func (h *TicketHandler) Summary(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	summary, err := h.service.UserSummary(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, summary)
}

func (h *TicketHandler) Get(c *gin.Context) {
	subject, ticketID, ok := userAndTicketID(c)
	if !ok {
		return
	}
	detail, err := h.service.GetUser(c.Request.Context(), subject.UserID, ticketID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, detail)
}

func (h *TicketHandler) MarkRead(c *gin.Context) {
	subject, ticketID, ok := userAndTicketID(c)
	if !ok {
		return
	}
	read, err := h.service.MarkUserRead(c.Request.Context(), subject.UserID, ticketID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, read)
}

func (h *TicketHandler) Reply(c *gin.Context) {
	subject, ticketID, ok := userAndTicketID(c)
	if !ok {
		return
	}
	var req replyTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	message, err := h.service.ReplyUser(c.Request.Context(), subject.UserID, ticketID, req.Content, req.RequestID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, message)
}

func (h *TicketHandler) UpdateStatus(c *gin.Context) {
	subject, ticketID, ok := userAndTicketID(c)
	if !ok {
		return
	}
	var req updateUserTicketStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	ticket, err := h.service.SetUserStatus(c.Request.Context(), subject.UserID, ticketID, req.Status)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, ticket)
}

func userAndTicketID(c *gin.Context) (middleware2.AuthSubject, int64, bool) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return middleware2.AuthSubject{}, 0, false
	}
	ticketID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || ticketID <= 0 {
		response.BadRequest(c, "Invalid ticket ID")
		return middleware2.AuthSubject{}, 0, false
	}
	return subject, ticketID, true
}
