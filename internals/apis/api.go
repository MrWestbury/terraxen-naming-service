package apis

import (
	"github.com/MrWestbury/terraxen-naming-service/internals/config"
	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type Api struct {
	router *gin.Engine
}

func NewApi(config *config.Config) *Api {
	orgService := services.NewOrganizationService(config)
	nsService := services.NewNamespaceService(config)
	schemaService := services.NewSchemaService(config)
	apiKeyService := services.NewApiKeyService(config)

	api := &Api{
		router: gin.Default(),
	}

	middleware := NewMiddlewares(apiKeyService)

	apiGroup := api.router.Group("/api")
	apiGroup.Use(middleware.ValidateRequest)
	v1Group := apiGroup.Group("/v1")

	RegisterOrganizationApi(v1Group, orgService)
	RegisterNamespaceApi(v1Group, nsService, orgService, schemaService)
	RegisterSchemaApi(v1Group, schemaService)
	RegisterApiKeyApi(v1Group, apiKeyService)

	return api
}

func (api *Api) Run(listener string) error {
	err := api.router.Run(listener)
	return err
}
