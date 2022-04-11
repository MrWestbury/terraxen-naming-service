package apis

import (
	"net/http"

	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type ApiKeyApi struct {
	akSvc *services.ApiKeyService
}

func RegisterApiKeyApi(parentGroup *gin.RouterGroup, apiKeySvc *services.ApiKeyService) {
	group := parentGroup.Group("/apikeys")
	akApi := &ApiKeyApi{
		akSvc: apiKeySvc,
	}

	group.GET("/", akApi.ListApiKeys)
	group.POST("/", akApi.CreateKey)
}

func (aka *ApiKeyApi) ListApiKeys(c *gin.Context) {
	orgId := c.GetString("x-organization-id")
	if orgId == "" {
		responseError(c, http.StatusForbidden, "Valid API key for an organization required")
		return
	}

	keys, err := aka.akSvc.ListKeys(orgId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, err.Error())
	}
	responseSingleItem(c, keys)
}

func (aka *ApiKeyApi) CreateKey(c *gin.Context) {
	orgId := "a3dd56d4-c470-4f12-b47f-9f29a0380fc5"

	key, err := aka.akSvc.GenerateNewApiKey(orgId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, err.Error())
	}
	responseSingleItem(c, key)
}
