package handler

import (
	"strconv"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	middleware2 "github.com/Wei-Shaw/sub2api/internal/server/middleware"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	service *service.MembershipService
}

func NewMembershipHandler(membershipService *service.MembershipService) *MembershipHandler {
	return &MembershipHandler{service: membershipService}
}

func (h *MembershipHandler) Summary(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	result, err := h.service.GetSummary(c.Request.Context(), subject.UserID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

func (h *MembershipHandler) Offers(c *gin.Context) {
	items, err := h.service.ListOffers(c.Request.Context(), true)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, items)
}

func (h *MembershipHandler) SetOveragePolicy(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "User not found in context")
		return
	}
	entitlementID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || entitlementID <= 0 {
		response.BadRequest(c, "Invalid entitlement ID")
		return
	}
	var req struct {
		Policy string `json:"policy" binding:"required,oneof=block wallet"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if err := h.service.SetOveragePolicy(c.Request.Context(), subject.UserID, entitlementID, req.Policy); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"updated": true})
}
