package apis

import (
	"net/http"

	"github.com/MrWestbury/terraxen-naming-service/internals/engine"
	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type NamespaceApi struct {
	orgSvc    *services.OrganizationService
	nsSvc     *services.NamespaceService
	schemaSvc *services.SchemaService
}

func RegisterNamespaceApi(parentGroup *gin.RouterGroup, svc *services.NamespaceService, oSvc *services.OrganizationService, sSvc *services.SchemaService) {
	group := parentGroup.Group("/namespaces")
	nsApi := &NamespaceApi{
		nsSvc:     svc,
		orgSvc:    oSvc,
		schemaSvc: sSvc,
	}

	group.GET("/", nsApi.NotImplemented)
	group.POST("/", nsApi.NotImplemented)

	group.GET("/:ns", nsApi.NotImplemented)
	group.PUT("/:ns", nsApi.NotImplemented)
	group.DELETE("/:ns", nsApi.NotImplemented)

	group.GET("/:ns/variables", nsApi.NotImplemented)
	group.POST("/:ns/variables", nsApi.NotImplemented)
	group.GET("/:ns/variables/:var", nsApi.NotImplemented)
	group.PUT("/:ns/variables/:var", nsApi.NotImplemented)
	group.DELETE("/:ns/variables/:var", nsApi.NotImplemented)

	// Resolve a name
	group.GET("/:ns/resolve/:resource", nsApi.Resolve)

}

func (nsApi *NamespaceApi) NotImplemented(c *gin.Context) {
	responseError(c, http.StatusNotImplemented, "Not yet implemented")
}

func (nsApi *NamespaceApi) Resolve(c *gin.Context) {
	nsId := c.Param("ns")
	orgId := c.GetString("x-organization-id")
	resourceName := c.Param("resource")

	ns, err := nsApi.nsSvc.GetNamespaceById(nsId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "failed to get namespace")
	}

	org, err := nsApi.orgSvc.GetOrganizationById(orgId)
	schemaVersion, err := nsApi.schemaSvc.GetSchemaVersionById(ns.SchemaId, ns.SchemaVersion)
	resource, found := schemaVersion.Resources[resourceName]
	if !found {
		responseError(c, http.StatusNotFound, "resource name not found")
		return
	}

	resultVars := make(map[string]string)
	for k, v := range org.OrgVars {
		result, err := engine.ResolvePattern(v, resultVars)
		if err != nil {
			resultVars[k] = v
		}
		resultVars[k] = result
	}

	for k, v := range ns.Variables {
		result, err := engine.ResolvePattern(v, resultVars)
		if err != nil {
			resultVars[k] = v
		}
		resultVars[k] = result
	}

	resolved, err := engine.ResolvePattern(resource.Pattern, resultVars)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "failed to resolve resource")
		return
	}
	item := ResolveResourceResponse{
		ResourceName: resourceName,
		Pattern:      resource.Pattern,
		Value:        resolved,
	}
	responseSingleItem(c, item)
}
