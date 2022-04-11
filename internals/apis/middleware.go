package apis

import (
	"net/http"

	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type Middlewares struct {
	apiSvc *services.ApiKeyService
}

func NewMiddlewares(apiKeySvc *services.ApiKeyService) *Middlewares {
	mdw := &Middlewares{
		apiSvc: apiKeySvc,
	}
	return mdw
}

func (m *Middlewares) ValidateRequest(c *gin.Context) {
	apiKey := c.GetHeader("X-Terraxen-API")
	if apiKey == "" {
		c.Set("x-organization-id", "")
	} else {
		ak := m.apiSvc.GetKey(apiKey)
		if ak == nil {
			responseError(c, http.StatusForbidden, "Api key rejected")
			return
		}
		c.Set("x-organization-id", ak.OrganizationId)
	}

	c.Next()
}
