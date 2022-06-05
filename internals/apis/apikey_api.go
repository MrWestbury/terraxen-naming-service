package apis

import (
	"net/http"

	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type ApiKeyApi struct {
	akSvc services.ApiKeyProvider
}

func RegisterApiKeyApi(parentGroup *gin.RouterGroup, apiKeySvc services.ApiKeyProvider) {
	group := parentGroup.Group("/apikeys")
	akApi := &ApiKeyApi{
		akSvc: apiKeySvc,
	}

	group.GET("/", akApi.ListApiKeys)
	group.POST("/", akApi.CreateKey)

	group.GET("/:apiId")
	group.DELETE("/:apiId", akApi.DeleteKey)
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

func (aka *ApiKeyApi) DeleteKey(c *gin.Context) {
	orgId := c.GetString("x-organization-id")
	if orgId == "" {
		responseError(c, http.StatusForbidden, "Valid API key for an organization required")
		return
	}

	apiId := c.Param("apidId")
	if apiId == "" {
		responseError(c, http.StatusBadRequest, "Invalid API ID")
		return
	}

	apikey := aka.akSvc.GetKey(apiId)
	if apikey == nil {
		responseError(c, http.StatusNotFound, "API key not found")
		return
	}

	aka.akSvc.DeleteKey(apiId)
	c.Status(http.StatusNoContent)
}
