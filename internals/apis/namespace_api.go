package apis

import (
	"net/http"

	"github.com/MrWestbury/terraxen-naming-service/internals/engine"
	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type NamespaceHandler struct {
	orgSvc    services.OrganizationServiceProvider
	nsSvc     services.NamespaceServiceProvider
	schemaSvc services.SchemaServiceProvider
}

func NewNamespaceHandler(svc services.NamespaceServiceProvider, oSvc services.OrganizationServiceProvider, sSvc services.SchemaServiceProvider) *NamespaceHandler {
	nsApi := &NamespaceHandler{
		nsSvc:     svc,
		orgSvc:    oSvc,
		schemaSvc: sSvc,
	}

	return nsApi
}

func (nsApi *NamespaceHandler) ListNamespaces(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)

	nsList, err := nsApi.nsSvc.ListNamespaces(orgId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to get namespace")
		return
	}

	responseSingleItem(c, nsList)
}

func (nsApi *NamespaceHandler) CreateNamespace(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	var nsRequest NewNamespaceRequest

	if err := DecodeBody(c, &nsRequest); err != nil {
		return
	}

	exists := nsApi.nsSvc.ExistsByName(orgId, nsRequest.Name)
	if exists {
		responseError(c, http.StatusConflict, "Namespace already exists")
		return
	}

	ns, err := nsApi.nsSvc.CreateNamespace(orgId, nsRequest.Name, nsRequest.Schema, nsRequest.SchemaVersion, map[string]string{})
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to created namespace")
		return
	}

	responseSingleItem(c, ns)
}

func (nsApi *NamespaceHandler) GetNamespace(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	nsId := c.Param("ns")

	ns, err := nsApi.nsSvc.GetNamespaceById(orgId, nsId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to created namespace")
		return
	}

	responseSingleItem(c, ns)
}

func (nsApi *NamespaceHandler) Resolve(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)

	nsId := c.Param("ns")
	resourceName := c.Param("resource")

	ns, err := nsApi.nsSvc.GetNamespaceById(orgId, nsId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to get namespace")
		return
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

	nsVars, err := nsApi.nsSvc.GetVariablesAsMap(orgId, nsId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "something went wrong")
		return
	}

	for k, v := range nsVars {
		result, err := engine.ResolvePattern(v, resultVars)
		if err != nil {
			resultVars[k] = v
		}
		resultVars[k] = result
	}

	resolved, err := engine.ResolvePattern(resource, resultVars)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to resolve resource")
		return
	}
	item := ResolveResourceResponse{
		ResourceName: resourceName,
		Pattern:      resource,
		Value:        resolved,
	}
	responseSingleItem(c, item)
}

func (nsApi *NamespaceHandler) UpdateNamespace(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	nsId := c.Param("ns")

	var updateBody UpdateNamespaceRequest
	err := DecodeBody(c, &updateBody)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	ns, err := nsApi.nsSvc.GetNamespaceById(orgId, nsId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to get namespace")
		return
	}

	ns.Name = updateBody.Name
	ns.SchemaVersion = updateBody.SchemaVersion

	nsApi.nsSvc.UpdateNamespace(orgId, nsId, updateBody.Name, updateBody.SchemaVersion)
}

func (nsApi *NamespaceHandler) DeleteNamespace(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	nsId := c.Param("ns")

	err := nsApi.nsSvc.DeleteNamespace(orgId, nsId)
	if err != nil {
		if err == services.ErrNamespaceNotFound {
			responseError(c, http.StatusNotFound, "Namespace not found")
			return
		}
		responseError(c, http.StatusInternalServerError, "Failed to get namespace")
		return
	}

	responseNoContent(c, http.StatusNoContent)
}

func (nsApi *NamespaceHandler) ListNamespaceVariables(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	nsId := c.Param("ns")

	nsApi.nsSvc.ListNamespaceVars(orgId, nsId)
	ns, err := nsApi.nsSvc.ListNamespaceVars(orgId, nsId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to get namespace")
		return
	}

	responseSingleItem(c, ns)
}

func (nsApi *NamespaceHandler) PostNamespaceVariable(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	nsId := c.Param("ns")

	var reqBody NewNamespaceVariable
	err := DecodeBody(c, &reqBody)
	if err != nil {
		responseError(c, http.StatusUnprocessableEntity, "Invalid request body")
		return
	}

	ns, err := nsApi.nsSvc.GetNamespaceById(orgId, nsId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to get namespace")
		return
	}

	exists := nsApi.nsSvc.NamespaceVariableExists(orgId, nsId, reqBody.Name)
	if exists {
		responseError(c, http.StatusConflict, "Variable already exists")
		return
	}

	nsVar, err := nsApi.nsSvc.CreateNamespaceVariable(ns.OrganizationId, ns.Id, reqBody.Name, reqBody.Value)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to create namespace variable")
		return
	}

	responseSingleItem(c, nsVar)
}

func (nsApi *NamespaceHandler) GetNamespaceVariable(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	nsId := c.Param("ns")
	varId := c.Param("var")

	nsVar, err := nsApi.nsSvc.GetNamespaceVariable(orgId, nsId, varId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to get namespace variable")
		return
	}

	if nsVar == nil {
		responseError(c, http.StatusNotFound, "Namespace variable not found")
		return
	}

	responseSingleItem(c, nsVar)
}

func (nsApi *NamespaceHandler) PutNamespaceVariable(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	nsId := c.Param("ns")
	varId := c.Param("var")

	var reqBody UpdateNamespaceVariable
	err := DecodeBody(c, &reqBody)
	if err != nil {
		responseError(c, http.StatusUnprocessableEntity, "Invalid request body")
		return
	}

	nsVar, err := nsApi.nsSvc.UpdateNamespaceVariable(orgId, nsId, varId, reqBody.Value)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to update namespace variable")
		return
	}

	responseSingleItem(c, nsVar)
}

func (nsApi *NamespaceHandler) DeleteNamespaceVariable(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	nsId := c.Param("ns")
	varId := c.Param("var")

	err := nsApi.nsSvc.DeleteNamespaceVariable(orgId, nsId, varId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Failed to delete namespace variable")
		return
	}

	responseNoContent(c, http.StatusNoContent)
}
