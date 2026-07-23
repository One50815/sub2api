package admin

import (
	"fmt"
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

func (h *MembershipHandler) Catalog(c *gin.Context) {
	result, err := h.service.GetCatalog(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, result)
}

func (h *MembershipHandler) GetGroupSetting(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("groupId"), 10, 64)
	if err != nil || groupID <= 0 {
		response.BadRequest(c, "Invalid group ID")
		return
	}
	setting, err := h.service.GetOmnioProGroupSetting(c.Request.Context(), groupID)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, setting)
}

func (h *MembershipHandler) UpsertGroupSetting(c *gin.Context) {
	groupID, err := strconv.ParseInt(c.Param("groupId"), 10, 64)
	if err != nil || groupID <= 0 {
		response.BadRequest(c, "Invalid group ID")
		return
	}
	var req struct {
		RateMultiplier *float64 `json:"rate_multiplier"`
		ProOnly        bool     `json:"pro_only"`
		DailyFreeUSD   *float64 `json:"daily_free_usd"`
		MonthlyFreeUSD *float64 `json:"monthly_free_usd"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	setting, err := h.service.UpsertOmnioProGroupSetting(c.Request.Context(), service.OmnioProGroupSetting{
		GroupID:        groupID,
		RateMultiplier: req.RateMultiplier,
		ProOnly:        req.ProOnly,
		DailyFreeUSD:   req.DailyFreeUSD,
		MonthlyFreeUSD: req.MonthlyFreeUSD,
	})
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, setting)
}

func (h *MembershipHandler) UpsertLevel(c *gin.Context) {
	var req service.MembershipLevel
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if value := c.Param("id"); value != "" {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			response.BadRequest(c, "Invalid level ID")
			return
		}
		req.ID = id
	}
	level, err := h.service.UpsertLevel(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, level)
}

func (h *MembershipHandler) DeleteLevel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid level ID")
		return
	}
	if err := h.service.DeleteLevel(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *MembershipHandler) UpsertOffer(c *gin.Context) {
	var req service.MembershipOffer
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if value := c.Param("id"); value != "" {
		id, err := strconv.ParseInt(value, 10, 64)
		if err != nil || id <= 0 {
			response.BadRequest(c, "Invalid offer ID")
			return
		}
		req.ID = id
	}
	offer, err := h.service.UpsertOffer(c.Request.Context(), req)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, offer)
}

func (h *MembershipHandler) DeleteOffer(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid offer ID")
		return
	}
	if err := h.service.DeleteOffer(c.Request.Context(), id); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *MembershipHandler) UpsertBenefit(c *gin.Context) {
	var req service.MembershipGroupBenefit
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if err := h.service.UpsertGroupBenefit(c.Request.Context(), req); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"updated": true})
}

func (h *MembershipHandler) DeleteBenefit(c *gin.Context) {
	levelID, err1 := strconv.ParseInt(c.Param("levelId"), 10, 64)
	groupID, err2 := strconv.ParseInt(c.Param("groupId"), 10, 64)
	if err1 != nil || err2 != nil || levelID <= 0 || groupID <= 0 {
		response.BadRequest(c, "Invalid benefit ID")
		return
	}
	if err := h.service.DeleteGroupBenefit(c.Request.Context(), levelID, groupID); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *MembershipHandler) SetPlanBenefit(c *gin.Context) {
	var req service.PlanMembershipBenefit
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if err := h.service.SetPlanBenefit(c.Request.Context(), req); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"updated": true})
}

func (h *MembershipHandler) DeletePlanBenefit(c *gin.Context) {
	planID, err := strconv.ParseInt(c.Param("planId"), 10, 64)
	if err != nil || planID <= 0 {
		response.BadRequest(c, "Invalid plan ID")
		return
	}
	if err := h.service.DeletePlanBenefit(c.Request.Context(), planID); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"deleted": true})
}

func (h *MembershipHandler) Grant(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "Administrator not found in context")
		return
	}
	var req struct {
		UserID  int64  `json:"user_id" binding:"required"`
		LevelID int64  `json:"level_id" binding:"required"`
		Days    int    `json:"days" binding:"required,min=1"`
		Notes   string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	grant, err := h.service.GrantManual(c.Request.Context(), req.UserID, req.LevelID, req.Days, req.Notes, fmt.Sprintf("admin:%d", subject.UserID))
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, grant)
}

func (h *MembershipHandler) RevokeGrant(c *gin.Context) {
	subject, ok := middleware2.GetAuthSubjectFromContext(c)
	if !ok {
		response.Unauthorized(c, "Administrator not found in context")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.BadRequest(c, "Invalid grant ID")
		return
	}
	if err := h.service.RevokeGrant(c.Request.Context(), id, fmt.Sprintf("admin:%d", subject.UserID)); err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, gin.H{"revoked": true})
}

func (h *MembershipHandler) AuditLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	items, err := h.service.ListAuditLogs(c.Request.Context(), limit)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}
	response.Success(c, items)
}
