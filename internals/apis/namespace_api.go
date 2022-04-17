package apis

import (
	"net/http"

	"github.com/MrWestbury/terraxen-naming-service/internals/engine"
	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type NamespaceApi struct {
	orgSvc    *services.OrganizationServiceInterface
	nsSvc     *services.NamespaceServiceInterface
	schemaSvc *services.SchemaService
}

func RegisterNamespaceApi(parentGroup *gin.RouterGroup, svc *services.NamespaceServiceInterface, oSvc *services.OrganizationServiceInterface, sSvc *services.SchemaService) {
	group := parentGroup.Group("/namespaces")
	nsApi := &NamespaceApi{
		nsSvc:     svc,
		orgSvc:    oSvc,
		schemaSvc: sSvc,
	}

	group.GET("/", NotImplemented)
	group.POST("/", NotImplemented)

	group.GET("/:ns", NotImplemented)
	group.PUT("/:ns", NotImplemented)
	group.DELETE("/:ns", NotImplemented)

	group.GET("/:ns/variables", NotImplemented)
	group.POST("/:ns/variables", NotImplemented)
	group.GET("/:ns/variables/:var", NotImplemented)
	group.PUT("/:ns/variables/:var", NotImplemented)
	group.DELETE("/:ns/variables/:var", NotImplemented)

	// Resolve a name
	group.GET("/:ns/resolve/:resource", nsApi.Resolve)

}

func (nsApi *NamespaceApi) ListNamespaces(c *gin.Context) {
	orgId := c.GetString("x-organization-id")

	nsList, err := nsApi.nsSvc.ListNamespaces(orgId)
}

func (nsApi *NamespaceApi) Resolve(c *gin.Context) {
	orgId := c.GetString("x-organization-id")

	nsId := c.Param("ns")
	resourceName := c.Param("resource")

	ns, err := nsApi.nsSvc.GetNamespaceById(nsId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to get namespace")
	}

	org, err := nsApi.orgSvc.GetOrganizationById(orgId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if org == nil {
		responseError(c, http.StatusNotFound, "Organization not found")
		return
	}

	schemaVersion, err := nsApi.schemaSvc.GetSchemaVersion(orgId, ns.SchemaId, ns.SchemaVersion)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
	}

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
