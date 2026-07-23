package middleware

import (
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

const ContextKeyTicketPermission = "ticket_permission"

// RequireTicketPermission keeps ticket actions separated at the route layer.
// The current Sub2API role model grants all three capabilities to administrators;
// a future RBAC provider can replace this single check without changing handlers.
func RequireTicketPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if permission != service.TicketPermissionRead && permission != service.TicketPermissionReply && permission != service.TicketPermissionManage {
			AbortWithError(c, 403, "TICKET_PERMISSION_DENIED", "Ticket permission denied")
			return
		}
		role, ok := c.Get(string(ContextKeyUserRole))
		if !ok || role != service.RoleAdmin {
			AbortWithError(c, 403, "TICKET_PERMISSION_DENIED", "Ticket permission denied")
			return
		}
		c.Set(ContextKeyTicketPermission, permission)
		c.Next()
	}
}
