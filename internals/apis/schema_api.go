package apis

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/MrWestbury/terraxen-naming-service/internals/services"
	"github.com/gin-gonic/gin"
)

type SchemaApi struct {
	schemaSvc *services.SchemaService
}

func RegisterSchemaApi(parentGroup *gin.RouterGroup, svc *services.SchemaService) {
	group := parentGroup.Group("/schemas")

	schemaApi := &SchemaApi{
		schemaSvc: svc,
	}

	group.GET("/", schemaApi.ListSchemas)
	group.POST("/", schemaApi.CreateSchema)

	group.GET("/:schema", schemaApi.GetSchema) // Schema details
	group.PUT("/:schema", schemaApi.UpdateSchema)
	group.DELETE("/:schema", schemaApi.DeleteSchema)

	group.GET("/:schema/versions/:version", schemaApi.NotImplemented)
	group.PUT("/:schema/versions/:version", schemaApi.NotImplemented)
	group.DELETE("/:schema/versions/:version", schemaApi.NotImplemented)

	// Resolve a resource
	group.POST("/:schema/versions/:version/resolve", schemaApi.ResolveResourceName)
}

func (sApi *SchemaApi) ListSchemas(c *gin.Context) {
	orgId := c.GetString("x-organization-id")
	if orgId == "" {
		responseError(c, http.StatusForbidden, "Valid API key for an organization required")
		return
	}

	schemaList, err := sApi.schemaSvc.ListSchemaInOrganization(orgId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong getting schemas")
		return
	}

	responseSingleItem(c, schemaList)
}

// Create a new schema and a new version 1 for that schema
func (sApi *SchemaApi) CreateSchema(c *gin.Context) {
	orgId := c.GetString("x-organization-id")
	if orgId == "" {
		responseError(c, http.StatusForbidden, "Valid API key for an organization required")
		return
	}

	var schemaReq NewSchemaRequest
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&schemaReq); err != nil {
		responseError(c, http.StatusBadRequest, "Unable to process request")
		return
	}

	schema, err := sApi.schemaSvc.CreateSchema(orgId, schemaReq.Name)
	if err.Error() == "schema already exists" {
		responseError(c, http.StatusConflict, "Schema name already exists in organization")
	} else if err != nil {
		log.Printf("failed to create schema: %v", err)
		responseError(c, http.StatusInternalServerError, "Failed to create schema")
		return
	}

	responseSingleItem(c, schema)
}

func (sApi *SchemaApi) GetSchema(c *gin.Context) {
	orgId := c.GetString("x-organization-id")
	schemaId := c.Param("schema")

	schema, err := sApi.schemaSvc.GetSchemaById(orgId, schemaId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseSingleItem(c, schema)
}

func (sApi *SchemaApi) UpdateSchema(c *gin.Context) {
	orgId := c.GetString("x-organization-id")
	schemaId := c.Param("schema")

	var schemaReq UpdateSchemaRequest
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&schemaReq); err != nil {
		responseError(c, http.StatusBadRequest, "Unable to process request")
		return
	}

	schema, err := sApi.schemaSvc.GetSchemaById(orgId, schemaId)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}
	schema.Name = schemaReq.Name

	err = sApi.schemaSvc.UpdateSchema(*schema)
	if err != nil {
		responseError(c, http.StatusInternalServerError, "Something went wrong")
		return
	}

	responseSingleItem(c, schema)
}

func (sApi *SchemaApi) DeleteSchema(c *gin.Context) {
	orgId := c.GetString("x-organization-id")
	schemaId := c.Param("schema")

	err := sApi.schemaSvc.DeleteSchema(orgId, schemaId)
	if err != nil {
		if err.Error() == "schema not found" {
			responseError(c, http.StatusNotFound, "Schema not found")

		} else {
			responseError(c, http.StatusInternalServerError, "Something went wrong")
		}
		return
	}
}

func (sApi *SchemaApi) NotImplemented(c *gin.Context) {
	responseError(c, http.StatusNotImplemented, "Not yet implemented")
}

func (sApi *SchemaApi) ResolveResourceName(c *gin.Context) {

}
