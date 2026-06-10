package middleware

import (
	"net/http"
	"strings"

	"github.com/example/ems/pkg/jwt"
	"github.com/example/ems/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Context keys for values stored by the auth middleware.
const (
	ctxUserID = "auth_user_id"
	ctxEmail  = "auth_email"
	ctxRole   = "auth_role"
)

// Authenticate verifies the Bearer access token and stores the caller's
// identity in the gin context. Requests without a valid token are rejected.
func Authenticate(mgr *jwt.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		token := extractBearer(header)
		if token == "" {
			response.Fail(c, http.StatusUnauthorized, "missing or malformed Authorization header")
			c.Abort()
			return
		}
		claims, err := mgr.VerifyAccess(token)
		if err != nil {
			response.Fail(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}
		c.Set(ctxUserID, claims.UserID)
		c.Set(ctxEmail, claims.Email)
		c.Set(ctxRole, claims.Role)
		c.Next()
	}
}

// RequireRoles authorizes the request only if the caller's role is in roles.
// Must run after Authenticate.
func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}
	return func(c *gin.Context) {
		role, _ := c.Get(ctxRole)
		roleStr, _ := role.(string)
		if _, ok := allowed[roleStr]; !ok {
			response.Fail(c, http.StatusForbidden, "you do not have permission to perform this action")
			c.Abort()
			return
		}
		c.Next()
	}
}

// CurrentUserID returns the authenticated user's ID from the context.
func CurrentUserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(ctxUserID)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}

func extractBearer(header string) string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}
	return strings.TrimSpace(parts[1])
}
