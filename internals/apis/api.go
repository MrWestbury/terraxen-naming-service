package apis

import (
	"github.com/MrWestbury/terraxen-naming-service/internals/config"
	"github.com/MrWestbury/terraxen-naming-service/internals/services/mongobackend"
	"github.com/gin-gonic/gin"
)

type Api struct {
	router *gin.Engine
}

func NewApi(config *config.Config) *Api {
	orgService := mongobackend.NewOrganizationService(config)
	nsService := mongobackend.NewNamespaceService(config)
	schemaService := mongobackend.NewSchemaService(config)
	apiKeyService := mongobackend.NewApiKeyService(config)

	api := &Api{
		router: gin.Default(),
	}

	middleware := NewMiddlewares(apiKeyService)

	apiGroup := api.router.Group("/api")
	apiGroup.Use(middleware.ValidateRequest)
	v1Group := apiGroup.Group("/v1")

	// Organization API
	orgHandler := NewOrganizationHandler(orgService)
	orgGroup := v1Group.Group("/organizations")
	orgGroup.GET("/", orgHandler.GetListOfOrganizations)
	orgGroup.POST("/", orgHandler.CreateOrganization)
	orgGroup.GET("/:orgId", orgHandler.GetOrganization)
	orgGroup.PUT("/:orgId", orgHandler.UpdateOrganization)
	orgGroup.DELETE("/:orgId")

	// Namespace API
	nsHandler := NewNamespaceHandler(nsService, orgService, schemaService)
	nsGroup := v1Group.Group("/namespaces")
	nsGroup.GET("/", nsHandler.ListNamespaces)
	nsGroup.POST("/", nsHandler.CreateNamespace)
	nsGroup.GET("/:ns", nsHandler.GetNamespace)
	nsGroup.PUT("/:ns", nsHandler.UpdateNamespace)
	nsGroup.DELETE("/:ns", nsHandler.DeleteNamespace)
	nsGroup.GET("/:ns/variables", nsHandler.ListNamespaceVariables)
	nsGroup.POST("/:ns/variables", nsHandler.PostNamespaceVariable)
	nsGroup.GET("/:ns/variables/:var", nsHandler.GetNamespaceVariable)
	nsGroup.PUT("/:ns/variables/:var", nsHandler.PutNamespaceVariable)
	nsGroup.DELETE("/:ns/variables/:var", nsHandler.DeleteNamespaceVariable)
	// Resolve a name
	nsGroup.GET("/:ns/resolve/:resource", nsHandler.Resolve)

	// Schema API
	schemaApiHandler := NewSchemaApiHandler(schemaService)
	schGroup := v1Group.Group("/schemas")
	schGroup.GET("/", schemaApiHandler.ListSchemas)
	schGroup.POST("/", schemaApiHandler.CreateSchema)

	schGroup.GET("/:schema", schemaApiHandler.GetSchema) // Schema details
	schGroup.PUT("/:schema", schemaApiHandler.UpdateSchema)
	schGroup.DELETE("/:schema", schemaApiHandler.DeleteSchema)

	schGroup.GET("/:schema/versions", schemaApiHandler.ListSchemaVersions)
	schGroup.POST("/:schema/versions", schemaApiHandler.CreateSchemaVersion)
	schGroup.GET("/:schema/versions/:version", schemaApiHandler.GetSchemaVersion)
	schGroup.PUT("/:schema/versions/:version", schemaApiHandler.UpdateSchemaVersion)
	schGroup.DELETE("/:schema/versions/:version", schemaApiHandler.DeleteSchemaVersion)
	schGroup.POST("/:schema/versions/:version/resolve", schemaApiHandler.ResolveResourceName)

	// API Key API
	apiKeyHandler := NewApiKeyHandler(apiKeyService)
	apiKeysGroup := v1Group.Group("/apikeys")
	apiKeysGroup.GET("/", apiKeyHandler.ListApiKeys)
	apiKeysGroup.POST("/", apiKeyHandler.CreateKey)

	apiKeysGroup.GET("/:apiId", apiKeyHandler.GetKey)
	apiKeysGroup.DELETE("/:apiId", apiKeyHandler.DeleteKey)

	return api
}

func (api *Api) Run(listener string) error {
	err := api.router.Run(listener)
	return err
}
