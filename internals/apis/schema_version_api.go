package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MrWestbury/terraxen-naming-service/internals/engine"
	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

func (sApi *SchemaApi) ListSchemaVersions(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	schemaId := c.Param("schema")

	versions, err := sApi.schemaSvc.ListSchemaVersions(orgId, schemaId)
	if err != nil {
		if err == services.ErrSchemaNotFound {
			responseError(c, http.StatusNotFound, "Schema not found")
			return
		}
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseSingleItem(c, versions)
}

func (sApi *SchemaApi) CreateSchemaVersion(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	schemaId := c.Param("schema")

	requestCreateSchemaVersion := CreateSchemaVersionRequest{
		FromVersion: -1,
		Resources:   map[string]string{},
	}
	err := DecodeBody(c, &requestCreateSchemaVersion)
	if err != nil {
		responseError(c, http.StatusNotAcceptable, "Invalid request body")
		return
	}

	newResources := make(map[string]string)

	if requestCreateSchemaVersion.FromVersion > 0 {
		schemaVer, err := sApi.schemaSvc.GetSchemaVersion(orgId, schemaId, fmt.Sprintf("%d", requestCreateSchemaVersion.FromVersion))
		if err != nil {
			responseError(c, http.StatusInternalServerError, "Something went wrong")
			return
		}

		for k, v := range schemaVer.Resources {
			newResources[k] = v
		}
	}

	for k, v := range requestCreateSchemaVersion.Resources {
		newResources[k] = v
	}

	sv, err := sApi.schemaSvc.CreateSchemaVersion(orgId, schemaId, newResources, false)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseSingleItem(c, sv)
}

func (sApi *SchemaApi) GetSchemaVersion(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	schemaId := c.Param("schema")
	schemaVersionId := c.Param("version")

	sv, err := sApi.schemaSvc.GetSchemaVersion(orgId, schemaId, schemaVersionId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseSingleItem(c, sv)
}

func (sApi *SchemaApi) UpdateSchemaVersion(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	schemaId := c.Param("schema")
	schemaVersion := c.Param("version")

	var req UpdateSchemaVersionRequest
	err := DecodeBody(c, &req)
	if err != nil {
		responseError(c, http.StatusNotAcceptable, "Invalid request body")
		return
	}

	schemaVer, err := sApi.schemaSvc.GetSchemaVersion(orgId, schemaId, schemaVersion)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	if schemaVer.Published {
		responseError(c, http.StatusLocked, "Published schema versions cannot be updated")
		return
	}

	updatedVer, err := sApi.schemaSvc.UpdateSchemaVersion(orgId, schemaId, schemaVersion, req.Resources, req.Published)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	responseSingleItem(c, updatedVer)
}

func (sApi *SchemaApi) DeleteSchemaVersion(c *gin.Context) {
	// TODO: If we delete, do we actually delete, or just disable?
}

func (sApi *SchemaApi) ResolveResourceName(c *gin.Context) {
	orgId := c.GetString(ORG_CONTEXT_NAME)
	schemaId := c.Param("schema")
	schemaVersionId := c.Param("version")

	var resolveReq ResolveSchemaVersionRequest
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&resolveReq); err != nil {
		responseError(c, http.StatusBadRequest, "Unable to process request")
		return
	}

	sv, err := sApi.schemaSvc.GetSchemaVersion(orgId, schemaId, schemaVersionId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	res, found := sv.Resources[resolveReq.ResouceName]
	if !found {
		responseError(c, http.StatusNotFound, "Resource name not found in schema")
		return
	}

	result, err := engine.ResolvePattern(res, resolveReq.Variables)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	response := ResolveResourceResponse{
		ResourceName: resolveReq.ResouceName,
		Pattern:      res,
		Value:        result,
	}

	responseSingleItem(c, response)
}
